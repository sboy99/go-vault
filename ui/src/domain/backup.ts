export type DatabaseType = 'POSTGRESQL'

export type StorageType = 'LOCAL' | 'CLOUD'

export type BackupStatus = 'running' | 'success' | 'failed'

export interface Backup {
  id: string
  name: string
  storage_key: string
  database_type: DatabaseType
  storage_type: StorageType
  status: BackupStatus
  error?: string
  started_at: string
  finished_at?: string
  created_at: string
  size_bytes: number
  sha256?: string
  pg_version?: string
  format: string
  verified: boolean
}
