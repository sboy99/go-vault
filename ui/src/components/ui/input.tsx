import { splitProps, type JSX } from 'solid-js'
import { cn } from '@/lib/utils/cn'

export interface InputProps extends Omit<JSX.InputHTMLAttributes<HTMLInputElement>, 'size'> {
  size?: 'sm' | 'md' | 'lg'
  invalid?: boolean
}

export function Input(props: InputProps) {
  const [local, rest] = splitProps(props, ['class', 'size', 'invalid'])
  const size = () => local.size ?? 'md'
  return (
    <input
      {...rest}
      class={cn(
        'w-full rounded-md border border-default bg-default text-highlighted appearance-none',
        'placeholder:text-dimmed transition-colors',
        'focus:outline-none focus:border-primary',
        'disabled:cursor-not-allowed disabled:opacity-50',
        size() === 'sm' && 'px-2 py-1 text-[12px] h-6',
        size() === 'md' && 'px-2 py-1 text-[13px] h-7',
        size() === 'lg' && 'px-2.5 py-1.5 text-[13px] h-8',
        local.invalid && 'border-error',
        local.class,
      )}
    />
  )
}
