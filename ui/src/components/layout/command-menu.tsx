import { For, Show, createMemo, createSignal, onCleanup, onMount } from 'solid-js'
import { useNavigate } from '@solidjs/router'
import {
  HardDrive,
  LayoutDashboard,
  ListTodo,
  Search,
  Settings,
} from 'lucide-solid'
import { Dialog } from '@kobalte/core/dialog'
import { cn } from '@/lib/utils/cn'
import { Kbd } from '@/components/ui/spinner'

const COMMANDS = [
  { id: 'overview', label: 'Overview', href: '/', icon: LayoutDashboard, keywords: 'home dashboard' },
  { id: 'backups', label: 'Backups', href: '/backups', icon: HardDrive, keywords: 'backup dump' },
  { id: 'jobs', label: 'Jobs', href: '/jobs', icon: ListTodo, keywords: 'job run task' },
  { id: 'settings', label: 'Settings', href: '/settings', icon: Settings, keywords: 'token config' },
] as const

const OPEN_EVENT = 'go-vault:open-command-menu'

export function openCommandMenu() {
  window.dispatchEvent(new CustomEvent(OPEN_EVENT))
}

export function CommandMenu() {
  const navigate = useNavigate()
  const [open, setOpen] = createSignal(false)
  const [query, setQuery] = createSignal('')
  const [active, setActive] = createSignal(0)

  const filtered = createMemo(() => {
    const q = query().trim().toLowerCase()
    if (!q) return [...COMMANDS]
    return COMMANDS.filter(
      (c) =>
        c.label.toLowerCase().includes(q) ||
        c.keywords.includes(q) ||
        c.href.includes(q),
    )
  })

  function go(href: string) {
    setOpen(false)
    setQuery('')
    setActive(0)
    navigate(href)
  }

  function onKeyDown(e: KeyboardEvent) {
    const meta = e.metaKey || e.ctrlKey
    if (meta && e.key.toLowerCase() === 'k') {
      e.preventDefault()
      setOpen((v) => !v)
      return
    }
    if (!open()) return
    if (e.key === 'Escape') {
      setOpen(false)
      return
    }
    const list = filtered()
    if (e.key === 'ArrowDown') {
      e.preventDefault()
      setActive((i) => (i + 1) % Math.max(list.length, 1))
    } else if (e.key === 'ArrowUp') {
      e.preventDefault()
      setActive((i) => (i - 1 + Math.max(list.length, 1)) % Math.max(list.length, 1))
    } else if (e.key === 'Enter') {
      e.preventDefault()
      const item = list[active()]
      if (item) go(item.href)
    }
  }

  onMount(() => {
    const onOpen = () => setOpen(true)
    window.addEventListener('keydown', onKeyDown)
    window.addEventListener(OPEN_EVENT, onOpen)
    onCleanup(() => {
      window.removeEventListener('keydown', onKeyDown)
      window.removeEventListener(OPEN_EVENT, onOpen)
    })
  })

  return (
    <Dialog
      open={open()}
      onOpenChange={(v) => {
        setOpen(v)
        if (!v) {
          setQuery('')
          setActive(0)
        }
      }}
    >
      <Dialog.Portal>
        <Dialog.Overlay class="fixed inset-0 z-50 bg-black/50" />
        <div class="fixed inset-0 z-50 flex items-start justify-center px-4 pt-[15vh]">
          <Dialog.Content class="w-full max-w-md overflow-hidden rounded-md border border-default bg-muted outline-none">
            <div class="flex items-center gap-2 border-b border-default px-3">
              <Search class="size-3.5 shrink-0 text-muted" />
              <input
                class="h-10 w-full bg-transparent text-[13px] text-highlighted placeholder:text-dimmed outline-none"
                placeholder="Jump to…"
                value={query()}
                onInput={(e) => {
                  setQuery(e.currentTarget.value)
                  setActive(0)
                }}
                autofocus
              />
              <Kbd>esc</Kbd>
            </div>
            <div class="max-h-64 overflow-y-auto p-1">
              <Show
                when={filtered().length > 0}
                fallback={
                  <p class="px-2 py-6 text-center text-[12px] text-muted">No results</p>
                }
              >
                <For each={filtered()}>
                  {(item, index) => {
                    const Icon = item.icon
                    return (
                      <button
                        type="button"
                        class={cn(
                          'flex w-full items-center gap-2 rounded-md px-2 py-1.5 text-left text-[13px]',
                          index() === active()
                            ? 'bg-elevated text-highlighted'
                            : 'text-default hover:bg-elevated',
                        )}
                        onMouseEnter={() => setActive(index())}
                        onClick={() => go(item.href)}
                      >
                        <Icon class="size-3.5 text-primary" />
                        <span class="flex-1 font-medium">{item.label}</span>
                        <span class="text-[11px] text-dimmed">{item.href}</span>
                      </button>
                    )
                  }}
                </For>
              </Show>
            </div>
            <div class="flex items-center gap-2 border-t border-default px-3 py-1.5 text-[11px] text-muted">
              <span class="inline-flex items-center gap-1">
                <Kbd>↑</Kbd>
                <Kbd>↓</Kbd>
                navigate
              </span>
              <span class="inline-flex items-center gap-1">
                <Kbd>↵</Kbd>
                open
              </span>
              <span class="ml-auto inline-flex items-center gap-1">
                <Kbd>⌘</Kbd>
                <Kbd>K</Kbd>
              </span>
            </div>
          </Dialog.Content>
        </div>
      </Dialog.Portal>
    </Dialog>
  )
}
