import { useCallback, useMemo, useState } from "react";
import type { Asset, Assignments } from "../types";
import { DUMMY_ASSETS, INITIAL_ASSIGNMENTS } from "../data/assets";

/**
 * Holds the asset pool and which space each asset is assigned to. An asset can
 * live in at most one space, so assigning it elsewhere moves it.
 */
export function useAssetAssignments() {
  const [assignments, setAssignments] =
    useState<Assignments>(INITIAL_ASSIGNMENTS);

  const byId = useMemo(() => {
    const map = new Map<string, Asset>();
    for (const asset of DUMMY_ASSETS) map.set(asset.id, asset);
    return map;
  }, []);

  // tru_id -> resolved Asset[]; stable per element unless its list changes.
  const assetsByElement = useMemo(() => {
    const result: Record<string, Asset[]> = {};
    for (const [elementId, ids] of Object.entries(assignments)) {
      result[elementId] = ids
        .map((id) => byId.get(id))
        .filter((a): a is Asset => Boolean(a));
    }
    return result;
  }, [assignments, byId]);

  const assignedIds = useMemo(
    () => new Set(Object.values(assignments).flat()),
    [assignments],
  );

  const available = useMemo(
    () => DUMMY_ASSETS.filter((a) => !assignedIds.has(a.id)),
    [assignedIds],
  );

  const assign = useCallback((elementId: string, assetId: string) => {
    setAssignments((prev) => {
      const next: Assignments = {};
      // Remove the asset from any space that currently holds it.
      for (const [id, ids] of Object.entries(prev)) {
        const filtered = ids.filter((a) => a !== assetId);
        if (filtered.length > 0) next[id] = filtered;
      }
      next[elementId] = [...(next[elementId] ?? []), assetId];
      return next;
    });
  }, []);

  const unassign = useCallback((elementId: string, assetId: string) => {
    setAssignments((prev) => {
      const current = prev[elementId];
      if (!current) return prev;
      const filtered = current.filter((a) => a !== assetId);
      const next = { ...prev };
      if (filtered.length > 0) next[elementId] = filtered;
      else delete next[elementId];
      return next;
    });
  }, []);

  return { assetsByElement, available, assign, unassign };
}
