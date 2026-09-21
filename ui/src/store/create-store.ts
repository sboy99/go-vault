import { createStore, type SetStoreFunction, type Store } from 'solid-js/store'

export type ActionMap = Record<string, (...args: never[]) => unknown>

export interface SliceDefinition<S, A extends ActionMap> {
  name: string
  initialState: S
  createActions: (ctx: {
    get: () => S
    set: SetStoreFunction<S>
    getRoot?: () => unknown
  }) => A
}

export function createSliceStore<S extends object, A extends ActionMap>(
  initialState: S,
  createActions: (get: () => S, set: SetStoreFunction<S>) => A,
): { state: Store<S>; setState: SetStoreFunction<S>; actions: A } {
  const [state, setState] = createStore(initialState)
  const actions = createActions(() => state, setState)
  return { state, setState, actions }
}

export function logAction(slice: string, action: string, payload?: unknown): void {
  if (import.meta.env.DEV) {
    // eslint-disable-next-line no-console
    console.debug(`[store:${slice}] ${action}`, payload)
  }
}
