import { splitProps, type JSX, type ParentProps } from 'solid-js'
import { cn } from '@/lib/utils/cn'

export function Table(props: ParentProps & JSX.HTMLAttributes<HTMLTableElement>) {
  const [local, rest] = splitProps(props, ['class', 'children'])
  return (
    <table {...rest} class={cn('w-full text-[12px]', local.class)}>
      {local.children}
    </table>
  )
}

export function THead(props: ParentProps & JSX.HTMLAttributes<HTMLTableSectionElement>) {
  const [local, rest] = splitProps(props, ['class', 'children'])
  return (
    <thead {...rest} class={cn(local.class)}>
      {local.children}
    </thead>
  )
}

export function TBody(props: ParentProps & JSX.HTMLAttributes<HTMLTableSectionElement>) {
  const [local, rest] = splitProps(props, ['class', 'children'])
  return (
    <tbody {...rest} class={cn(local.class)}>
      {local.children}
    </tbody>
  )
}

export function TR(props: ParentProps & JSX.HTMLAttributes<HTMLTableRowElement>) {
  const [local, rest] = splitProps(props, ['class', 'children'])
  return (
    <tr
      {...rest}
      class={cn(
        'border-b border-default transition-colors hover:bg-elevated',
        local.class,
      )}
    >
      {local.children}
    </tr>
  )
}

export function TH(props: ParentProps & JSX.ThHTMLAttributes<HTMLTableCellElement>) {
  const [local, rest] = splitProps(props, ['class', 'children'])
  return (
    <th
      {...rest}
      class={cn(
        'px-2 py-1.5 text-left text-[12px] font-medium text-muted whitespace-nowrap',
        local.class,
      )}
    >
      {local.children}
    </th>
  )
}

export function TD(props: ParentProps & JSX.TdHTMLAttributes<HTMLTableCellElement>) {
  const [local, rest] = splitProps(props, ['class', 'children'])
  return (
    <td {...rest} class={cn('px-2 py-1 text-default align-middle', local.class)}>
      {local.children}
    </td>
  )
}
