import { BranchMapView } from "./features/branch-map";

function App() {
  return (
    <div className="flex h-screen w-screen flex-col bg-slate-100">
      <header className="flex items-center justify-between border-b border-slate-200 bg-white px-6 py-3">
        <h1 className="text-lg font-semibold text-slate-800">GoMap</h1>
        <span className="text-sm text-slate-400">Office Layout</span>
      </header>
      <main className="min-h-0 flex-1">
        <BranchMapView branchName="stoneledge" />
      </main>
    </div>
  );
}

export default App;
