import { For, Show, createSignal, onMount } from 'solid-js'
import { useNavigate } from '@solidjs/router'
import { Download, HardDrive, Info, RefreshCw, RotateCcw } from 'lucide-solid'
import type { BackupStatus } from '@/domain/backup'
import { useStore } from '@/store/root'
import { PageHeader } from '@/components/layout/page-header'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Input } from '@/components/ui/input'
import { Modal } from '@/components/ui/modal'
import { EmptyState } from '@/components/feedback/empty-state'
import { Spinner } from '@/components/ui/spinner'
import { isApiError } from '@/lib/http/errors'
import { formatRelativeTime } from '@/lib/utils/format'
import { cn } from '@/lib/utils/cn'

const FILTER_TABS: { value: BackupStatus | ''; label: string }[] = [
  { value: '', label: 'All' },
  { value: 'success', label: 'Success' },
  { value: 'failed', label: 'Failed' },
]

export function BackupsPage() {
  const { state, actions, selectors } = useStore()
  const navigate = useNavigate()
  const [restoreId, setRestoreId] = createSignal<string | null>(null)
  const [confirm, setConfirm] = createSignal('')
  const [restoreError, setRestoreError] = createSignal<string | null>(null)
  const [restoring, setRestoring] = createSignal(false)

  onMount(() => {
    void actions.backups.fetchList()
  })

  const list = () => state.backups.list.data ?? []
  const loading = () => state.backups.list.status === 'loading' && list().length === 0

  async function onRestore() {
    const id = restoreId()
    if (!id) return
    setRestoreError(null)
    setRestoring(true)
    try {
      const job = await actions.restores.restore(id, confirm())
      if (job) {
        setRestoreId(null)
        setConfirm('')
        navigate('/jobs')
      }
    } catch (e) {
      if (isApiError(e) && e.status === 400) {
        setRestoreError(e.message)
      }
    } finally {
      setRestoring(false)
    }
  }

  return (
    <>
      <PageHeader
        title="Database Backups"
        description="Artifacts stored by go-vault"
        icon={<HardDrive class="size-4" />}
        actions={
          <Button
            variant="outline"
            color="neutral"
            size="sm"
            onClick={() => void actions.backups.fetchList()}
          >
            <RefreshCw class="size-3.5" />
            Refresh
          </Button>
        }
      />

      <div class="flex-1 overflow-auto">
        <div class="flex items-center gap-1 border-b border-default px-3 py-2">
          <For each={FILTER_TABS}>
            {(tab) => (
              <button
                type="button"
                class={cn(
                  'rounded-md px-2.5 py-1 text-[12px] font-medium transition-colors',
                  state.backups.statusFilter === tab.value
                    ? 'bg-elevated text-highlighted'
                    : 'text-muted hover:bg-elevated hover:text-highlighted',
                )}
                onClick={() => {
                  actions.backups.setStatusFilter(tab.value)
                  void actions.backups.fetchList()
                }}
              >
                {tab.label}
              </button>
            )}
          </For>
        </div>

        <div class="mx-3 mt-3 flex items-start gap-2 rounded-md border border-default bg-elevated/50 px-3 py-2 text-[12px] text-muted">
          <Info class="mt-0.5 size-3.5 shrink-0 text-primary" />
          <p>
            Restore replaces data in the configured Postgres database. Confirm with the database
            name before running.
          </p>
        </div>

        <Show when={loading()}>
          <div class="flex justify-center py-16">
            <Spinner />
          </div>
        </Show>
        <Show when={!loading() && list().length === 0}>
          <EmptyState
            title="No backups"
            description="Use Create Backup in the sidebar to populate this list."
          />
        </Show>
        <Show when={!loading() && list().length > 0}>
          <ul class="mt-2 divide-y divide-default border-t border-default">
            <For each={list()}>
              {(b) => (
                <li>
                  <div
                    role="button"
                    tabindex={0}
                    class="flex cursor-pointer items-center gap-3 px-3 py-3 transition-colors hover:bg-elevated/60"
                    onClick={() => navigate(`/backups/${encodeURIComponent(b.id)}`)}
                    onKeyDown={(e) => {
                      if (e.key === 'Enter' || e.key === ' ') {
                        e.preventDefault()
                        navigate(`/backups/${encodeURIComponent(b.id)}`)
                      }
                    }}
                  >
                    <div class="min-w-0 flex-1">
                      <p class="truncate font-mono text-[13px] font-medium text-highlighted">
                        {b.name}
                      </p>
                      <p class="mt-0.5 text-[11px] text-muted">
                        {formatRelativeTime(b.created_at)}
                      </p>
                    </div>
                    <Badge
                      color={
                        b.status === 'success'
                          ? 'success'
                          : b.status === 'failed'
                            ? 'error'
                            : 'info'
                      }
                      size="sm"
                      class="uppercase"
                    >
                      {b.status}
                    </Badge>
                    <div
                      class="flex shrink-0 items-center gap-1.5"
                      onClick={(e) => e.stopPropagation()}
                      onKeyDown={(e) => e.stopPropagation()}
                    >
                      <Button
                        variant="outline"
                        color="neutral"
                        size="sm"
                        disabled={selectors.hasActiveJob() || b.status !== 'success'}
                        onClick={() => {
                          setRestoreError(null)
                          setConfirm('')
                          setRestoreId(b.id)
                        }}
                      >
                        <RotateCcw class="size-3.5" />
                        Restore
                      </Button>
                      <Button
                        variant="outline"
                        color="neutral"
                        size="sm"
                        onClick={() => void actions.backups.download(b.id, b.name)}
                      >
                        <Download class="size-3.5" />
                        Download
                      </Button>
                    </div>
                  </div>
                </li>
              )}
            </For>
          </ul>
        </Show>
      </div>

      <Modal
        open={restoreId() !== null}
        onOpenChange={(open) => {
          if (!open) {
            setRestoreId(null)
            setConfirm('')
            setRestoreError(null)
          }
        }}
        title="Restore backup"
        description="This will restore into the configured database. Type the database name to confirm."
        footer={
          <>
            <Button
              variant="outline"
              color="neutral"
              size="sm"
              onClick={() => {
                setRestoreId(null)
                setConfirm('')
                setRestoreError(null)
              }}
            >
              Cancel
            </Button>
            <Button
              size="sm"
              color="error"
              loading={restoring()}
              disabled={!confirm().trim() || selectors.hasActiveJob()}
              onClick={() => void onRestore()}
            >
              Restore
            </Button>
          </>
        }
      >
        <label class="mb-1 block text-xs font-medium text-toned">Database name</label>
        <Input
          value={confirm()}
          onInput={(e) => setConfirm(e.currentTarget.value)}
          placeholder="Type the database name"
          invalid={!!restoreError()}
        />
        <Show when={restoreError()}>
          <p class="mt-2 text-xs text-error">{restoreError()}</p>
        </Show>
      </Modal>
    </>
  )
}
