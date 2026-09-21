import { Show, createSignal, onMount } from 'solid-js'
import { Settings as SettingsIcon } from 'lucide-solid'
import { useStore } from '@/store/root'
import { PageHeader } from '@/components/layout/page-header'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Card, CardBody, CardHeader } from '@/components/ui/card'

export function SettingsPage() {
  const { state, actions } = useStore()
  const [tokenDraft, setTokenDraft] = createSignal(state.settings.apiToken)
  const [testing, setTesting] = createSignal(false)
  const [testResult, setTestResult] = createSignal<string | null>(null)

  onMount(() => setTokenDraft(state.settings.apiToken))

  async function saveToken() {
    actions.settings.setApiToken(tokenDraft())
    actions.toast.push('success', 'Token saved')
  }

  async function testConnection() {
    actions.settings.setApiToken(tokenDraft())
    setTesting(true)
    setTestResult(null)
    try {
      const ok = await actions.health.check()
      if (ok) {
        setTestResult(`Connected — health: ${state.health.status}`)
        actions.toast.push('success', 'Connection OK')
      } else {
        setTestResult(state.health.error ?? 'Connection failed')
        actions.toast.push('error', 'Connection failed', state.health.error ?? undefined)
      }
    } finally {
      setTesting(false)
    }
  }

  return (
    <>
      <PageHeader
        title="Settings"
        description="API token for go-vault"
        icon={<SettingsIcon class="size-4" />}
      />
      <div class="flex-1 overflow-y-auto p-4">
        <div class="mx-auto max-w-2xl space-y-4">
          <Card>
            <CardHeader>
              <h2 class="text-sm font-semibold text-highlighted">API access</h2>
              <p class="text-xs text-muted">
                Bearer token sent with every request. Leave empty if the server has no token configured.
              </p>
            </CardHeader>
            <CardBody class="space-y-3">
              <Show when={state.settings.unauthorized}>
                <div class="rounded-md bg-warning/10 px-3 py-2 text-xs text-warning">
                  The API returned 401 Unauthorized. Update the token below.
                </div>
              </Show>
              <div>
                <label class="mb-1 block text-xs font-medium text-toned">API token</label>
                <Input
                  type="password"
                  autocomplete="off"
                  value={tokenDraft()}
                  onInput={(e) => setTokenDraft(e.currentTarget.value)}
                  placeholder="Bearer token"
                />
              </div>
              <div class="flex flex-wrap gap-2">
                <Button size="sm" onClick={() => void saveToken()}>
                  Save token
                </Button>
                <Button
                  size="sm"
                  variant="outline"
                  color="neutral"
                  loading={testing()}
                  onClick={() => void testConnection()}
                >
                  Test connection
                </Button>
              </div>
              <Show when={testResult()}>
                <p class="text-xs text-muted">{testResult()}</p>
              </Show>
            </CardBody>
          </Card>
        </div>
      </div>
    </>
  )
}
