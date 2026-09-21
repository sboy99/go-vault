"use client";

import { useEffect, useRef, useState } from "react";
import { Palette } from "lucide-react";
import { cn } from "@/lib/cn";
import {
  applyTheme,
  DEFAULT_NEUTRAL,
  DEFAULT_PRIMARY,
  isNeutralId,
  isPrimaryId,
  NEUTRAL_OPTIONS,
  NEUTRAL_STORAGE_KEY,
  PRIMARY_OPTIONS,
  PRIMARY_STORAGE_KEY,
  type NeutralId,
  type PrimaryId,
} from "@/lib/theme/palettes";

export function ThemePicker() {
  const [open, setOpen] = useState(false);
  const [primary, setPrimary] = useState<PrimaryId>(DEFAULT_PRIMARY);
  const [neutral, setNeutral] = useState<NeutralId>(DEFAULT_NEUTRAL);
  const panelRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const storedPrimary = localStorage.getItem(PRIMARY_STORAGE_KEY);
    const storedNeutral = localStorage.getItem(NEUTRAL_STORAGE_KEY);
    const nextPrimary =
      storedPrimary && isPrimaryId(storedPrimary)
        ? storedPrimary
        : DEFAULT_PRIMARY;
    const nextNeutral =
      storedNeutral && isNeutralId(storedNeutral)
        ? storedNeutral
        : DEFAULT_NEUTRAL;
    setPrimary(nextPrimary);
    setNeutral(nextNeutral);
    applyTheme(nextPrimary, nextNeutral);
  }, []);

  useEffect(() => {
    if (!open) return;
    const onPointer = (e: MouseEvent) => {
      if (!panelRef.current?.contains(e.target as Node)) setOpen(false);
    };
    const onKey = (e: KeyboardEvent) => {
      if (e.key === "Escape") setOpen(false);
    };
    document.addEventListener("mousedown", onPointer);
    document.addEventListener("keydown", onKey);
    return () => {
      document.removeEventListener("mousedown", onPointer);
      document.removeEventListener("keydown", onKey);
    };
  }, [open]);

  const selectPrimary = (id: PrimaryId) => {
    setPrimary(id);
    localStorage.setItem(PRIMARY_STORAGE_KEY, id);
    applyTheme(id, neutral);
  };

  const selectNeutral = (id: NeutralId) => {
    setNeutral(id);
    localStorage.setItem(NEUTRAL_STORAGE_KEY, id);
    applyTheme(primary, id);
  };

  return (
    <div className="relative" ref={panelRef}>
      <button
        type="button"
        aria-label="Theme"
        aria-expanded={open}
        onClick={() => setOpen((v) => !v)}
        className="inline-flex size-8 items-center justify-center rounded-none border border-dashed border-edge bg-transparent text-foreground transition-colors duration-200 ease-out hover:text-brand"
      >
        <Palette className="size-4" strokeWidth={1.75} />
      </button>

      {open ? (
        <div className="absolute top-full right-0 z-30 mt-1.5 w-[17.5rem] rounded-none border border-dashed border-edge bg-surface p-2">
          <p className="mb-1.5 px-0.5 text-[11px] font-medium text-foreground-muted">
            Primary
          </p>
          <div className="grid grid-cols-3 gap-1">
            {PRIMARY_OPTIONS.map((opt) => (
              <button
                key={opt.id}
                type="button"
                onClick={() => selectPrimary(opt.id)}
                className={cn(
                  "flex items-center gap-1.5 rounded border px-1.5 py-1 text-left text-[11px] transition-colors",
                  primary === opt.id
                    ? "border-control bg-surface-control text-foreground"
                    : "border-transparent text-foreground-light hover:bg-surface-control",
                )}
              >
                <span
                  aria-hidden
                  className="size-2.5 shrink-0 rounded-full"
                  style={{ backgroundColor: opt.swatch }}
                />
                {opt.label}
              </button>
            ))}
          </div>

          <p className="mt-2.5 mb-1.5 px-0.5 text-[11px] font-medium text-foreground-muted">
            Neutral
          </p>
          <div className="grid grid-cols-3 gap-1">
            {NEUTRAL_OPTIONS.map((opt) => (
              <button
                key={opt.id}
                type="button"
                onClick={() => selectNeutral(opt.id)}
                className={cn(
                  "flex items-center gap-1.5 rounded border px-1.5 py-1 text-left text-[11px] transition-colors",
                  neutral === opt.id
                    ? "border-control bg-surface-control text-foreground"
                    : "border-transparent text-foreground-light hover:bg-surface-control",
                )}
              >
                <span
                  aria-hidden
                  className="size-2.5 shrink-0 rounded-full"
                  style={{ backgroundColor: opt.swatch }}
                />
                {opt.label}
              </button>
            ))}
          </div>
        </div>
      ) : null}
    </div>
  );
}
