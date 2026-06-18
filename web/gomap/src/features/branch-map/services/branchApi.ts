import type { Asset, AssetKind, Branch } from "../types";

/**
 * Fetches a branch layout. Requests go to "/api/..." which Vite proxies to the
 * backend at http://localhost:8080 (see vite.config.ts) to avoid CORS in dev.
 */
export async function fetchBranch(
  branchName: string,
  signal?: AbortSignal,
): Promise<Branch> {
  const res = await fetch(`/api/branches/${encodeURIComponent(branchName)}/layout`, {
    signal,
  });
  if (!res.ok) {
    throw new Error(`Failed to load branch "${branchName}" (${res.status})`);
  }
  return (await res.json()) as Branch;
}

interface AssetApiResponse {
  id: number;
  name: string;
  type: AssetKind;
  branch_id: string;
}

export async function fetchAssets(
  branchName: string,
  signal?: AbortSignal,
): Promise<Asset[]> {
  const res = await fetch(`/api/branches/${encodeURIComponent(branchName)}/assets`, {
    signal,
  });
  if (!res.ok) {
    throw new Error(`Failed to load assets for "${branchName}" (${res.status})`);
  }
  const raw = (await res.json()) as AssetApiResponse[];
  return raw.map((a) => ({ id: String(a.id), kind: a.type, name: a.name }));
}
