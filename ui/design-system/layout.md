# Layout

## App chrome

```
┌──────────────────────────────────────┐
│ AppHeader (h = --ui-header-height)   │  40px
├──────────────────────────────────────┤
│ PageHeader (optional)                │  40px
├──────────────────────────────────────┤
│ Main scroll region                   │
│  padding: 12–16px                    │
└──────────────────────────────────────┘
```

- No persistent sidebar (nav via dropdown + Cmd+K)
- Header: brand, nav dropdown, theme picker, mode toggle
- `--ui-header-height: 2.5rem`

## Density

Prefer single-column content with max readable width for settings (~`max-w-2xl`). Data tables use full width.
