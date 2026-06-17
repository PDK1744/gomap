import { ELEMENT_TYPES, TYPE_LABELS, styleFor } from "../utils/style";
import { GeneratorIcon } from "./GeneratorBadge";

export function Legend() {
  return (
    <div className="rounded-lg border border-slate-200 bg-white/90 p-3 shadow-sm backdrop-blur">
      <h3 className="mb-2 text-xs font-semibold uppercase tracking-wide text-slate-500">
        Legend
      </h3>
      <ul className="grid grid-cols-2 gap-x-4 gap-y-1.5">
        {ELEMENT_TYPES.map((type) => {
          const s = styleFor(type);
          return (
            <li key={type} className="flex items-center gap-2 text-sm text-slate-700">
              <span
                className="inline-block h-3.5 w-3.5 rounded-sm border"
                style={{
                  backgroundColor: s.fill === "transparent" ? "#ffffff" : s.fill,
                  borderColor: s.stroke,
                }}
              />
              {TYPE_LABELS[type]}
            </li>
          );
        })}
      </ul>
      <div className="mt-2 flex items-center gap-2 border-t border-slate-200 pt-2 text-sm text-slate-700">
        <GeneratorIcon />
        Generator power
      </div>
    </div>
  );
}
