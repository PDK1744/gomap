import type { Bounds, BranchElement } from "../types";

/** Bounding box that encloses every element, used to fit the map in the view. */
export function computeBounds(elements: BranchElement[]): Bounds {
  if (elements.length === 0) {
    return { minX: 0, minY: 0, maxX: 0, maxY: 0, width: 0, height: 0 };
  }

  let minX = Infinity;
  let minY = Infinity;
  let maxX = -Infinity;
  let maxY = -Infinity;

  for (const el of elements) {
    minX = Math.min(minX, el.x);
    minY = Math.min(minY, el.y);
    maxX = Math.max(maxX, el.x + el.width);
    maxY = Math.max(maxY, el.y + el.height);
  }

  return { minX, minY, maxX, maxY, width: maxX - minX, height: maxY - minY };
}
