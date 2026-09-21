import type { SetStoreFunction } from 'solid-js/store'

export type ToastKind = 'success' | 'error' | 'info'

export interface Toast {
  id: string
  kind: ToastKind
  title: string
  description?: string
  createdAt: number
}

export interface ToastState {
  items: Toast[]
}

export const initialToastState: ToastState = {
  items: [],
}

let toastSeq = 0

export function createToastActions(
  _get: () => ToastState,
  set: SetStoreFunction<ToastState>,
) {
  return {
    push(kind: ToastKind, title: string, description?: string) {
      const id = `toast_${++toastSeq}_${Date.now()}`
      const toast: Toast = { id, kind, title, description, createdAt: Date.now() }
      set('items', (prev) => [...prev, toast])
      window.setTimeout(() => {
        set('items', (prev) => prev.filter((t) => t.id !== id))
      }, 4500)
      return id
    },
    dismiss(id: string) {
      set('items', (prev) => prev.filter((t) => t.id !== id))
    },
    clear() {
      set('items', [])
    },
  }
}

export type ToastActions = ReturnType<typeof createToastActions>
