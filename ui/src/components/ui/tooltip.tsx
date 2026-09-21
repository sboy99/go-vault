import { Tooltip as KTooltip } from '@kobalte/core/tooltip'
import type { JSX, ParentProps } from 'solid-js'
import { cn } from '@/lib/utils/cn'

export function Tooltip(props: ParentProps & { content: string; class?: string }) {
  return (
    <KTooltip openDelay={200}>
      <KTooltip.Trigger as="span" class={cn('inline-flex', props.class)}>
        {props.children}
      </KTooltip.Trigger>
      <KTooltip.Portal>
        <KTooltip.Content class="z-50 rounded-md bg-inverted px-2 py-1 text-xs text-inverted shadow">
          {props.content}
          <KTooltip.Arrow />
        </KTooltip.Content>
      </KTooltip.Portal>
    </KTooltip>
  )
}

export type _TooltipChildren = JSX.Element
