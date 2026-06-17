import type { Asset, AssetKind, BranchElement } from "../types";
import { ASSET_ICON, ASSET_KIND_LABEL } from "../data/assets";
import { canHoldAssets, TYPE_LABELS } from "../utils/style";
import { GeneratorIcon } from "./GeneratorBadge";

interface DetailsPanelProps {
  element: BranchElement | null;
  assets: Asset[];
  available: Asset[];
  onAssign: (elementId: string, assetId: string) => void;
  onUnassign: (elementId: string, assetId: string) => void;
  onClose: () => void;
}

const KIND_ORDER: AssetKind[] = ["pc", "phone", "printer"];

export function DetailsPanel({
  element,
  assets,
  available,
  onAssign,
  onUnassign,
  onClose,
}: DetailsPanelProps) {
  if (!element) {
    return (
      <div className="rounded-lg border border-slate-200 bg-white/90 p-4 text-sm text-slate-500 shadow-sm backdrop-blur">
        Click any space on the map to see its details.
      </div>
    );
  }

  const holdsAssets = canHoldAssets(element);

  return (
    <div className="rounded-lg border border-slate-200 bg-white/90 p-4 shadow-sm backdrop-blur">
      <div className="flex items-start justify-between gap-2">
        <h3 className="text-base font-semibold text-slate-800">{element.name}</h3>
        <button
          onClick={onClose}
          className="rounded p-1 text-slate-400 hover:bg-slate-100 hover:text-slate-600"
          aria-label="Close details"
        >
          ✕
        </button>
      </div>

      <dl className="mt-3 space-y-2 text-sm">
        <Row label="Type" value={TYPE_LABELS[element.type]} />
        <Row label="ID" value={element.tru_id} mono />
        <Row label="Size" value={`${element.width} × ${element.height}`} mono />
        {element.generator && (
          <Row
            label="Power"
            value={
              <span className="flex items-center gap-1.5 rounded-full bg-amber-100 px-2 py-0.5 text-xs font-medium text-amber-700">
                <GeneratorIcon size={12} />
                Generator
              </span>
            }
          />
        )}
      </dl>

      {holdsAssets && (
        <div className="mt-4 border-t border-slate-200 pt-3">
          <h4 className="mb-2 text-xs font-semibold uppercase tracking-wide text-slate-500">
            Assets ({assets.length})
          </h4>

          {assets.length === 0 ? (
            <p className="text-sm italic text-slate-400">No assets assigned.</p>
          ) : (
            <ul className="space-y-1.5">
              {assets.map((asset) => (
                <li
                  key={asset.id}
                  className="flex items-center justify-between rounded-md bg-slate-100 px-2 py-1.5 text-sm"
                >
                  <span className="flex items-center gap-2 text-slate-700">
                    <span>{ASSET_ICON[asset.kind]}</span>
                    {asset.name}
                  </span>
                  <button
                    onClick={() => onUnassign(element.tru_id, asset.id)}
                    className="rounded p-0.5 text-slate-400 hover:bg-slate-200 hover:text-red-500"
                    aria-label={`Remove ${asset.name}`}
                  >
                    ✕
                  </button>
                </li>
              ))}
            </ul>
          )}

          <AssetPicker
            available={available}
            onPick={(assetId) => onAssign(element.tru_id, assetId)}
          />
        </div>
      )}
    </div>
  );
}

function AssetPicker({
  available,
  onPick,
}: {
  available: Asset[];
  onPick: (assetId: string) => void;
}) {
  return (
    <select
      value=""
      onChange={(e) => {
        if (e.target.value) onPick(e.target.value);
        e.target.value = "";
      }}
      className="mt-2 w-full rounded-md border border-slate-300 bg-white px-2 py-1.5 text-sm text-slate-700 focus:border-indigo-400 focus:outline-none"
    >
      <option value="">+ Assign an asset…</option>
      {KIND_ORDER.map((kind) => {
        const items = available.filter((a) => a.kind === kind);
        if (items.length === 0) return null;
        return (
          <optgroup key={kind} label={ASSET_KIND_LABEL[kind]}>
            {items.map((asset) => (
              <option key={asset.id} value={asset.id}>
                {asset.name}
              </option>
            ))}
          </optgroup>
        );
      })}
    </select>
  );
}

function Row({
  label,
  value,
  mono,
}: {
  label: string;
  value: React.ReactNode;
  mono?: boolean;
}) {
  return (
    <div className="flex items-center justify-between gap-3">
      <dt className="text-slate-500">{label}</dt>
      <dd className={mono ? "font-mono text-slate-700" : "text-slate-700"}>{value}</dd>
    </div>
  );
}
