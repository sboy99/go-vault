import { splitProps, type JSX, type ParentProps } from 'solid-js'
import { cn } from '@/lib/utils/cn'

export interface BadgeProps extends ParentProps, JSX.HTMLAttributes<HTMLSpanElement> {
  color?: 'primary' | 'neutral' | 'success' | 'warning' | 'error' | 'info'
  variant?: 'solid' | 'soft' | 'outline'
  size?: 'sm' | 'md'
}

const soft: Record<string, string> = {
  primary: 'bg-primary/10 text-primary',
  neutral: 'bg-elevated text-muted',
  success: 'bg-success/10 text-success',
  warning: 'bg-warning/10 text-warning',
  error: 'bg-error/10 text-error',
  info: 'bg-info/10 text-info',
}

const solid: Record<string, string> = {
  primary: 'bg-primary text-inverted',
  neutral: 'bg-inverted text-inverted',
  success: 'bg-success text-white',
  warning: 'bg-warning text-black',
  error: 'bg-error text-white',
  info: 'bg-info text-white',
}

export function Badge(props: BadgeProps) {
  const [local, rest] = splitProps(props, ['class', 'color', 'variant', 'size', 'children'])
  const color = () => local.color ?? 'neutral'
  const variant = () => local.variant ?? 'soft'
  return (
    <span
      {...rest}
      class={cn(
        'inline-flex items-center rounded-md font-medium',
        'px-1.5 py-0.5 text-[11px] leading-none',
        variant() === 'solid' && solid[color()],
        variant() === 'soft' && soft[color()],
        variant() === 'outline' && `border border-current ${soft[color()]} bg-transparent`,
        local.class,
      )}
    >
      {local.children}
    </span>
  )
}
