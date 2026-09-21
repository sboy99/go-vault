import { Select as KSelect } from '@kobalte/core/select'
import { Check, ChevronDown } from 'lucide-solid'
import { Show } from 'solid-js'
import { cn } from '@/lib/utils/cn'

export interface SelectOption {
  value: string
  label: string
}

export interface SelectProps {
  options: SelectOption[]
  value?: string
  placeholder?: string
  onChange?: (value: string | null) => void
  class?: string
  disabled?: boolean
}

export function Select(props: SelectProps) {
  const selected = () => props.options.find((o) => o.value === props.value) ?? null

  return (
    <KSelect
      options={props.options}
      optionValue="value"
      optionTextValue="label"
      value={selected()}
      onChange={(opt) => props.onChange?.(opt?.value ?? null)}
      placeholder={props.placeholder ?? 'Select…'}
      disabled={props.disabled}
      itemComponent={(itemProps) => (
        <KSelect.Item
          item={itemProps.item}
          class={cn(
            'relative flex cursor-pointer select-none items-center rounded-md px-2 py-1.5 text-sm outline-none',
            'data-[highlighted]:bg-elevated data-[selected]:text-primary',
          )}
        >
          <KSelect.ItemLabel>{itemProps.item.rawValue.label}</KSelect.ItemLabel>
          <KSelect.ItemIndicator class="absolute right-2">
            <Check class="size-3.5 text-primary" />
          </KSelect.ItemIndicator>
        </KSelect.Item>
      )}
    >
      <KSelect.Trigger
        class={cn(
          'inline-flex h-8 w-full items-center justify-between gap-2 rounded-md px-2.5 text-sm',
          'bg-default text-highlighted shadow-[inset_0_0_0_1px_var(--ui-border-accented)]',
          'focus:outline-none focus:ring-2 focus:ring-primary/40',
          'disabled:opacity-75 disabled:cursor-not-allowed',
          props.class,
        )}
      >
        <KSelect.Value<SelectOption>>
          {(state) => (
            <Show when={state.selectedOption()} fallback={<span class="text-dimmed">{props.placeholder}</span>}>
              {(opt) => opt().label}
            </Show>
          )}
        </KSelect.Value>
        <KSelect.Icon>
          <ChevronDown class="size-4 text-muted" />
        </KSelect.Icon>
      </KSelect.Trigger>
      <KSelect.Portal>
        <KSelect.Content class="z-50 min-w-[var(--kb-select-trigger-width)] overflow-hidden rounded-md border border-default bg-default shadow-lg">
          <KSelect.Listbox class="max-h-60 overflow-auto p-1" />
        </KSelect.Content>
      </KSelect.Portal>
    </KSelect>
  )
}
