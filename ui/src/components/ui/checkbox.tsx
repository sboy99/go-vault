import { Checkbox as KCheckbox } from '@kobalte/core/checkbox'
import { Check } from 'lucide-solid'
import { splitProps } from 'solid-js'
import { cn } from '@/lib/utils/cn'

export interface CheckboxProps {
  checked?: boolean
  defaultChecked?: boolean
  onChange?: (checked: boolean) => void
  disabled?: boolean
  class?: string
  id?: string
  'aria-label'?: string
}

export function Checkbox(props: CheckboxProps) {
  const [local, rest] = splitProps(props, ['class', 'checked', 'defaultChecked', 'onChange', 'disabled'])
  return (
    <KCheckbox
      checked={local.checked}
      defaultChecked={local.defaultChecked}
      onChange={local.onChange}
      disabled={local.disabled}
      class={cn('inline-flex items-center', local.class)}
      {...rest}
    >
      <KCheckbox.Input class="peer sr-only" />
      <KCheckbox.Control
        class={cn(
          'size-4 rounded border border-accented bg-default transition-colors',
          'data-[checked]:bg-primary data-[checked]:border-primary',
          'peer-focus-visible:ring-2 peer-focus-visible:ring-primary/40',
        )}
      >
        <KCheckbox.Indicator class="flex items-center justify-center text-inverted">
          <Check class="size-3" stroke-width={3} />
        </KCheckbox.Indicator>
      </KCheckbox.Control>
    </KCheckbox>
  )
}

export interface SwitchProps {
  checked?: boolean
  onChange?: (checked: boolean) => void
  disabled?: boolean
  class?: string
  label?: string
}

export function Switch(props: SwitchProps) {
  return (
    <label class={cn('inline-flex items-center gap-2 cursor-pointer', props.class)}>
      <button
        type="button"
        role="switch"
        aria-checked={props.checked}
        disabled={props.disabled}
        onClick={() => props.onChange?.(!props.checked)}
        class={cn(
          'relative h-5 w-9 rounded-full transition-colors',
          props.checked ? 'bg-primary' : 'bg-accented',
          'disabled:opacity-75 disabled:cursor-not-allowed',
        )}
      >
        <span
          class={cn(
            'absolute top-0.5 left-0.5 size-4 rounded-full bg-white shadow transition-transform',
            props.checked && 'translate-x-4',
          )}
        />
      </button>
      {props.label ? <span class="text-sm text-default">{props.label}</span> : null}
    </label>
  )
}
