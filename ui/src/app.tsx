import { Route, Router } from '@solidjs/router'
import { lazy, type ParentProps } from 'solid-js'
import { StoreProvider } from '@/store/root'
import { AppLayout } from '@/components/layout/app-shell'

const OverviewPage = lazy(() => import('@/pages/overview').then((m) => ({ default: m.OverviewPage })))
const BackupsPage = lazy(() => import('@/pages/backups').then((m) => ({ default: m.BackupsPage })))
const BackupDetailPage = lazy(() =>
  import('@/pages/backup-detail').then((m) => ({ default: m.BackupDetailPage })),
)
const JobsPage = lazy(() => import('@/pages/jobs').then((m) => ({ default: m.JobsPage })))
const SettingsPage = lazy(() => import('@/pages/settings').then((m) => ({ default: m.SettingsPage })))

function Layout(props: ParentProps) {
  return <AppLayout>{props.children}</AppLayout>
}

export function App() {
  return (
    <StoreProvider>
      <Router root={Layout}>
        <Route path="/" component={OverviewPage} />
        <Route path="/backups" component={BackupsPage} />
        <Route path="/backups/:id" component={BackupDetailPage} />
        <Route path="/jobs" component={JobsPage} />
        <Route path="/settings" component={SettingsPage} />
      </Router>
    </StoreProvider>
  )
}
