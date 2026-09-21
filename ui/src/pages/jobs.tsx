import { For, Show, onMount } from 'solid-js'
import { A } from '@solidjs/router'
import { ListTodo, RefreshCw } from 'lucide-solid'
import { useStore } from '@/store/root'
import { PageHeader } from '@/components/layout/page-header'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Table, THead, TBody, TR, TH, TD } from '@/components/ui/table'
import { EmptyState } from '@/components/feedback/empty-state'
import { Spinner } from '@/components/ui/spinner'
import { formatDuration, formatRelativeTime } from '@/lib/utils/format'
import type { JobStatus } from '@/domain/job'

function statusColor(status: JobStatus): 'success' | 'error' | 'info' | 'warning' {
  if (status === 'succeeded') return 'success'
  if (status === 'failed') return 'error'
  if (status === 'running') return 'info'
  return 'warning'
}

export function JobsPage() {
  const { state, actions } = useStore()

  onMount(() => {
    void actions.jobs.fetchList()
  })

  const list = () => state.jobs.list.data ?? []
  const loading = () => state.jobs.list.status === 'loading' && list().length === 0

  return (
    <>
      <PageHeader
        title="Jobs"
        description="Async backup and restore operations"
        icon={<ListTodo class="size-4" />}
        actions={
          <Button
            variant="outline"
            color="neutral"
            size="sm"
            onClick={() => void actions.jobs.fetchList()}
          >
            <RefreshCw class="size-3.5" />
            Refresh
            <Show when={state.jobs.polling}>
              <span class="ml-1 size-1.5 rounded-full bg-primary animate-pulse" />
            </Show>
          </Button>
        }
      />
      <div class="flex-1 overflow-auto">
        <Show when={loading()}>
          <div class="flex justify-center py-16">
            <Spinner />
          </div>
        </Show>
        <Show when={!loading() && list().length === 0}>
          <EmptyState title="No jobs" description="Jobs appear when you create a backup or restore." />
        </Show>
        <Show when={!loading() && list().length > 0}>
          <Table>
            <THead>
              <TR class="hover:bg-transparent">
                <TH>ID</TH>
                <TH>Type</TH>
                <TH>Status</TH>
                <TH>Backup</TH>
                <TH>Duration</TH>
                <TH>Created</TH>
                <TH>Error</TH>
              </TR>
            </THead>
            <TBody>
              <For each={list()}>
                {(j) => (
                  <TR>
                    <TD class="font-mono text-xs">{j.id}</TD>
                    <TD class="capitalize">{j.type}</TD>
                    <TD>
                      <Badge color={statusColor(j.status)} size="sm">
                        {j.status}
                      </Badge>
                    </TD>
                    <TD>
                      <Show when={j.backup_id} fallback={<span class="text-dimmed">—</span>}>
                        <A
                          href={`/backups/${encodeURIComponent(j.backup_id!)}`}
                          class="font-mono text-xs text-primary hover:underline"
                        >
                          {j.backup_id}
                        </A>
                      </Show>
                    </TD>
                    <TD class="text-xs text-muted">
                      {formatDuration(j.started_at, j.finished_at)}
                    </TD>
                    <TD class="text-xs text-muted">{formatRelativeTime(j.created_at)}</TD>
                    <TD class="max-w-xs truncate text-xs text-error">{j.error || '—'}</TD>
                  </TR>
                )}
              </For>
            </TBody>
          </Table>
        </Show>
      </div>
    </>
  )
}
