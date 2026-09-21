import { Dialog } from '@kobalte/core/dialog'
import { X } from 'lucide-solid'
import type { JSX } from 'solid-js'
import { Show } from 'solid-js'
import { cn } from '@/lib/utils/cn'

export interface ModalProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  title?: string
  description?: string
  children?: JSX.Element
  footer?: JSX.Element
  class?: string
}

export function Modal(props: ModalProps) {
  return (
    <Dialog open={props.open} onOpenChange={props.onOpenChange}>
      <Dialog.Portal>
        <Dialog.Overlay class="fixed inset-0 z-50 bg-black/40" />
        <div class="fixed inset-0 z-50 flex items-center justify-center p-4">
          <Dialog.Content
            class={cn(
              'w-full max-w-lg rounded-lg border border-default bg-default shadow-lg outline-none',
              props.class,
            )}
          >
            <div class="flex items-start justify-between gap-3 border-b border-default px-4 py-3">
              <div>
                <Show when={props.title}>
                  <Dialog.Title class="text-sm font-semibold text-highlighted">
                    {props.title}
                  </Dialog.Title>
                </Show>
                <Show when={props.description}>
                  <Dialog.Description class="mt-0.5 text-sm text-muted">
                    {props.description}
                  </Dialog.Description>
                </Show>
              </div>
              <Dialog.CloseButton
                class="inline-flex size-7 items-center justify-center rounded-md text-muted transition-colors hover:bg-elevated hover:text-highlighted"
                aria-label="Close"
              >
                <X class="size-4" />
              </Dialog.CloseButton>
            </div>
            <div class="px-4 py-3">{props.children}</div>
            <Show when={props.footer}>
              <div class="flex justify-end gap-2 border-t border-default px-4 py-3">
                {props.footer}
              </div>
            </Show>
          </Dialog.Content>
        </div>
      </Dialog.Portal>
    </Dialog>
  )
}
