import { memo } from "react";
import type { BranchElement } from "../types";
import { displayName } from "../utils/style";

interface ContainerLabelProps {
  element: BranchElement;
}

/**
 * A room's name, drawn on a top layer above its children. A translucent pill
 * keeps the text legible while letting the cubicles underneath show through.
 */
function ContainerLabelImpl({ element }: ContainerLabelProps) {
  const { x, y, width, height } = element;
  const label = displayName(element.name);

  const fontSize = Math.min(width * 0.06, height * 0.09, 32);
  if (fontSize < 7) return null;

  // Estimated text box, anchored near the top-center of the room.
  const padX = fontSize * 0.7;
  const padY = fontSize * 0.4;
  const pillWidth = Math.min(label.length * fontSize * 0.56 + padX * 2, width - 8);
  const pillHeight = fontSize + padY * 2;
  const cx = x + width / 2;
  const top = y + Math.min(height * 0.06, 18);

  return (
    <g pointerEvents="none" style={{ userSelect: "none" }}>
      <rect
        x={cx - pillWidth / 2}
        y={top}
        width={pillWidth}
        height={pillHeight}
        rx={pillHeight / 2}
        fill="#ffffff"
        opacity={0.6}
      />
      <text
        x={cx}
        y={top + pillHeight / 2}
        textAnchor="middle"
        dominantBaseline="central"
        fontSize={fontSize}
        fontWeight={600}
        fill="#0f172a"
      >
        {label}
      </text>
    </g>
  );
}

export const ContainerLabel = memo(ContainerLabelImpl);
