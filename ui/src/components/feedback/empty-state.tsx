import { Show, type JSX } from 'solid-js'
import { Inbox } from 'lucide-solid'
import { cn } from '@/lib/utils/cn'
import { Button } from '@/components/ui/button'

export function EmptyState(props: {
  title: string
  description?: string
  icon?: JSX.Element
  action?: JSX.Element
  class?: string
}) {
  return (
    <div
      class={cn(
        'flex flex-col items-center justify-center gap-2 py-16 text-center',
        props.class,
      )}
    >
      <div class="mb-1 text-dimmed">
        {props.icon ?? <Inbox class="size-10" stroke-width={1.5} />}
      </div>
      <h3 class="text-sm font-semibold text-highlighted">{props.title}</h3>
      <Show when={props.description}>
        <p class="max-w-sm text-sm text-muted">{props.description}</p>
      </Show>
      <Show when={props.action}>
        <div class="mt-3">{props.action}</div>
      </Show>
    </div>
  )
}

export function ErrorBoundaryFallback(props: {
  error: Error
  reset: () => void
}) {
  return (
    <div class="flex flex-col items-center justify-center gap-3 p-8 text-center">
      <h2 class="text-sm font-semibold text-error">Something went wrong</h2>
      <p class="max-w-md text-sm text-muted">{props.error.message}</p>
      <Button size="sm" variant="outline" color="neutral" onClick={props.reset}>
        Try again
      </Button>
    </div>
  )
}

export function AsyncBoundary(props: {
  status: 'idle' | 'loading' | 'success' | 'error'
  error?: string | null
  loading?: JSX.Element
  onRetry?: () => void
  children: JSX.Element
  empty?: boolean
  emptyFallback?: JSX.Element
}) {
  return (
    <Show
      when={props.status !== 'loading' || props.children}
      fallback={props.loading ?? <div class="p-8 text-center text-sm text-muted">Loading…</div>}
    >
      <Show when={props.status === 'error'}>
        <div class="flex flex-col items-center gap-2 p-8 text-center">
          <p class="text-sm text-error">{props.error ?? 'Failed to load'}</p>
          <Show when={props.onRetry}>
            <Button size="sm" variant="outline" color="neutral" onClick={props.onRetry}>
              Retry
            </Button>
          </Show>
        </div>
      </Show>
      <Show when={props.status !== 'error'}>
        <Show when={!props.empty} fallback={props.emptyFallback}>
          {props.children}
        </Show>
      </Show>
    </Show>
  )
}
