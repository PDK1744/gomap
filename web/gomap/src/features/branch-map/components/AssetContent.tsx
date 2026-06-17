import { memo } from "react";
import type { Asset, BranchElement } from "../types";
import { ASSET_ICON } from "../data/assets";
import { displayName } from "../utils/style";

interface AssetContentProps {
  element: BranchElement;
  assets: Asset[];
  /** Draw the space's name as a header. Off for rooms — the overlay pill already names them. */
  showName: boolean;
}

/** Header text: cubicles show their number, others their (de-prefixed) name. */
function header(element: BranchElement): string {
  const name = displayName(element.name);
  if (element.type === "cubicle") return name.match(/\d+/)?.[0] ?? name;
  return name;
}

const MAX_LINE_H = 34; // world units; keeps rows top-anchored in tall rooms

/**
 * Rendered inside a cell once it's zoomed in enough to be legible (see the
 * detail threshold in BranchMap). Shows the space's name and its assets as an
 * icon + label list, sized as a fraction of the cell so it scales with zoom.
 */
function AssetContentImpl({ element, assets, showName }: AssetContentProps) {
  const { x, y, width, height } = element;
  const pad = Math.min(width, height) * 0.12;
  const left = x + pad;
  const availH = height - pad * 2;
  const availW = width - pad * 2;

  // header (optional) + asset rows, with "Empty" standing in on leaf cells
  const rowCount = Math.max(assets.length, showName ? 1 : 0);
  const lines = (showName ? 1 : 0) + rowCount;
  if (lines === 0) return null;

  const lineH = Math.min(availH / lines, MAX_LINE_H);
  const headFont = Math.min(
    lineH * 0.78,
    availW / Math.max(header(element).length * 0.6, 1),
    22,
  );
  const rowFont = Math.min(lineH * 0.66, 15);
  const rowsTop = y + pad + (showName ? lineH : 0);

  return (
    <g pointerEvents="none" style={{ userSelect: "none" }}>
      {showName && (
        <text
          x={left}
          y={y + pad + lineH * 0.5}
          fontSize={headFont}
          fontWeight={700}
          dominantBaseline="central"
          fill="#0f172a"
        >
          {header(element)}
        </text>
      )}

      {assets.length === 0 ? (
        showName && (
          <text
            x={left}
            y={rowsTop + lineH * 0.5}
            fontSize={rowFont}
            fontStyle="italic"
            dominantBaseline="central"
            fill="#94a3b8"
          >
            Empty
          </text>
        )
      ) : (
        assets.map((asset, i) => {
          const cy = rowsTop + lineH * (i + 0.5);
          return (
            <g key={asset.id}>
              <text x={left} y={cy} fontSize={rowFont * 1.05} dominantBaseline="central">
                {ASSET_ICON[asset.kind]}
              </text>
              <text
                x={left + rowFont * 1.5}
                y={cy}
                fontSize={rowFont}
                dominantBaseline="central"
                fill="#334155"
              >
                {asset.name}
              </text>
            </g>
          );
        })
      )}
    </g>
  );
}

export const AssetContent = memo(AssetContentImpl);
