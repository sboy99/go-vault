import type { Backup } from '@/domain/backup'
import type { BackupsState } from '@/store/slices/backups.slice'

export function selectBackups(state: BackupsState): Backup[] {
  return state.list.data ?? []
}

export function selectBackupCounts(state: BackupsState): {
  total: number
  success: number
  failed: number
  running: number
  storageBytes: number
} {
  const list = selectBackups(state)
  return {
    total: list.length,
    success: list.filter((b) => b.status === 'success').length,
    failed: list.filter((b) => b.status === 'failed').length,
    running: list.filter((b) => b.status === 'running').length,
    storageBytes: list.reduce((sum, b) => sum + (b.size_bytes || 0), 0),
  }
}

export function selectLatestBackup(state: BackupsState): Backup | null {
  const list = selectBackups(state)
  if (list.length === 0) return null
  return [...list].sort(
    (a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime(),
  )[0] ?? null
}
