import type { Job } from '@/domain/job'
import { isActiveJob } from '@/domain/job'
import type { JobsState } from '@/store/slices/jobs.slice'

export function selectJobs(state: JobsState): Job[] {
  return state.list.data ?? []
}

export function selectActiveJobs(state: JobsState): Job[] {
  return selectJobs(state).filter(isActiveJob)
}

export function selectHasActiveJob(state: JobsState): boolean {
  return selectActiveJobs(state).length > 0
}

export function selectLatestActiveJob(state: JobsState): Job | null {
  const active = selectActiveJobs(state)
  return active[0] ?? null
}
