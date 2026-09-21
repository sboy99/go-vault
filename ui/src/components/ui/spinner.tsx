import { cn } from '@/lib/utils/cn'
import type { JSX } from 'solid-js'

export function Spinner(props: { class?: string; size?: 'sm' | 'md' | 'lg' }) {
  const size = () => props.size ?? 'md'
  return (
    <span
      role="status"
      aria-label="Loading"
      class={cn(
        'inline-block animate-spin rounded-full border-2 border-current border-r-transparent text-muted',
        size() === 'sm' && 'size-3.5',
        size() === 'md' && 'size-5',
        size() === 'lg' && 'size-8',
        props.class,
      )}
    />
  )
}

export function Skeleton(props: { class?: string } & JSX.HTMLAttributes<HTMLDivElement>) {
  return (
    <div
      {...props}
      class={cn('animate-pulse rounded-md bg-elevated', props.class)}
    />
  )
}

export function Kbd(props: { children: string; class?: string }) {
  return (
    <kbd
      class={cn(
        'inline-flex items-center rounded border border-default bg-elevated px-1.5 py-0.5',
        'font-mono text-[10px] text-muted',
        props.class,
      )}
    >
      {props.children}
    </kbd>
  )
}
