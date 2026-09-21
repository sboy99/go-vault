import { For, Show } from 'solid-js'
import { CheckCircle2, Info, X, XCircle } from 'lucide-solid'
import { useStore } from '@/store/root'
import { cn } from '@/lib/utils/cn'
import type { ToastKind } from '@/store/slices/toast.slice'

const icons: Record<ToastKind, typeof Info> = {
  success: CheckCircle2,
  error: XCircle,
  info: Info,
}

const colors: Record<ToastKind, string> = {
  success: 'border-success/30 bg-default',
  error: 'border-error/30 bg-default',
  info: 'border-info/30 bg-default',
}

export function Toaster() {
  const { state, actions } = useStore()
  return (
    <div class="pointer-events-none fixed right-4 bottom-4 z-[100] flex w-80 flex-col gap-2">
      <For each={state.toast.items}>
        {(toast) => {
          const Icon = icons[toast.kind]
          return (
            <div
              class={cn(
                'pointer-events-auto flex gap-2 rounded-lg border p-3 shadow-lg',
                colors[toast.kind],
              )}
              role="status"
            >
              <Icon
                class={cn(
                  'mt-0.5 size-4 shrink-0',
                  toast.kind === 'success' && 'text-success',
                  toast.kind === 'error' && 'text-error',
                  toast.kind === 'info' && 'text-info',
                )}
              />
              <div class="min-w-0 flex-1">
                <p class="text-sm font-medium text-highlighted">{toast.title}</p>
                <Show when={toast.description}>
                  <p class="mt-0.5 text-xs text-muted">{toast.description}</p>
                </Show>
              </div>
              <button
                type="button"
                class="text-muted hover:text-highlighted"
                aria-label="Dismiss"
                onClick={() => actions.toast.dismiss(toast.id)}
              >
                <X class="size-3.5" />
              </button>
            </div>
          )
        }}
      </For>
    </div>
  )
}
