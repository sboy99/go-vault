import { describe, expect, it } from 'vitest'
import { applyTheme } from '@/theme/apply-theme'
import { cycleColorMode, resolveColorMode } from '@/theme/color-mode'

describe('applyTheme', () => {
  it('writes CSS variables on :root', () => {
    applyTheme({ primary: 'indigo', neutral: 'linear', radius: 0.125 })
    const root = document.documentElement
    expect(root.style.getPropertyValue('--ui-radius')).toBe('0.125rem')
    expect(root.style.getPropertyValue('--ui-color-primary-500')).toBe('#5e6ad2')
    expect(root.style.getPropertyValue('--ui-color-neutral-950')).toBe('#08090a')
    expect(root.style.getPropertyValue('--ui-primary')).toBe('#5e6ad2')
  })
})

describe('color-mode', () => {
  it('cycles modes', () => {
    expect(cycleColorMode('system')).toBe('light')
    expect(cycleColorMode('light')).toBe('dark')
    expect(cycleColorMode('dark')).toBe('system')
  })

  it('resolves light/dark directly', () => {
    expect(resolveColorMode('light')).toBe('light')
    expect(resolveColorMode('dark')).toBe('dark')
  })
})
