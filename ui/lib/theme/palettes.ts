import { applyFavicon } from "@/lib/theme/favicon";

export type PrimaryId =
  | "black"
  | "red"
  | "orange"
  | "amber"
  | "yellow"
  | "lime"
  | "green"
  | "emerald"
  | "teal"
  | "cyan"
  | "sky"
  | "blue"
  | "indigo"
  | "violet"
  | "purple"
  | "fuchsia"
  | "pink"
  | "rose";

export type NeutralId =
  | "slate"
  | "gray"
  | "zinc"
  | "neutral"
  | "stone"
  | "taupe"
  | "mauve"
  | "mist"
  | "olive";

export const PRIMARY_OPTIONS: { id: PrimaryId; label: string; swatch: string }[] =
  [
    { id: "black", label: "Black", swatch: "#a1a1aa" },
    { id: "red", label: "Red", swatch: "#ef4444" },
    { id: "orange", label: "Orange", swatch: "#f97316" },
    { id: "amber", label: "Amber", swatch: "#f59e0b" },
    { id: "yellow", label: "Yellow", swatch: "#eab308" },
    { id: "lime", label: "Lime", swatch: "#84cc16" },
    { id: "green", label: "Green", swatch: "#22c55e" },
    { id: "emerald", label: "Emerald", swatch: "#10b981" },
    { id: "teal", label: "Teal", swatch: "#14b8a6" },
    { id: "cyan", label: "Cyan", swatch: "#06b6d4" },
    { id: "sky", label: "Sky", swatch: "#0ea5e9" },
    { id: "blue", label: "Blue", swatch: "#3b82f6" },
    { id: "indigo", label: "Indigo", swatch: "#6366f1" },
    { id: "violet", label: "Violet", swatch: "#8b5cf6" },
    { id: "purple", label: "Purple", swatch: "#a855f7" },
    { id: "fuchsia", label: "Fuchsia", swatch: "#d946ef" },
    { id: "pink", label: "Pink", swatch: "#ec4899" },
    { id: "rose", label: "Rose", swatch: "#f43f5e" },
  ];

export const NEUTRAL_OPTIONS: {
  id: NeutralId;
  label: string;
  swatch: string;
}[] = [
  { id: "slate", label: "Slate", swatch: "#64748b" },
  { id: "gray", label: "Gray", swatch: "#6b7280" },
  { id: "zinc", label: "Zinc", swatch: "#71717a" },
  { id: "neutral", label: "Neutral", swatch: "#737373" },
  { id: "stone", label: "Stone", swatch: "#78716c" },
  { id: "taupe", label: "Taupe", swatch: "#8a7f72" },
  { id: "mauve", label: "Mauve", swatch: "#8b7d8e" },
  { id: "mist", label: "Mist", swatch: "#7a8a99" },
  { id: "olive", label: "Olive", swatch: "#7d8468" },
];

export type ColorMode = "dark" | "light";

export const DEFAULT_PRIMARY: PrimaryId = "emerald";
export const DEFAULT_NEUTRAL: NeutralId = "neutral";
export const DEFAULT_MODE: ColorMode = "dark";

export const PRIMARY_STORAGE_KEY = "govault-primary";
export const NEUTRAL_STORAGE_KEY = "govault-neutral";
export const MODE_STORAGE_KEY = "govault-theme";

export function isPrimaryId(value: string): value is PrimaryId {
  return PRIMARY_OPTIONS.some((o) => o.id === value);
}

export function isNeutralId(value: string): value is NeutralId {
  return NEUTRAL_OPTIONS.some((o) => o.id === value);
}

export function isColorMode(value: string): value is ColorMode {
  return value === "dark" || value === "light";
}

export function applyColorMode(mode: ColorMode) {
  const root = document.documentElement;
  root.dataset.theme = mode;
  root.classList.toggle("dark", mode === "dark");
}

export function applyTheme(primary: PrimaryId, neutral: NeutralId) {
  const root = document.documentElement;
  root.dataset.primary = primary;
  root.dataset.neutral = neutral;
  const swatch =
    PRIMARY_OPTIONS.find((o) => o.id === primary)?.swatch ?? "#10b981";
  applyFavicon(swatch);
}
