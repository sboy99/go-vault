# States

| State | Treatment |
|---|---|
| Default | Surface + border |
| Hover | `bg-elevated` or opacity 0.9 on solid |
| Active / pressed | `bg-accented` |
| Focus | 1px outline primary |
| Disabled | opacity 50%, `cursor-not-allowed` |
| Loading | spinner in button / async boundary |
| Error | `text-error` / `border-error` |
| Success | `text-success` badges/icons |

Keep transitions short (`transition-colors`, ~150ms).
