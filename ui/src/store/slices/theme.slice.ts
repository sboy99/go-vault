import { applyTheme } from '@/theme/apply-theme'
import { applyColorModeClass } from '@/theme/color-mode'
import type { SetStoreFunction } from 'solid-js/store'
import { logAction } from '@/store/create-store'

/** Fixed Linear-inspired look — no multi-theme UI. */
const FIXED_PRIMARY = 'indigo' as const
const FIXED_NEUTRAL = 'linear' as const
const FIXED_RADIUS = 0.125

export interface ThemeState {
  resolvedMode: 'dark'
  primary: typeof FIXED_PRIMARY
  neutral: typeof FIXED_NEUTRAL
  radius: number
}

export const initialThemeState: ThemeState = {
  resolvedMode: 'dark',
  primary: FIXED_PRIMARY,
  neutral: FIXED_NEUTRAL,
  radius: FIXED_RADIUS,
}

export function createThemeActions(
  _get: () => ThemeState,
  set: SetStoreFunction<ThemeState>,
) {
  return {
    init() {
      logAction('theme', 'init')
      applyColorModeClass('dark')
      set('resolvedMode', 'dark')
      applyTheme({
        primary: FIXED_PRIMARY,
        neutral: FIXED_NEUTRAL,
        radius: FIXED_RADIUS,
      })
    },
  }
}

export type ThemeActions = ReturnType<typeof createThemeActions>
