import {
  NEUTRAL_PALETTES,
  PRIMARY_PALETTES,
  type NeutralPaletteId,
  type PrimaryPaletteId,
} from '@/theme/palettes'

const STEPS = ['50', '100', '200', '300', '400', '500', '600', '700', '800', '900', '950'] as const

export function applyPrimaryPalette(id: PrimaryPaletteId): void {
  const ramp = PRIMARY_PALETTES[id]
  const root = document.documentElement
  for (const step of STEPS) {
    root.style.setProperty(`--ui-color-primary-${step}`, ramp[step])
  }
  // Linear-inspired: keep accent at 500 in both modes
  root.style.setProperty('--ui-primary', ramp['500'])
  root.style.setProperty('--ui-primary-hover', ramp['400'])
}

export function applyNeutralPalette(id: NeutralPaletteId): void {
  const ramp = NEUTRAL_PALETTES[id]
  const root = document.documentElement
  for (const step of STEPS) {
    root.style.setProperty(`--ui-color-neutral-${step}`, ramp[step])
  }
}

export function applyRadius(rem: number): void {
  document.documentElement.style.setProperty('--ui-radius', `${rem}rem`)
}

export function applyTheme(options: {
  primary: PrimaryPaletteId
  neutral: NeutralPaletteId
  radius: number
}): void {
  applyPrimaryPalette(options.primary)
  applyNeutralPalette(options.neutral)
  applyRadius(options.radius)
}
