import { useCallback, useRef, useState } from "react";

export interface Transform {
  scale: number;
  tx: number;
  ty: number;
}

const MIN_SCALE = 0.05;
const MAX_SCALE = 8;
const DRAG_THRESHOLD_PX = 4; // movement below this is treated as a click

function clampScale(scale: number): number {
  return Math.min(MAX_SCALE, Math.max(MIN_SCALE, scale));
}

/**
 * Pan/zoom controller for an SVG <g> transform. Zoom is anchored at the pointer
 * so the point under the cursor stays put while scaling.
 *
 * Pointer capture starts only after the pointer moves past a small threshold —
 * capturing on pointerdown would retarget the click event to the SVG and steal
 * clicks from the map elements. Use `wasDrag()` in click handlers to ignore the
 * click that fires at the end of a pan.
 */
export function usePanZoom(initial: Transform) {
  const [transform, setTransform] = useState<Transform>(initial);
  const drag = useRef<{ x: number; y: number; moved: boolean } | null>(null);
  const didDrag = useRef(false);
  const [isPanning, setIsPanning] = useState(false);

  const onWheel = useCallback((e: React.WheelEvent<SVGSVGElement>) => {
    const rect = e.currentTarget.getBoundingClientRect();
    const px = e.clientX - rect.left;
    const py = e.clientY - rect.top;
    const factor = Math.exp(-e.deltaY * 0.0015);

    setTransform((t) => {
      const nextScale = clampScale(t.scale * factor);
      const ratio = nextScale / t.scale;
      return {
        scale: nextScale,
        tx: px - (px - t.tx) * ratio,
        ty: py - (py - t.ty) * ratio,
      };
    });
  }, []);

  const onPointerDown = useCallback((e: React.PointerEvent<SVGSVGElement>) => {
    if (e.button !== 0) return;
    drag.current = { x: e.clientX, y: e.clientY, moved: false };
    didDrag.current = false;
  }, []);

  const onPointerMove = useCallback((e: React.PointerEvent<SVGSVGElement>) => {
    const d = drag.current;
    if (!d) return;

    const dx = e.clientX - d.x;
    const dy = e.clientY - d.y;

    if (!d.moved) {
      if (Math.hypot(dx, dy) < DRAG_THRESHOLD_PX) return;
      d.moved = true;
      setIsPanning(true);
      e.currentTarget.setPointerCapture(e.pointerId);
    }

    d.x = e.clientX;
    d.y = e.clientY;
    setTransform((t) => ({ ...t, tx: t.tx + dx, ty: t.ty + dy }));
  }, []);

  const endDrag = useCallback((e: React.PointerEvent<SVGSVGElement>) => {
    if (e.currentTarget.hasPointerCapture(e.pointerId)) {
      e.currentTarget.releasePointerCapture(e.pointerId);
    }
    if (drag.current?.moved) didDrag.current = true;
    drag.current = null;
    setIsPanning(false);
  }, []);

  /** True when the click being handled is the tail end of a pan gesture. */
  const wasDrag = useCallback(() => didDrag.current, []);

  const reset = useCallback(() => setTransform(initial), [initial]);

  return {
    transform,
    setTransform,
    isPanning,
    wasDrag,
    reset,
    handlers: {
      onWheel,
      onPointerDown,
      onPointerMove,
      onPointerUp: endDrag,
      onPointerLeave: endDrag,
    },
  };
}
