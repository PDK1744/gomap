package app

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/PDK1744/gomap/internal/config"
	"github.com/PDK1744/gomap/internal/handlers"
	"github.com/PDK1744/gomap/internal/service"
	"github.com/PDK1744/gomap/internal/storage"
	"github.com/PDK1744/gomap/internal/store"
	"github.com/PDK1744/gomap/internal/worker"
	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	HttpServer    *http.Server
	Branches      *store.BranchStore // needed for the sync worker
	branchHandler *handlers.BranchHandler
	sync          *worker.AssetSyncService
	mainPool      *pgxpool.Pool
	workerPool    *pgxpool.Pool
}

type Configs struct {
	apiCfg    *config.APIConfig
	dbCfg     *config.DBConfig
	workerCfg *config.WorkerDBConfig
}

func New() (*App, error) {
	apiCfg := config.NewAPIConfig()
	dbCfg := config.NewDBConfig()
	workerDBCfg := config.NewWorkerDBConfig()

	cfgs := &Configs{
		apiCfg:    apiCfg,
		dbCfg:     dbCfg,
		workerCfg: workerDBCfg,
	}

	ctx := context.Background()

	app := loadApp(ctx, cfgs)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/branches/{branch_name}/layout", app.branchHandler.GetBranchLayout)
	mux.HandleFunc("GET /api/branches/{branch_name}/asset/data", app.branchHandler.GetBranchAssetsAndAssignments) // will return Assets and their assignments
	// mux.HandleFunc("GET /api/branches/{branch_name}/assets", app.branchHandler.GetAssetsByBranch)
	// mux.HandleFunc("GET /api/branches/{branch_name}/assignments", app.branchHandler.GetAssetRoomAssignments) //load saved assignments on page open

	// PUT    /api/branches/{branch_name}/assignments        # debounced sync (replace full state)
	// DELETE /api/branches/{branch_name}/assignments/{room_id}  # unassign a specific room (optional)

	app.HttpServer = &http.Server{
		Addr:         ":" + apiCfg.Addr,
		Handler:      enableCORS(mux),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	return app, nil

}

func (a *App) Close() {
	a.mainPool.Close()
	a.workerPool.Close()
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173") // or "*" for local labs
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Instantly answer preflight checks
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (a *App) StartAssetSyncWorker(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	if err := a.sync.SyncAssets(ctx); err != nil {
		log.Printf("initial sync failed: %v", err)
	}

	for {
		select {
		case <-ctx.Done():
			log.Println("asset sync worker stopping")
			return
		case <-ticker.C:
			if err := a.sync.SyncAssets(ctx); err != nil {
				log.Printf("asset sync failed: %v", err)
			}
		}

	}
}

func loadApp(ctx context.Context, cfgs *Configs) *App {
	mainPool, err := storage.Connect(ctx, cfgs.dbCfg.ConnStr)
	if err != nil {
		log.Fatalf("main db pool failure: %v", err)
	}

	workerPool, err := storage.Connect(ctx, cfgs.workerCfg.ConnStr)
	if err != nil {
		mainPool.Close()
		log.Fatalf("worker db pool failure: %v", err)
	}
	// init App struct with services, handlers, etc.

	branchStore := store.NewBranchStore(mainPool)
	assetStore := store.NewAssetStore(mainPool)

	branchService := service.NewBranchService(branchStore)
	assetService := service.NewAssetService(assetStore)

	syncWorker := worker.NewAssetSyncService(mainPool, workerPool)

	branchHandler := handlers.NewHandler(branchService, assetService)

	return &App{
		Branches:      branchStore,
		branchHandler: branchHandler,
		sync:          syncWorker,
		mainPool:      mainPool,
		workerPool:    workerPool,
	}
}
