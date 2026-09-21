export type JobStatus = 'queued' | 'running' | 'succeeded' | 'failed'

export type JobType = 'backup' | 'restore'

export interface Job {
  id: string
  type: JobType
  status: JobStatus
  error?: string
  backup_id?: string
  created_at: string
  started_at?: string
  finished_at?: string
}

export function isActiveJob(job: Job): boolean {
  return job.status === 'queued' || job.status === 'running'
}
