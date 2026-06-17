import { useEffect, useState } from "react";
import type { Branch } from "../types";
import { fetchBranch } from "../services/branchApi";

interface BranchState {
  branch: Branch | null;
  loading: boolean;
  error: string | null;
}

export function useBranch(branchName: string): BranchState {
  const [state, setState] = useState<BranchState>({
    branch: null,
    loading: true,
    error: null,
  });

  useEffect(() => {
    const controller = new AbortController();
    setState({ branch: null, loading: true, error: null });

    fetchBranch(branchName, controller.signal)
      .then((branch) => setState({ branch, loading: false, error: null }))
      .catch((err: unknown) => {
        if (controller.signal.aborted) return;
        const message = err instanceof Error ? err.message : "Unknown error";
        setState({ branch: null, loading: false, error: message });
      });

    return () => controller.abort();
  }, [branchName]);

  return state;
}
