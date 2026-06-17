import type { Asset, AssetKind, Assignments } from "../types";

/**
 * Dummy asset pool. This stands in for a list that will eventually come from the
 * backend — swap `DUMMY_ASSETS` for a fetched list and the rest keeps working.
 */
function make(kind: AssetKind, names: string[]): Asset[] {
  return names.map((name) => ({ id: name, kind, name }));
}

const PCS = make(
  "pc",
  Array.from({ length: 30 }, (_, i) => `PC${500 + i}`),
);

const PRINTERS = make(
  "printer",
  Array.from({ length: 8 }, (_, i) => `PR${5001 + i}`),
);

const PHONES = make(
  "phone",
  Array.from({ length: 40 }, (_, i) => `Ext. ${1021 + i}`),
);

export const DUMMY_ASSETS: Asset[] = [...PCS, ...PRINTERS, ...PHONES];

/** dataTransfer MIME type used when dragging an asset from the sidebar list. */
export const ASSET_DRAG_TYPE = "application/x-gomap-asset";

/** Pre-seeded assignments so the demo shows assets without manual setup. */
export const INITIAL_ASSIGNMENTS: Assignments = {
  sl_cubicle_100: ["PC500", "Ext. 1021"],
  sl_cubicle_101: ["PC501", "Ext. 1022"],
  sl_cubicle_102: ["PC502", "Ext. 1023"],
  sl_office_101: ["PC520", "Ext. 1040", "PR5001"],
  sl_office_102: ["PC521", "Ext. 1041"],
  sl_counter_1: ["PC510"],
};

export const ASSET_ICON: Record<AssetKind, string> = {
  pc: "🖥",
  printer: "🖨",
  phone: "☎",
};

export const ASSET_KIND_LABEL: Record<AssetKind, string> = {
  pc: "PCs",
  printer: "Printers",
  phone: "Extensions",
};
