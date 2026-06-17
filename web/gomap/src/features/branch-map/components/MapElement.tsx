import { memo, useState } from "react";
import type { Asset, BranchElement } from "../types";
import { ASSET_DRAG_TYPE } from "../data/assets";
import {
  canHoldAssets,
  displayName,
  LEAF_CARD_TYPES,
  styleFor,
} from "../utils/style";
import { AssetContent } from "./AssetContent";
import { GeneratorBadge } from "./GeneratorBadge";

interface MapElementProps {
  element: BranchElement;
  selected: boolean;
  detailed: boolean;
  assets: Asset[];
  onSelect: (element: BranchElement) => void;
  onDropAsset: (elementId: string, assetId: string) => void;
}

/** Short in-cell label: cubicles show just their number, others their name. */
function inlineLabel(element: BranchElement): string {
  const name = displayName(element.name);
  if (element.type === "cubicle") {
    return name.match(/\d+/)?.[0] ?? name;
  }
  return name;
}

function MapElementImpl({
  element,
  selected,
  detailed,
  assets,
  onSelect,
  onDropAsset,
}: MapElementProps) {
  const style = styleFor(element.type);
  const { x, y, width, height } = element;
  const interactive = element.assignable;
  const droppable = canHoldAssets(element);
  // Leaf cells always show their card when zoomed in; containers (rooms,
  // branch) only once they hold assets, to avoid placeholder clutter.
  const isLeafCard = LEAF_CARD_TYPES.has(element.type);
  const showContent =
    detailed && droppable && (isLeafCard || assets.length > 0);
  const [dragOver, setDragOver] = useState(false);

  // Compact label (shown when zoomed out): anchored along the bottom of the cell.
  const label = inlineLabel(element);
  const padding = Math.min(width, height) * 0.12;
  const fontSize = Math.min(
    height * 0.24,
    (width - padding * 2) / Math.max(label.length * 0.58, 1),
    26,
  );

  const stroke = dragOver ? "#6366f1" : selected ? "#6366f1" : style.stroke;
  const strokeWidth =
    dragOver || selected ? style.strokeWidth + 3 : style.strokeWidth;

  return (
    <g
      onClick={(e) => {
        e.stopPropagation();
        onSelect(element);
      }}
      onDragOver={
        droppable
          ? (e) => {
              e.preventDefault();
              e.dataTransfer.dropEffect = "move";
              setDragOver(true);
            }
          : undefined
      }
      onDragLeave={droppable ? () => setDragOver(false) : undefined}
      onDrop={
        droppable
          ? (e) => {
              e.preventDefault();
              setDragOver(false);
              const assetId = e.dataTransfer.getData(ASSET_DRAG_TYPE);
              if (assetId) onDropAsset(element.tru_id, assetId);
            }
          : undefined
      }
      style={{ cursor: interactive ? "pointer" : "default" }}
      role={interactive ? "button" : undefined}
      aria-label={element.name}
    >
      <rect
        x={x}
        y={y}
        width={width}
        height={height}
        rx={element.type === "cubicle" || element.type === "counter" ? 4 : 8}
        fill={dragOver ? "#e0e7ff" : style.fill}
        stroke={stroke}
        strokeWidth={strokeWidth}
        opacity={interactive ? 1 : 0.85}
      />
      {showContent ? (
        <AssetContent element={element} assets={assets} showName={isLeafCard} />
      ) : (
        style.label === "inline" &&
        fontSize >= 6 && (
          <text
            x={x + width / 2}
            y={y + height - padding - fontSize * 0.35}
            textAnchor="middle"
            dominantBaseline="central"
            fontSize={fontSize}
            fill="#1e293b"
            fontWeight={element.type === "cubicle" ? 600 : 500}
            pointerEvents="none"
            style={{ userSelect: "none" }}
          >
            {label}
          </text>
        )
      )}
      {element.generator && <GeneratorBadge element={element} />}
    </g>
  );
}

export const MapElement = memo(MapElementImpl);
