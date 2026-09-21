import { A } from '@solidjs/router'
import { Activity, DatabaseBackup, Search } from 'lucide-solid'
import { Show } from 'solid-js'
import { useStore } from '@/store/root'
import { APP_NAME } from '@/config/constants'
import { openCommandMenu } from '@/components/layout/command-menu'

export function AppHeader() {
  const { state } = useStore()

  return (
    <header class="flex h-[var(--ui-header-height)] shrink-0 items-center justify-between gap-2 border-b border-default bg-default px-3">
      <div class="flex min-w-0 items-center gap-2">
        <div class="flex size-6 shrink-0 items-center justify-center rounded-md bg-primary text-inverted">
          <DatabaseBackup class="size-3.5" />
        </div>
        <span class="text-[13px] font-semibold text-highlighted">{APP_NAME}</span>
      </div>

      <div class="flex items-center gap-2">
        <Show when={state.settings.unauthorized}>
          <A
            href="/settings"
            class="flex items-center gap-1.5 rounded-md bg-warning/10 px-2 py-1 text-[11px] text-warning"
          >
            <Activity class="size-3.5" />
            <span class="hidden sm:inline">API token required</span>
          </A>
        </Show>

        <button
          type="button"
          class="inline-flex h-7 items-center gap-1.5 rounded-md border border-default bg-elevated px-2 text-[11px] text-muted transition-colors hover:text-highlighted"
          onClick={() => openCommandMenu()}
          aria-label="Open command menu"
        >
          <Search class="size-3 text-primary" />
          <span class="hidden sm:inline">Search</span>
          <span class="ml-1 hidden items-center gap-0.5 sm:inline-flex">
            <kbd class="rounded border border-default bg-default px-1 font-mono text-[10px]">⌘</kbd>
            <kbd class="rounded border border-default bg-default px-1 font-mono text-[10px]">K</kbd>
          </span>
        </button>
      </div>
    </header>
  )
}
