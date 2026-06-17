import { useEffect, useMemo, useState } from "react";
import type { Asset, Branch, BranchElement } from "../types";
import { computeBounds } from "../utils/geometry";
import { styleFor } from "../utils/style";
import { usePanZoom } from "../hooks/usePanZoom";
import { useElementSize } from "../hooks/useElementSize";
import { useAssetAssignments } from "../hooks/useAssetAssignments";
import { MapElement } from "./MapElement";
import { ContainerLabel } from "./ContainerLabel";
import { Legend } from "./Legend";
import { DetailsPanel } from "./DetailsPanel";
import { AssetList } from "./AssetList";

interface BranchMapProps {
  branch: Branch;
}

const PADDING = 80; // world-unit margin around the layout when fitting to view
const DETAIL_MIN_PX = 120; // on-screen cell width at which asset content appears
const NO_ASSETS: Asset[] = [];

export function BranchMap({ branch }: BranchMapProps) {
  const { ref, size } = useElementSize<HTMLDivElement>();
  const [selected, setSelected] = useState<BranchElement | null>(null);
  const [showRoomNames, setShowRoomNames] = useState(true);
  const { assetsByElement, available, assign, unassign } =
    useAssetAssignments();

  const bounds = useMemo(
    () => computeBounds(branch.elements),
    [branch.elements],
  );

  // Draw large containers first so small fixtures (cubicles, doors) sit on top.
  const ordered = useMemo(
    () =>
      [...branch.elements].sort(
        (a, b) => styleFor(a.type).z - styleFor(b.type).z,
      ),
    [branch.elements],
  );

  // Room names render on a separate top layer so they overlay their children.
  const overlayLabels = useMemo(
    () => branch.elements.filter((e) => styleFor(e.type).label === "overlay"),
    [branch.elements],
  );

  const fitTransform = useMemo(() => {
    if (size.width === 0 || size.height === 0) {
      return { scale: 1, tx: 0, ty: 0 };
    }
    const scale = Math.min(
      size.width / (bounds.width + PADDING * 2),
      size.height / (bounds.height + PADDING * 2),
    );
    const tx = (size.width - bounds.width * scale) / 2 - bounds.minX * scale;
    const ty = (size.height - bounds.height * scale) / 2 - bounds.minY * scale;
    return { scale, tx, ty };
  }, [size, bounds]);

  const { transform, setTransform, isPanning, wasDrag, handlers } =
    usePanZoom(fitTransform);

  // Re-fit whenever the branch or container size changes (e.g. first measure).
  const fitKey = `${branch.branch_name}:${size.width}x${size.height}`;
  const [lastFit, setLastFit] = useState("");
  useEffect(() => {
    if (fitKey !== lastFit && size.width > 0) {
      setTransform(fitTransform);
      setLastFit(fitKey);
    }
  }, [fitKey, lastFit, fitTransform, size.width, setTransform]);

  const zoomBy = (factor: number) =>
    setTransform((t) => {
      const cx = size.width / 2;
      const cy = size.height / 2;
      const scale = Math.min(8, Math.max(0.05, t.scale * factor));
      const ratio = scale / t.scale;
      return {
        scale,
        tx: cx - (cx - t.tx) * ratio,
        ty: cy - (cy - t.ty) * ratio,
      };
    });

  const assignableCount = branch.elements.filter((e) => e.assignable).length;

  return (
    <div className="relative flex h-full w-full overflow-hidden bg-slate-100">
      <div ref={ref} className="relative h-full flex-1">
        <svg
          className="h-full w-full touch-none"
          style={{ cursor: isPanning ? "grabbing" : "grab" }}
          onClick={() => {
            if (!wasDrag()) setSelected(null);
          }}
          {...handlers}
        >
          <g
            transform={`translate(${transform.tx} ${transform.ty}) scale(${transform.scale})`}
          >
            {ordered.map((el) => (
              <MapElement
                key={el.tru_id}
                element={el}
                selected={selected?.tru_id === el.tru_id}
                detailed={transform.scale * el.width >= DETAIL_MIN_PX}
                assets={assetsByElement[el.tru_id] ?? NO_ASSETS}
                onSelect={setSelected}
                onDropAsset={assign}
              />
            ))}
            {/* Top layer: room names float above the cubicles inside them. */}
            {showRoomNames &&
              overlayLabels.map((el) => (
                <ContainerLabel key={`label-${el.tru_id}`} element={el} />
              ))}
          </g>
        </svg>

        {/* Zoom controls */}
        <div className="absolute bottom-4 left-4 flex flex-col overflow-hidden rounded-lg border border-slate-200 bg-white shadow-sm">
          <button
            onClick={() => zoomBy(1.25)}
            className="px-3 py-2 text-lg leading-none text-slate-600 hover:bg-slate-100"
            aria-label="Zoom in"
          >
            +
          </button>
          <button
            onClick={() => zoomBy(0.8)}
            className="border-t border-slate-200 px-3 py-2 text-lg leading-none text-slate-600 hover:bg-slate-100"
            aria-label="Zoom out"
          >
            −
          </button>
          <button
            onClick={() => setTransform(fitTransform)}
            className="border-t border-slate-200 px-3 py-1.5 text-xs text-slate-600 hover:bg-slate-100"
            aria-label="Fit to view"
          >
            Fit
          </button>
        </div>
      </div>

      {/* Sidebar */}
      <aside className="flex w-72 flex-col gap-3 overflow-y-auto border-l border-slate-200 bg-slate-50 p-4">
        <div>
          <h2 className="text-lg font-semibold capitalize text-slate-800">
            {branch.branch_name}
          </h2>
          <p className="text-sm text-slate-500">
            {branch.elements.length} spaces · {assignableCount} assignable
          </p>
        </div>
        <label className="flex items-center gap-2 rounded-lg border border-slate-200 bg-white/90 px-3 py-2 text-sm text-slate-700 shadow-sm">
          <input
            type="checkbox"
            checked={showRoomNames}
            onChange={(e) => setShowRoomNames(e.target.checked)}
            className="h-4 w-4 accent-indigo-500"
          />
          Show room names
        </label>
        <DetailsPanel
          element={selected}
          assets={selected ? assetsByElement[selected.tru_id] ?? NO_ASSETS : NO_ASSETS}
          available={available}
          onAssign={assign}
          onUnassign={unassign}
          onClose={() => setSelected(null)}
        />
        <AssetList available={available} />
        <Legend />
        <p className="mt-auto text-xs text-slate-400">
          Scroll to zoom · drag to pan · click a space for details.
        </p>
      </aside>
    </div>
  );
}
