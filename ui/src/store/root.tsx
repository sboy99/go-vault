import {
  createContext,
  createEffect,
  onCleanup,
  onMount,
  useContext,
  type ParentProps,
} from 'solid-js'
import { createStore, type SetStoreFunction } from 'solid-js/store'
import type { Job } from '@/domain/job'
import { isApiError } from '@/lib/http/errors'
import { createRepositories } from '@/repositories'
import type { Repositories } from '@/repositories/types'
import { selectBackupCounts, selectLatestBackup } from '@/store/selectors/backups.selectors'
import { selectActiveJobs, selectHasActiveJob } from '@/store/selectors/jobs.selectors'
import {
  createBackupsActions,
  initialBackupsState,
  type BackupsActions,
  type BackupsState,
} from '@/store/slices/backups.slice'
import {
  createJobsActions,
  initialJobsState,
  type JobsActions,
  type JobsState,
} from '@/store/slices/jobs.slice'
import {
  createSettingsActions,
  initialSettingsState,
  type SettingsActions,
  type SettingsState,
} from '@/store/slices/settings.slice'
import {
  createThemeActions,
  initialThemeState,
  type ThemeActions,
  type ThemeState,
} from '@/store/slices/theme.slice'
import {
  createToastActions,
  initialToastState,
  type ToastActions,
  type ToastState,
} from '@/store/slices/toast.slice'

export interface RootState {
  backups: BackupsState
  jobs: JobsState
  settings: SettingsState
  theme: ThemeState
  toast: ToastState
  health: {
    status: string | null
    ready: string | null
    error: string | null
    checking: boolean
  }
}

export interface RootActions {
  backups: BackupsActions
  jobs: JobsActions
  settings: SettingsActions
  theme: ThemeActions
  toast: ToastActions
  health: { check: () => Promise<boolean> }
  restores: { restore: (backupId: string, confirm: string) => Promise<Job | null> }
}

export interface StoreContextValue {
  state: RootState
  actions: RootActions
  selectors: {
    backupCounts: () => ReturnType<typeof selectBackupCounts>
    latestBackup: () => ReturnType<typeof selectLatestBackup>
    hasActiveJob: () => boolean
    activeJobs: () => Job[]
  }
  repos: Repositories
}

const StoreContext = createContext<StoreContextValue>()

function nestSet<T extends object>(
  setRoot: SetStoreFunction<RootState>,
  key: keyof RootState,
): SetStoreFunction<T> {
  return ((...args: unknown[]) => {
    ;(setRoot as (...a: unknown[]) => void)(key, ...args)
  }) as SetStoreFunction<T>
}

export function StoreProvider(props: ParentProps & { repos?: Repositories }) {
  const [state, setState] = createStore<RootState>({
    backups: initialBackupsState,
    jobs: initialJobsState,
    settings: initialSettingsState,
    theme: initialThemeState,
    toast: initialToastState,
    health: { status: null, ready: null, error: null, checking: false },
  })

  const onUnauthorized = () => setState('settings', 'unauthorized', true)

  const repos =
    props.repos ??
    createRepositories({
      getToken: () => state.settings.apiToken || null,
      onUnauthorized,
    })

  const setToast = nestSet<ToastState>(setState, 'toast')
  const toastActions = createToastActions(() => state.toast, setToast)
  const notify = (kind: 'success' | 'error' | 'info', title: string, description?: string) => {
    toastActions.push(kind, title, description)
  }

  const setJobs = nestSet<JobsState>(setState, 'jobs')
  const setBackups = nestSet<BackupsState>(setState, 'backups')
  const setSettings = nestSet<SettingsState>(setState, 'settings')
  const setTheme = nestSet<ThemeState>(setState, 'theme')

  const backupsActions = createBackupsActions(
    () => state.backups,
    setBackups,
    repos,
    notify,
    onUnauthorized,
  )

  const jobsActions = createJobsActions(
    () => state.jobs,
    setJobs,
    repos,
    notify,
    onUnauthorized,
    (job) => {
      if (job.type === 'backup') void backupsActions.fetchList()
    },
  )

  const settingsActions = createSettingsActions(() => state.settings, setSettings)
  const themeActions = createThemeActions(() => state.theme, setTheme)

  const actions: RootActions = {
    backups: backupsActions,
    jobs: jobsActions,
    settings: settingsActions,
    theme: themeActions,
    toast: toastActions,
    health: {
      async check() {
        setState('health', 'checking', true)
        setState('health', 'error', null)
        try {
          const h = await repos.health.healthz()
          setState('health', 'status', h.status)
          try {
            const r = await repos.health.readyz()
            setState('health', 'ready', r.status)
          } catch (e) {
            setState('health', 'ready', null)
            setState('health', 'error', e instanceof Error ? e.message : 'not ready')
          }
          return true
        } catch (e) {
          if (isApiError(e) && e.isUnauthorized) onUnauthorized()
          setState('health', 'error', e instanceof Error ? e.message : 'health check failed')
          setState('health', 'status', null)
          return false
        } finally {
          setState('health', 'checking', false)
        }
      },
    },
    restores: {
      async restore(backupId, confirm) {
        try {
          const job = await repos.restores.create({ backup_id: backupId, confirm })
          jobsActions.upsertJob(job)
          toastActions.push('success', 'Restore started', `Job ${job.id}`)
          return job
        } catch (e) {
          if (isApiError(e) && e.isUnauthorized) onUnauthorized()
          if (isApiError(e) && e.isConflict) {
            toastActions.push('info', 'Job already running', e.message)
            return null
          }
          toastActions.push('error', 'Restore', e instanceof Error ? e.message : 'Restore failed')
          throw e
        }
      },
    },
  }

  onMount(() => {
    themeActions.init()
    onCleanup(() => {
      jobsActions.stopPolling()
    })
  })

  createEffect(() => {
    void state.settings.apiToken
  })

  const value: StoreContextValue = {
    state,
    actions,
    selectors: {
      backupCounts: () => selectBackupCounts(state.backups),
      latestBackup: () => selectLatestBackup(state.backups),
      hasActiveJob: () => selectHasActiveJob(state.jobs),
      activeJobs: () => selectActiveJobs(state.jobs),
    },
    repos,
  }

  return <StoreContext.Provider value={value}>{props.children}</StoreContext.Provider>
}

export function useStore(): StoreContextValue {
  const ctx = useContext(StoreContext)
  if (!ctx) throw new Error('useStore must be used within StoreProvider')
  return ctx
}
