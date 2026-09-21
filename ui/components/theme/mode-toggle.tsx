"use client";

import { useEffect, useState } from "react";
import { Moon, Sun } from "lucide-react";
import {
  applyColorMode,
  DEFAULT_MODE,
  isColorMode,
  MODE_STORAGE_KEY,
  type ColorMode,
} from "@/lib/theme/palettes";

function readStoredMode(): ColorMode {
  const fromDom = document.documentElement.dataset.theme;
  if (fromDom && isColorMode(fromDom)) return fromDom;
  const stored = localStorage.getItem(MODE_STORAGE_KEY);
  if (stored && isColorMode(stored)) return stored;
  return DEFAULT_MODE;
}

export function ModeToggle() {
  const [mode, setMode] = useState<ColorMode>(DEFAULT_MODE);
  const [ready, setReady] = useState(false);

  useEffect(() => {
    const next = readStoredMode();
    setMode(next);
    applyColorMode(next);
    setReady(true);
  }, []);

  const toggle = () => {
    const next: ColorMode = mode === "dark" ? "light" : "dark";
    setMode(next);
    localStorage.setItem(MODE_STORAGE_KEY, next);
    applyColorMode(next);
  };

  const toLight = mode === "dark";

  return (
    <button
      type="button"
      aria-label={
        ready
          ? toLight
            ? "Switch to light mode"
            : "Switch to dark mode"
          : "Toggle color mode"
      }
      onClick={toggle}
      className="inline-flex size-8 items-center justify-center rounded-none border border-dashed border-edge bg-transparent text-foreground transition-colors duration-200 ease-out hover:text-brand"
    >
      {!ready ? (
        <span className="size-4" aria-hidden />
      ) : toLight ? (
        <Sun className="size-4" strokeWidth={1.75} />
      ) : (
        <Moon className="size-4" strokeWidth={1.75} />
      )}
    </button>
  );
}
