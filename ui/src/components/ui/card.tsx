import { splitProps, type JSX, type ParentProps } from 'solid-js'
import { cn } from '@/lib/utils/cn'

export function Card(props: ParentProps & JSX.HTMLAttributes<HTMLDivElement>) {
  const [local, rest] = splitProps(props, ['class', 'children'])
  return (
    <div
      {...rest}
      class={cn('rounded-md border border-default bg-muted', local.class)}
    >
      {local.children}
    </div>
  )
}

export function CardHeader(props: ParentProps & JSX.HTMLAttributes<HTMLDivElement>) {
  const [local, rest] = splitProps(props, ['class', 'children'])
  return (
    <div {...rest} class={cn('px-3 py-2 border-b border-default', local.class)}>
      {local.children}
    </div>
  )
}

export function CardBody(props: ParentProps & JSX.HTMLAttributes<HTMLDivElement>) {
  const [local, rest] = splitProps(props, ['class', 'children'])
  return (
    <div {...rest} class={cn('px-3 py-2', local.class)}>
      {local.children}
    </div>
  )
}
