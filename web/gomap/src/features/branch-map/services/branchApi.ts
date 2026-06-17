import type { Branch } from "../types";

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
