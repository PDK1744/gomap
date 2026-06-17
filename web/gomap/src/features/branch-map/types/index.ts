export type ElementType =
  | "branch"
  | "door"
  | "hallway"
  | "room"
  | "office"
  | "cubicle"
  | "counter";

export interface BranchElement {
  tru_id: string;
  name: string;
  type: ElementType;
  assignable: boolean;
  /** True when the space is backed by generator power. */
  generator?: boolean;
  x: number;
  y: number;
  width: number;
  height: number;
}

export interface Branch {
  branch_name: string;
  elements: BranchElement[];
}

export type AssetKind = "pc" | "printer" | "phone";

export interface Asset {
  id: string;
  kind: AssetKind;
  name: string; // e.g. "PC500", "PR5001", "Ext. 1021"
}

/** element tru_id -> list of asset ids assigned to that space */
export type Assignments = Record<string, string[]>;

export interface Bounds {
  minX: number;
  minY: number;
  maxX: number;
  maxY: number;
  width: number;
  height: number;
}
