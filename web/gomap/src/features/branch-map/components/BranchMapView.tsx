import { useBranch } from "../hooks/useBranch";
import { useAssets } from "../hooks/useAssets";
import { BranchMap } from "./BranchMap";

interface BranchMapViewProps {
  branchName: string;
}

/** Loads a branch and renders the map, with loading/error fallbacks. */
export function BranchMapView({ branchName }: BranchMapViewProps) {
  const { branch, loading: branchLoading, error: branchError } = useBranch(branchName);
  const { assets, loading: assetsLoading, error: assetsError } = useAssets(branchName);

  if (branchLoading || assetsLoading) {
    return (
      <Centered>
        <div className="flex items-center gap-3 text-slate-500">
          <span className="h-4 w-4 animate-spin rounded-full border-2 border-slate-300 border-t-slate-600" />
          Loading {branchName}…
        </div>
      </Centered>
    );
  }

  if (branchError || assetsError || !branch) {
    return (
      <Centered>
        <div className="max-w-md text-center">
          <p className="text-lg font-semibold text-slate-800">
            Couldn’t load this branch
          </p>
          <p className="mt-1 text-sm text-slate-500">{branchError ?? assetsError}</p>
          <p className="mt-3 text-xs text-slate-400">
            Make sure the backend is running at http://localhost:8080.
          </p>
        </div>
      </Centered>
    );
  }

  return <BranchMap branch={branch} assets={assets} />;
}

function Centered({ children }: { children: React.ReactNode }) {
  return (
    <div className="flex h-full w-full items-center justify-center bg-slate-100">
      {children}
    </div>
  );
}
