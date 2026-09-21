import { splitProps, type JSX, type ParentProps } from 'solid-js'
import { cn } from '@/lib/utils/cn'

export type ButtonVariant = 'solid' | 'outline' | 'ghost' | 'soft' | 'link'
export type ButtonColor = 'primary' | 'neutral' | 'error' | 'success'
export type ButtonSize = 'xs' | 'sm' | 'md' | 'lg'

export interface ButtonProps extends ParentProps, Omit<JSX.ButtonHTMLAttributes<HTMLButtonElement>, 'color'> {
  variant?: ButtonVariant
  color?: ButtonColor
  size?: ButtonSize
  loading?: boolean
}

const sizeClass: Record<ButtonSize, string> = {
  xs: 'h-6 px-1.5 text-[11px] gap-1',
  sm: 'h-7 px-2 text-[12px] gap-1',
  md: 'h-7 px-2.5 text-[13px] gap-1.5',
  lg: 'h-8 px-3 text-[13px] gap-1.5',
}

function colorClasses(variant: ButtonVariant, color: ButtonColor): string {
  if (variant === 'link') {
    return 'bg-transparent underline-offset-4 hover:underline text-primary p-0 h-auto'
  }
  if (variant === 'ghost') {
    const map = {
      primary: 'text-primary hover:bg-elevated',
      neutral: 'text-default hover:bg-elevated',
      error: 'text-error hover:bg-elevated',
      success: 'text-success hover:bg-elevated',
    }
    return map[color]
  }
  if (variant === 'outline') {
    const map = {
      primary: 'border border-primary/40 text-primary hover:bg-primary/10',
      neutral: 'border border-default text-default hover:bg-elevated',
      error: 'border border-error/40 text-error hover:bg-error/10',
      success: 'border border-success/40 text-success hover:bg-success/10',
    }
    return map[color]
  }
  if (variant === 'soft') {
    const map = {
      primary: 'bg-primary/10 text-primary hover:bg-primary/15',
      neutral: 'bg-elevated text-default hover:bg-accented',
      error: 'bg-error/10 text-error hover:bg-error/15',
      success: 'bg-success/10 text-success hover:bg-success/15',
    }
    return map[color]
  }
  const map = {
    primary: 'bg-primary text-inverted hover:opacity-90',
    neutral: 'bg-inverted text-inverted hover:opacity-90',
    error: 'bg-error text-white hover:opacity-90',
    success: 'bg-success text-white hover:opacity-90',
  }
  return map[color]
}

export function Button(props: ButtonProps) {
  const [local, rest] = splitProps(props, [
    'variant',
    'color',
    'size',
    'loading',
    'class',
    'children',
    'disabled',
  ])
  const variant = () => local.variant ?? 'solid'
  const color = () => local.color ?? 'primary'
  const size = () => local.size ?? 'md'

  return (
    <button
      type="button"
      {...rest}
      disabled={local.disabled || local.loading}
      class={cn(
        'inline-flex items-center justify-center rounded-md font-medium transition-colors',
        'disabled:cursor-not-allowed disabled:opacity-50',
        sizeClass[size()],
        colorClasses(variant(), color()),
        local.class,
      )}
    >
      {local.loading ? (
        <span class="size-3 animate-spin rounded-full border-2 border-current border-r-transparent" />
      ) : null}
      {local.children}
    </button>
  )
}
