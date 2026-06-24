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
	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	HttpServer *http.Server
	Rooms      *store.RoomStore
	Assets     *store.AssetStore
	Branches   *store.BranchStore

	Sync *service.AssetSyncService
	// May noy need the pools in the App struct
	// I'm passing the pool into the stores so unless I need the pools elsewhere.. I can remove from App
	mainPool   *pgxpool.Pool
	workerPool *pgxpool.Pool
}

func New() (*App, error) {
	apiCfg := config.NewAPIConfig()
	dbCfg := config.NewDBConfig()
	workerDBCfg := config.NewWorkerDBConfig()

	ctx := context.Background()

	mainPool, err := storage.Connect(ctx, dbCfg.ConnStr)
	if err != nil {
		log.Fatalf("main db pool failure: %v", err)
	}

	workerPool, err := storage.Connect(ctx, workerDBCfg.ConnStr)
	if err != nil {
		mainPool.Close()
		log.Fatalf("worker db pool failure: %v", err)
	}

	app := &App{
		Rooms:      store.NewRoomStore(mainPool),
		Assets:     store.NewAssetStore(mainPool),
		Branches:   store.NewBranchStore(mainPool),
		Sync:       service.NewAssetSyncService(mainPool, workerPool),
		mainPool:   mainPool,
		workerPool: workerPool,
	}
	branchHandler := handlers.NewBranchHandler(app.Branches, app.Assets)
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/branches/{branch_name}/layout", branchHandler.GetBranchLayout)
	mux.HandleFunc("GET /api/branches/{branch_name}/assets", branchHandler.GetAssetsByBranch)
	mux.HandleFunc("GET /api/branches/{branch_name}/assignments", branchHandler.GetAssetRoomAssignments) //load saved assignments on page open

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

	if err := a.Sync.SyncAssets(ctx); err != nil {
		log.Printf("initial sync failed: %v", err)
	}

	for {
		select {
		case <-ctx.Done():
			log.Println("asset sync worker stopping")
			return
		case <-ticker.C:
			if err := a.Sync.SyncAssets(ctx); err != nil {
				log.Printf("asset sync failed: %v", err)
			}
		}

	}
}
