import { useMemo, useState } from "react";
import type { Asset, AssetKind } from "../types";
import { ASSET_DRAG_TYPE, ASSET_ICON, ASSET_KIND_LABEL } from "../data/assets";

interface AssetListProps {
  /** Unassigned assets, available to drag onto the map. */
  available: Asset[];
}

const KIND_ORDER: AssetKind[] = ["pc", "phone", "printer"];

/**
 * Searchable list of unassigned assets. Each row is draggable — drop it on a
 * cubicle, office, or counter to assign it there.
 */
export function AssetList({ available }: AssetListProps) {
  const [query, setQuery] = useState("");

  const groups = useMemo(() => {
    const q = query.trim().toLowerCase();
    const matches = q
      ? available.filter((a) => a.name.toLowerCase().includes(q))
      : available;
    return KIND_ORDER.map((kind) => ({
      kind,
      items: matches.filter((a) => a.kind === kind),
    })).filter((g) => g.items.length > 0);
  }, [available, query]);

  return (
    <div className="flex min-h-0 flex-col rounded-lg border border-slate-200 bg-white/90 shadow-sm backdrop-blur">
      <div className="border-b border-slate-200 p-3">
        <h3 className="mb-2 text-xs font-semibold uppercase tracking-wide text-slate-500">
          Available Assets ({available.length})
        </h3>
        <input
          type="search"
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          placeholder="Search PC500, Ext. 1021…"
          className="w-full rounded-md border border-slate-300 bg-white px-2 py-1.5 text-sm text-slate-700 placeholder:text-slate-400 focus:border-indigo-400 focus:outline-none"
        />
        <p className="mt-1.5 text-xs text-slate-400">
          Drag an asset onto any assignable space (cubicle, office, room…).
        </p>
      </div>

      <div className="max-h-64 overflow-y-auto p-2">
        {groups.length === 0 ? (
          <p className="px-1 py-2 text-sm italic text-slate-400">
            {available.length === 0 ? "All assets are assigned." : "No matches."}
          </p>
        ) : (
          groups.map(({ kind, items }) => (
            <div key={kind} className="mb-2 last:mb-0">
              <h4 className="px-1 py-1 text-xs font-medium text-slate-400">
                {ASSET_KIND_LABEL[kind]} ({items.length})
              </h4>
              <ul className="space-y-1">
                {items.map((asset) => (
                  <li
                    key={asset.id}
                    draggable
                    onDragStart={(e) => {
                      e.dataTransfer.setData(ASSET_DRAG_TYPE, asset.id);
                      e.dataTransfer.effectAllowed = "move";
                    }}
                    className="flex cursor-grab items-center gap-2 rounded-md border border-slate-200 bg-white px-2 py-1.5 text-sm text-slate-700 shadow-sm hover:border-indigo-300 hover:bg-indigo-50 active:cursor-grabbing"
                  >
                    <span>{ASSET_ICON[asset.kind]}</span>
                    {asset.name}
                  </li>
                ))}
              </ul>
            </div>
          ))
        )}
      </div>
    </div>
  );
}
