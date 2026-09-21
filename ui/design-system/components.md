# Components

Assembly order:

```
background → border → typography → state → interaction
```

Avoid:

- Large shadows
- Gradients
- Overly rounded cards
- Unnecessary animation
- Competing accent colors

SolidJS conventions:

- Primitives in `src/components/ui`
- Layout in `src/components/layout`
- Feature pages compose primitives; no one-off color hexes in pages
- Theme via CSS variables + store actions (`setPrimary`, `setNeutral`, `setColorMode`)
