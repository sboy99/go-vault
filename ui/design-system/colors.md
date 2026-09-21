# Colors

Source: `src/styles/tokens.css`. Runtime overrides via theme picker update `--ui-color-primary-*` and `--ui-color-neutral-*`.

## Semantic (use these in components)

| Token | Role | Dark default |
|---|---|---|
| `--ui-bg` | App canvas | `#08090A` (`neutral-950`) |
| `--ui-bg-muted` | Sidebar / card surface | `#0D0E10` |
| `--ui-bg-elevated` | Hover / elevated | `#17191C` |
| `--ui-bg-accented` | Active / pressed | `#1D2024` |
| `--ui-border` | Default hairline | `#272A30` |
| `--ui-border-muted` | Softer divider | `#1C1F23` |
| `--ui-border-accented` | Stronger edge | `#343840` |
| `--ui-text-highlighted` | Headings | `#F7F8F8` |
| `--ui-text` | Body | `#C8CACD` |
| `--ui-text-toned` | Secondary | `#A7A9AD` |
| `--ui-text-muted` | Labels / meta | `#6F737A` |
| `--ui-text-dimmed` | Disabled / hint | `#4D5158` |
| `--ui-primary` | Accent | `#5E6AD2` |
| `--ui-primary-hover` | Accent hover | `#6872D9` |
| `--ui-primary-muted` | Accent wash | `rgba(94,106,210,0.15)` |
| `--ui-success` | Success | `#3FB950` |
| `--ui-warning` | Warning | `#D29922` |
| `--ui-error` | Error | `#F85149` |
| `--ui-info` | Info | `#58A6FF` |

Utility classes: `bg-default`, `bg-muted`, `bg-elevated`, `text-highlighted`, `text-muted`, `border-default`, `text-primary`, etc. (`src/styles/utilities.css`).

## Defaults

- Primary palette id: `indigo`
- Neutral palette id: `linear`
