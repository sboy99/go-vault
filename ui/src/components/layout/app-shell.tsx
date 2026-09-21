import { ErrorBoundary, Suspense, type JSX } from 'solid-js'
import { AppHeader } from '@/components/layout/app-header'
import { Sidebar } from '@/components/layout/sidebar'
import { MainPanel } from '@/components/layout/page-header'
import { CommandMenu } from '@/components/layout/command-menu'
import { Toaster } from '@/components/feedback/toaster'
import { ErrorBoundaryFallback } from '@/components/feedback/empty-state'

export function AppLayout(props: { children: JSX.Element }) {
  return (
    <div class="flex h-full min-h-0 flex-col">
      <AppHeader />
      <div class="flex min-h-0 flex-1">
        <Sidebar />
        <MainPanel>
          <ErrorBoundary
            fallback={(err, reset) => <ErrorBoundaryFallback error={err} reset={reset} />}
          >
            <Suspense
              fallback={
                <div class="flex flex-1 items-center justify-center text-[13px] text-muted">
                  Loading…
                </div>
              }
            >
              {props.children}
            </Suspense>
          </ErrorBoundary>
        </MainPanel>
      </div>
      <CommandMenu />
      <Toaster />
    </div>
  )
}
