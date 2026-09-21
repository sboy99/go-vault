import { Show, createSignal, onCleanup, onMount } from 'solid-js'
import { A, useNavigate, useParams } from '@solidjs/router'
import { ArrowLeft, Download, HardDrive, RotateCcw } from 'lucide-solid'
import { useStore } from '@/store/root'
import { PageHeader } from '@/components/layout/page-header'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Card, CardBody, CardHeader } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Modal } from '@/components/ui/modal'
import { Spinner } from '@/components/ui/spinner'
import { isApiError } from '@/lib/http/errors'
import {
  formatBytes,
  formatDuration,
  formatRelativeTime,
} from '@/lib/utils/format'

export function BackupDetailPage() {
  const params = useParams()
  const navigate = useNavigate()
  const { state, actions, selectors } = useStore()
  const [restoreOpen, setRestoreOpen] = createSignal(false)
  const [confirm, setConfirm] = createSignal('')
  const [restoreError, setRestoreError] = createSignal<string | null>(null)
  const [restoring, setRestoring] = createSignal(false)

  onMount(() => {
    void actions.backups.fetchById(params.id!)
  })
  onCleanup(() => actions.backups.clearSelected())

  const backup = () => state.backups.selected.data
  const loading = () => state.backups.selected.status === 'loading'

  async function onRestore() {
    setRestoreError(null)
    setRestoring(true)
    try {
      const job = await actions.restores.restore(params.id!, confirm())
      if (job) {
        setRestoreOpen(false)
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
        title={backup()?.name ?? 'Backup'}
        description={backup()?.id}
        icon={<HardDrive class="size-4" />}
        actions={
          <>
            <Button variant="ghost" color="neutral" size="sm" onClick={() => navigate('/backups')}>
              <ArrowLeft class="size-3.5" />
              Back
            </Button>
            <Show when={backup()}>
              <Button
                variant="outline"
                color="neutral"
                size="sm"
                onClick={() => void actions.backups.download(backup()!.id, backup()!.name)}
              >
                <Download class="size-3.5" />
                Download
              </Button>
              <Button
                variant="outline"
                color="neutral"
                size="sm"
                disabled={selectors.hasActiveJob() || backup()?.status !== 'success'}
                onClick={() => setRestoreOpen(true)}
              >
                <RotateCcw class="size-3.5" />
                Restore
              </Button>
            </Show>
          </>
        }
      />
      <div class="flex-1 overflow-y-auto p-4">
        <Show when={loading()}>
          <div class="flex justify-center py-16">
            <Spinner />
          </div>
        </Show>
        <Show when={state.backups.selected.status === 'error'}>
          <p class="text-sm text-error">{state.backups.selected.error}</p>
          <A href="/backups" class="mt-2 inline-block text-sm text-primary">
            Back to list
          </A>
        </Show>
        <Show when={backup()}>
          {(b) => (
            <div class="mx-auto max-w-3xl space-y-4">
              <Card>
                <CardHeader class="flex items-center justify-between">
                  <h2 class="text-sm font-semibold text-highlighted">Metadata</h2>
                  <Badge
                    color={
                      b().status === 'success'
                        ? 'success'
                        : b().status === 'failed'
                          ? 'error'
                          : 'info'
                    }
                  >
                    {b().status}
                  </Badge>
                </CardHeader>
                <CardBody>
                  <dl class="grid gap-3 sm:grid-cols-2">
                    <Field label="ID" mono value={b().id} />
                    <Field label="Storage key" mono value={b().storage_key} />
                    <Field label="Database" value={b().database_type} />
                    <Field label="Storage" value={b().storage_type} />
                    <Field label="Format" value={b().format} />
                    <Field label="Size" value={formatBytes(b().size_bytes)} />
                    <Field label="Duration" value={formatDuration(b().started_at, b().finished_at)} />
                    <Field label="Verified" value={b().verified ? 'Yes' : 'No'} />
                    <Field label="PG version" value={b().pg_version || '—'} />
                    <Field label="Created" value={formatRelativeTime(b().created_at)} />
                    <Field label="SHA256" mono value={b().sha256 || '—'} class="sm:col-span-2" />
                    <Show when={b().error}>
                      <Field label="Error" value={b().error!} class="sm:col-span-2 text-error" />
                    </Show>
                  </dl>
                </CardBody>
              </Card>
            </div>
          )}
        </Show>
      </div>

      <Modal
        open={restoreOpen()}
        onOpenChange={(v) => {
          setRestoreOpen(v)
          if (!v) {
            setConfirm('')
            setRestoreError(null)
          }
        }}
        title="Restore backup"
        description="This will restore into the configured database. Type the database name to confirm."
        footer={
          <>
            <Button variant="outline" color="neutral" size="sm" onClick={() => setRestoreOpen(false)}>
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

function Field(props: { label: string; value: string; mono?: boolean; class?: string }) {
  return (
    <div class={props.class}>
      <dt class="text-xs text-muted">{props.label}</dt>
      <dd class={`mt-0.5 text-sm text-highlighted break-all ${props.mono ? 'font-mono' : ''}`}>
        {props.value}
      </dd>
    </div>
  )
}
