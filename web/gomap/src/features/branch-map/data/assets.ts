import type { AssetKind } from "../types";

/** dataTransfer MIME type used when dragging an asset from the sidebar list. */
export const ASSET_DRAG_TYPE = "application/x-gomap-asset";

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
