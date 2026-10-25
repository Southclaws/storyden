# Directory layouts

The page's Directory block stores one shared layout in its metadata. Table remains the default; Grid and Gallery use the same fields, controls and page sequence.

- Grid uses aligned rows and a 4:3 image crop with `object-fit: cover`.
- Gallery uses balanced CSS columns, read down each column. Images keep their natural height between 32px and `min(70svh, 40rem)`. Ordinary images use `cover`; portraits taller than three times their width use `contain` to avoid excessive cropping.
- Both layouts use a 200px minimum column width. Text-only pages stay in the same collection and source order. Resizing or loading images can rebalance Gallery columns.

Image variants come from `getAssetImageProps`, including portrait-aware responsive candidates. Width and height are omitted while the asset job is processing; the browser uses the image's intrinsic dimensions instead. The Gallery fit rule uses asset metadata and works during server rendering, without measuring or repositioning cards in JavaScript.
