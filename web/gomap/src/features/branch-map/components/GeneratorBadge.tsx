import { memo } from "react";
import type { BranchElement } from "../types";

export const GENERATOR_COLOR = "#f59e0b"; // amber-500

/** Lightning bolt path in a 24×24 box, shared by the map badge and legend icon. */
const BOLT_PATH = "M13 2 L5 13.5 H11 L9.5 22 L19 9.5 H12.5 L13 2 Z";

/**
 * Amber lightning badge pinned to the top-right corner of a space that is
 * backed by generator power. Sized relative to the cell so it scales with zoom.
 */
function GeneratorBadgeImpl({ element }: { element: BranchElement }) {
  const { x, y, width, height } = element;
  const r = Math.min(Math.min(width, height) * 0.18, 26);
  const cx = x + width - r * 1.35;
  const cy = y + r * 1.35;

  return (
    <g pointerEvents="none">
      <circle cx={cx} cy={cy} r={r} fill={GENERATOR_COLOR} stroke="#ffffff" strokeWidth={r * 0.12} />
      <path
        d={BOLT_PATH}
        fill="#ffffff"
        transform={`translate(${cx} ${cy}) scale(${(r * 1.4) / 24}) translate(-12 -12)`}
      />
    </g>
  );
}

export const GeneratorBadge = memo(GeneratorBadgeImpl);

/** Small inline version of the same badge for HTML contexts (legend, panels). */
export function GeneratorIcon({ size = 14 }: { size?: number }) {
  return (
    <svg width={size} height={size} viewBox="0 0 24 24" aria-hidden>
      <circle cx={12} cy={12} r={11} fill={GENERATOR_COLOR} />
      <path d={BOLT_PATH} fill="#ffffff" transform="translate(12 12) scale(0.7) translate(-12 -12)" />
    </svg>
  );
}
