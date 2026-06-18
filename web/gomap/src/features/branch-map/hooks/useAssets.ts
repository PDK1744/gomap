import { useEffect, useState } from "react";
import type { Asset } from "../types";
import { fetchAssets } from "../services/branchApi";

interface AssetsState {
  assets: Asset[];
  loading: boolean;
  error: string | null;
}

export function useAssets(branchName: string): AssetsState {
  const [state, setState] = useState<AssetsState>({
    assets: [],
    loading: true,
    error: null,
  });

  useEffect(() => {
    const controller = new AbortController();
    setState({ assets: [], loading: true, error: null });

    fetchAssets(branchName, controller.signal)
      .then((assets) => setState({ assets, loading: false, error: null }))
      .catch((err: unknown) => {
        if (controller.signal.aborted) return;
        const message = err instanceof Error ? err.message : "Unknown error";
        setState({ assets: [], loading: false, error: message });
      });

    return () => controller.abort();
  }, [branchName]);

  return state;
}
