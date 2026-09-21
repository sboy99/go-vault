import { DropdownMenu as KDropdownMenu } from '@kobalte/core/dropdown-menu'
import type { ComponentProps, JSX, ParentProps, ValidComponent } from 'solid-js'
import { splitProps } from 'solid-js'
import { cn } from '@/lib/utils/cn'

export function DropdownMenu(props: ParentProps & { modal?: boolean }) {
  return <KDropdownMenu modal={props.modal}>{props.children}</KDropdownMenu>
}

export function DropdownMenuTrigger(
  props: ParentProps & {
    class?: string
    as?: ValidComponent
  } & Record<string, unknown>,
) {
  const [local, rest] = splitProps(props, ['class', 'children', 'as'])
  return (
    <KDropdownMenu.Trigger
      as={local.as}
      {...rest}
      class={cn('inline-flex outline-none', local.class)}
    >
      {local.children}
    </KDropdownMenu.Trigger>
  )
}

export function DropdownMenuContent(
  props: ParentProps & { class?: string } & Record<string, unknown>,
) {
  const [local, rest] = splitProps(props, ['class', 'children'])
  return (
    <KDropdownMenu.Portal>
      <KDropdownMenu.Content
        {...rest}
        class={cn(
          'z-50 min-w-40 overflow-hidden rounded-md border border-default bg-muted p-1 outline-none',
          local.class,
        )}
      >
        {local.children}
      </KDropdownMenu.Content>
    </KDropdownMenu.Portal>
  )
}

export function DropdownMenuGroup(props: ParentProps & { class?: string }) {
  return (
    <KDropdownMenu.Group class={cn('p-0', props.class)}>
      {props.children}
    </KDropdownMenu.Group>
  )
}

export function DropdownMenuLabel(props: ParentProps & { class?: string }) {
  return (
    <KDropdownMenu.GroupLabel
      class={cn('px-2 py-1.5 text-xs font-medium text-muted', props.class)}
    >
      {props.children}
    </KDropdownMenu.GroupLabel>
  )
}

export function DropdownMenuItem(
  props: ParentProps & {
    onSelect?: () => void
    destructive?: boolean
    class?: string
    disabled?: boolean
    closeOnSelect?: boolean
  },
) {
  const [local, rest] = splitProps(props, [
    'class',
    'children',
    'onSelect',
    'destructive',
    'disabled',
    'closeOnSelect',
  ])
  return (
    <KDropdownMenu.Item
      {...rest}
      disabled={local.disabled}
      closeOnSelect={local.closeOnSelect}
      class={cn(
        'relative flex cursor-pointer select-none items-center gap-2 rounded-md px-2 py-1.5 text-sm outline-none',
        'data-[highlighted]:bg-elevated data-[disabled]:pointer-events-none data-[disabled]:opacity-50',
        local.destructive ? 'text-error' : 'text-default',
        local.class,
      )}
      onSelect={local.onSelect}
    >
      {local.children}
    </KDropdownMenu.Item>
  )
}

export function DropdownMenuSeparator(props: { class?: string }) {
  return (
    <KDropdownMenu.Separator
      class={cn('my-1 h-px bg-[var(--ui-border)]', props.class)}
    />
  )
}

/** Convenience wrapper kept for existing call sites. */
export function Dropdown(props: {
  trigger: JSX.Element
  children: JSX.Element
  class?: string
}) {
  return (
    <DropdownMenu>
      <DropdownMenuTrigger>{props.trigger}</DropdownMenuTrigger>
      <DropdownMenuContent class={props.class}>{props.children}</DropdownMenuContent>
    </DropdownMenu>
  )
}

export function DropdownItem(
  props: ParentProps & { onSelect?: () => void; destructive?: boolean; class?: string },
) {
  return <DropdownMenuItem {...props} />
}

export function DropdownSeparator() {
  return <DropdownMenuSeparator />
}

export type _DropdownProps = ComponentProps<'div'>
