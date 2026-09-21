import { STORAGE_KEYS } from '@/config/constants'
import { env } from '@/config/env'
import { readStorage, removeStorage, writeStorage } from '@/lib/utils/storage'
import type { SetStoreFunction } from 'solid-js/store'
import { logAction } from '@/store/create-store'

export interface SettingsState {
  apiToken: string
  unauthorized: boolean
}

function loadToken(): string {
  return readStorage(STORAGE_KEYS.apiToken) ?? env.apiToken ?? ''
}

export const initialSettingsState: SettingsState = {
  apiToken: loadToken(),
  unauthorized: false,
}

export function createSettingsActions(
  _get: () => SettingsState,
  set: SetStoreFunction<SettingsState>,
) {
  return {
    setApiToken(token: string) {
      logAction('settings', 'setApiToken')
      set('apiToken', token)
      if (token) writeStorage(STORAGE_KEYS.apiToken, token)
      else removeStorage(STORAGE_KEYS.apiToken)
      set('unauthorized', false)
    },
    setUnauthorized(value: boolean) {
      set('unauthorized', value)
    },
  }
}

export type SettingsActions = ReturnType<typeof createSettingsActions>
