import type { ElementType } from "../types";

/** How an element's name is drawn (if at all). */
export type LabelMode =
  | "none" // no name (doors, hallways, counters)
  | "inline" // name sits inside the cell, bottom-anchored, leaving room on top
  | "overlay"; // name floats on a top layer above any children (rooms)

interface ElementStyle {
  fill: string;
  stroke: string;
  strokeWidth: number;
  label: LabelMode;
  /** Stacking order — larger containers drawn first, small fixtures on top. */
  z: number;
}

const STYLES: Record<ElementType, ElementStyle> = {
  branch: { fill: "transparent", stroke: "#334155", strokeWidth: 6, label: "none", z: 0 },
  room: { fill: "#f1f5f9", stroke: "#94a3b8", strokeWidth: 2, label: "overlay", z: 1 },
  hallway: { fill: "#e2e8f0", stroke: "#cbd5e1", strokeWidth: 1, label: "none", z: 1 },
  office: { fill: "#dbeafe", stroke: "#60a5fa", strokeWidth: 2, label: "inline", z: 2 },
  cubicle: { fill: "#dcfce7", stroke: "#4ade80", strokeWidth: 1.5, label: "inline", z: 3 },
  counter: { fill: "#fef9c3", stroke: "#facc15", strokeWidth: 1.5, label: "none", z: 3 },
  door: { fill: "#fed7aa", stroke: "#fb923c", strokeWidth: 1.5, label: "none", z: 4 },
};

export function styleFor(type: ElementType): ElementStyle {
  return STYLES[type] ?? STYLES.room;
}

export const ELEMENT_TYPES: ElementType[] = [
  "branch",
  "room",
  "office",
  "cubicle",
  "counter",
  "hallway",
  "door",
];

export const TYPE_LABELS: Record<ElementType, string> = {
  branch: "Branch",
  room: "Room",
  office: "Office",
  cubicle: "Cubicle",
  counter: "Counter",
  hallway: "Hallway",
  door: "Door",
};

/** Strips the redundant branch-name prefix from an element's display name. */
export function displayName(name: string): string {
  return name.replace(/^Stoneledge\s+/i, "");
}

/**
 * Whether a space can hold assets. The backend's `assignable` flag is the
 * source of truth (e.g. conference rooms are assignable, restrooms are not).
 */
export function canHoldAssets(element: { assignable: boolean }): boolean {
  return element.assignable;
}

/**
 * Small leaf cells always show their asset card when zoomed in (including an
 * "Empty" placeholder). Larger containers only show a card once they actually
 * have assets, so the floor plan doesn't fill with placeholder text.
 */
export const LEAF_CARD_TYPES: ReadonlySet<ElementType> = new Set<ElementType>([
  "cubicle",
  "office",
  "counter",
]);
