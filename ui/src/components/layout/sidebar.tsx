import { A, useLocation, useNavigate } from '@solidjs/router'
import { For, createSignal } from 'solid-js'
import {
  HardDrive,
  LayoutDashboard,
  ListTodo,
  Plus,
  Settings,
} from 'lucide-solid'
import { useStore } from '@/store/root'
import { Button } from '@/components/ui/button'
import { cn } from '@/lib/utils/cn'

const NAV = [
  { href: '/', label: 'Overview', icon: LayoutDashboard },
  { href: '/backups', label: 'Backups', icon: HardDrive },
  { href: '/jobs', label: 'Jobs', icon: ListTodo },
  { href: '/settings', label: 'Settings', icon: Settings },
] as const

export function Sidebar() {
  const { state, actions, selectors } = useStore()
  const location = useLocation()
  const navigate = useNavigate()
  const [starting, setStarting] = createSignal(false)

  async function onCreate() {
    setStarting(true)
    try {
      const id = await actions.backups.createBackup()
      if (id) {
        await actions.jobs.fetchList()
        navigate('/jobs')
      }
    } finally {
      setStarting(false)
    }
  }

  return (
    <aside class="flex w-52 shrink-0 flex-col border-r border-default bg-muted">
      <div class="border-b border-default p-3">
        <Button
          size="sm"
          class="w-full justify-center"
          loading={starting() || state.backups.creating}
          disabled={selectors.hasActiveJob()}
          onClick={() => void onCreate()}
        >
          <Plus class="size-3.5" />
          Create Backup
        </Button>
      </div>

      <nav class="flex flex-1 flex-col gap-0.5 p-2">
        <p class="mb-1 px-2 text-[10px] font-semibold uppercase tracking-wider text-dimmed">
          Platform
        </p>
        <For each={NAV}>
          {(item) => {
            const active = () =>
              item.href === '/'
                ? location.pathname === '/'
                : location.pathname.startsWith(item.href)
            const Icon = item.icon
            return (
              <A
                href={item.href}
                class={cn(
                  'flex items-center gap-2 rounded-md px-2 py-1.5 text-[13px] font-medium transition-colors',
                  active()
                    ? 'bg-elevated text-highlighted'
                    : 'text-toned hover:bg-elevated hover:text-highlighted',
                )}
              >
                <Icon class={cn('size-3.5', active() ? 'text-primary' : 'text-muted')} />
                {item.label}
              </A>
            )
          }}
        </For>
      </nav>
    </aside>
  )
}
