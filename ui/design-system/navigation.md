# Navigation

## Patterns in go-vault

1. **Header nav dropdown** — current page label + chevron; items jump routes
2. **Command menu (⌘K)** — fuzzy jump to Overview / Backups / Jobs / Settings
3. **In-page links** — `text-primary` for cross-links (View all, job → backup)

## Rules

- Active item: `bg-elevated` + `text-primary` icon
- Inactive icons: `text-muted`
- No heavy sidebar
- Keep one clear “where am I” label in the header trigger
