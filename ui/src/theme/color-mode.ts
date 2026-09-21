import { STORAGE_KEYS } from '@/config/constants'
import { readStorage, writeStorage } from '@/lib/utils/storage'
import type { ColorMode } from '@/theme/palettes'

export type { ColorMode }

function systemPrefersDark(): boolean {
  return window.matchMedia('(prefers-color-scheme: dark)').matches
}

export function resolveColorMode(mode: ColorMode): 'light' | 'dark' {
  if (mode === 'system') return systemPrefersDark() ? 'dark' : 'light'
  return mode
}

export function applyColorModeClass(mode: ColorMode): 'light' | 'dark' {
  const resolved = resolveColorMode(mode)
  const root = document.documentElement
  root.classList.remove('light', 'dark')
  root.classList.add(resolved)
  return resolved
}

export function readStoredColorMode(): ColorMode {
  const raw = readStorage(STORAGE_KEYS.colorMode)
  if (raw === 'light' || raw === 'dark' || raw === 'system') return raw
  return 'system'
}

export function persistColorMode(mode: ColorMode): void {
  writeStorage(STORAGE_KEYS.colorMode, mode)
}

export function cycleColorMode(current: ColorMode): ColorMode {
  if (current === 'system') return 'light'
  if (current === 'light') return 'dark'
  return 'system'
}
