import { For, Show, onMount } from 'solid-js'
import { A } from '@solidjs/router'
import {
  Activity,
  CheckCircle2,
  DatabaseBackup,
  HardDrive,
  RefreshCw,
  XCircle,
} from 'lucide-solid'
import { useStore } from '@/store/root'
import { PageHeader } from '@/components/layout/page-header'
import { Button } from '@/components/ui/button'
import { Card, CardBody, CardHeader } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { formatBytes, formatRelativeTime } from '@/lib/utils/format'

export function OverviewPage() {
  const { state, actions, selectors } = useStore()

  onMount(() => {
    void actions.health.check()
    void actions.backups.fetchList()
    void actions.jobs.fetchList()
  })

  const counts = () => selectors.backupCounts()
  const latest = () => selectors.latestBackup()
  const active = () => selectors.activeJobs()

  return (
    <>
      <PageHeader
        title="Overview"
        description="Backup health and recent activity"
        icon={<HardDrive class="size-4" />}
        actions={
          <Button
            variant="outline"
            color="neutral"
            size="sm"
            onClick={() => {
              void actions.health.check()
              void actions.backups.fetchList()
              void actions.jobs.fetchList()
            }}
          >
            <RefreshCw class="size-3.5" />
            Refresh
          </Button>
        }
      />
      <div class="flex-1 overflow-y-auto p-4 space-y-4">
        <Show when={active().length > 0}>
          <div class="flex items-center gap-2 rounded-lg border border-primary/30 bg-primary/5 px-3 py-2 text-sm">
            <Activity class="size-4 animate-pulse text-primary" />
            <span class="text-highlighted">
              {active().length} job{active().length === 1 ? '' : 's'} running
            </span>
            <A href="/jobs" class="ml-auto text-primary hover:underline">
              View jobs
            </A>
          </div>
        </Show>

        <div class="grid gap-3 sm:grid-cols-2 lg:grid-cols-4">
          <StatCard
            label="Health"
            value={state.health.status ?? '—'}
            hint={state.health.error ?? (state.health.ready ? `ready: ${state.health.ready}` : undefined)}
            tone={state.health.status === 'ok' ? 'success' : state.health.error ? 'error' : 'neutral'}
          />
          <StatCard label="Backups" value={String(counts().total)} hint={`${counts().success} ok · ${counts().failed} failed`} />
          <StatCard label="Storage used" value={formatBytes(counts().storageBytes)} />
          <StatCard
            label="Latest backup"
            value={latest() ? formatRelativeTime(latest()!.created_at) : 'None'}
            hint={latest()?.name}
          />
        </div>

        <Card>
          <CardHeader class="flex items-center justify-between">
            <h2 class="text-sm font-semibold text-highlighted">Recent backups</h2>
            <A href="/backups" class="text-xs text-primary hover:underline">
              View all
            </A>
          </CardHeader>
          <CardBody class="p-0">
            <Show
              when={(state.backups.list.data ?? []).length > 0}
              fallback={<p class="px-4 py-8 text-center text-sm text-muted">No backups yet</p>}
            >
              <ul class="divide-y divide-default">
                <For each={(state.backups.list.data ?? []).slice(0, 5)}>
                  {(b) => (
                    <li>
                      <A
                        href={`/backups/${encodeURIComponent(b.id)}`}
                        class="flex items-center gap-3 px-4 py-2.5 hover:bg-elevated/30"
                      >
                        <StatusIcon status={b.status} />
                        <div class="min-w-0 flex-1">
                          <p class="truncate font-mono text-sm text-highlighted">{b.name}</p>
                          <p class="text-xs text-muted">{formatRelativeTime(b.created_at)}</p>
                        </div>
                        <span class="text-xs text-muted">{formatBytes(b.size_bytes)}</span>
                        <Badge
                          color={
                            b.status === 'success' ? 'success' : b.status === 'failed' ? 'error' : 'info'
                          }
                          size="sm"
                        >
                          {b.status}
                        </Badge>
                      </A>
                    </li>
                  )}
                </For>
              </ul>
            </Show>
          </CardBody>
        </Card>
      </div>
    </>
  )
}

function StatusIcon(props: { status: string }) {
  if (props.status === 'success') return <CheckCircle2 class="size-4 text-success" />
  if (props.status === 'failed') return <XCircle class="size-4 text-error" />
  return <DatabaseBackup class="size-4 text-info" />
}

function StatCard(props: {
  label: string
  value: string
  hint?: string
  tone?: 'success' | 'error' | 'neutral'
}) {
  return (
    <Card>
      <CardBody>
        <p class="text-xs text-muted">{props.label}</p>
        <p
          class={`mt-0.5 text-[13px] font-semibold ${
            props.tone === 'success'
              ? 'text-success'
              : props.tone === 'error'
                ? 'text-error'
                : 'text-highlighted'
          }`}
        >
          {props.value}
        </p>
        <Show when={props.hint}>
          <p class="mt-0.5 truncate text-xs text-dimmed">{props.hint}</p>
        </Show>
      </CardBody>
    </Card>
  )
}
