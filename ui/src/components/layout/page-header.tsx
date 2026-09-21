import type { JSX } from 'solid-js'
import { Show } from 'solid-js'
import { cn } from '@/lib/utils/cn'

export function PageHeader(props: {
  title: string
  description?: string
  actions?: JSX.Element
  icon?: JSX.Element
  class?: string
}) {
  return (
    <div
      class={cn(
        'flex h-[var(--ui-header-height)] shrink-0 items-center justify-between gap-2 border-b border-default px-3',
        props.class,
      )}
    >
      <div class="flex min-w-0 items-center gap-2">
        <Show when={props.icon}>
          <span class="text-primary">{props.icon}</span>
        </Show>
        <div class="min-w-0">
          <h1 class="truncate text-[13px] font-semibold text-highlighted">{props.title}</h1>
          <Show when={props.description}>
            <p class="truncate text-[12px] text-muted">{props.description}</p>
          </Show>
        </div>
      </div>
      <Show when={props.actions}>
        <div class="flex shrink-0 items-center gap-2">{props.actions}</div>
      </Show>
    </div>
  )
}

export function AppShell(props: { children: JSX.Element }) {
  return (
    <div class="flex h-full min-h-0 bg-default text-default">
      {/* sidebar injected by App */}
      {props.children}
    </div>
  )
}

export function MainPanel(props: { children: JSX.Element; class?: string }) {
  return (
    <main class={cn('flex min-w-0 flex-1 flex-col overflow-hidden', props.class)}>
      {props.children}
    </main>
  )
}
