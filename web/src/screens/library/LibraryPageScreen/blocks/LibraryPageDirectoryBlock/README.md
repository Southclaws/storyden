# Directory layouts

The page's Directory block stores one shared layout in its metadata. Table remains the default; Rows, Grid and Gallery use the same fields, controls and page sequence.

The display menu in the Directory toolbar and the block's Layout menu configure the draft. Existing page autosave persists these settings.

- Rows uses the Surface row shape, with compact text and an optional image on the right.
- Grid offers Minimal, Compact (Surface row), Editorial (Surface box), and Image overlay (Surface fill). Image overlay uses Editorial for entries without a cover.
- Gallery offers Minimal and Editorial. Editorial uses the asset's natural image ratio. Switching to Gallery temporarily renders other saved card styles as Editorial; switching back to Grid restores the saved style.
- Each custom property can be hidden or placed in the header, details, bottom-left badge, or bottom-right badge. Its placement is independent of Table visibility: visible properties remain ordinary columns in Table. Fixed fields keep their existing show/hide controls.
- All card presentations retain the Directory's sortable item IDs and page-menu drag handles during direct editing. Property values, names, descriptions and links remain editable.
- Rich Card presentations offer a Rich content toggle under Options in the display menu. It passes the page's HTML to the Surface preview and its existing Show more disclosure. Content is hidden by default, configured independently of the description, and read-only in directory previews; edit the page itself to change its rich content.

`card_style` is optional, so existing directories retain Minimal. `show_content` is optional and defaults to false. Column `placement` is optional and defaults to details. These options do not change the property's schema or stored value.

The interactive `Compositions/Library/Directory Layouts` story is the R&D playground for mixed image ratios, missing images, property placement, editing, and dragging. Its fixture edits stay in memory. Its page-router setup and navigation spies avoid the installed Storybook adapter's outdated Next App Router context.

## Minimal media behavior

- Grid uses aligned rows and a 4:3 image crop with `object-fit: cover`.
- Gallery uses balanced CSS columns, read down each column. Images keep their natural height between 32px and `min(70svh, 40rem)`. Ordinary images use `cover`; portraits taller than three times their width use `contain` to avoid excessive cropping.
- Both layouts use a 200px minimum column width. Text-only pages stay in the same collection and source order. Resizing or loading images can rebalance Gallery columns.

Image variants come from `getAssetImageProps`, including portrait-aware responsive candidates. Width and height are omitted while the asset job is processing; the browser uses the image's intrinsic dimensions instead. The Gallery fit rule uses asset metadata and works during server rendering, without measuring or repositioning cards in JavaScript.

## Card implementation

- `DirectoryRichCard.tsx` composes Surface slots; `DirectoryRichCardContent.tsx` owns their header, description, and link rendering.
- All directory cards use Surface. Its `minimal` shape owns the media-first layout, flat surface, and image placeholder styling. Responsive image sources and directory-specific fallback content are supplied by `DirectoryRichCard.tsx`.
- `DirectoryCardFields.tsx` shares property lists, value editors, and property-label menus. `DirectoryPropertyBadges.tsx` shares badge editors and footer layout.
- `directoryCard.ts` contains pure shape, aspect-ratio, visibility, badge-selection, and draft-overlay rules, covered by `directoryCard.test.ts`.
- `LibraryPageDirectoryBlockGrid.tsx` owns collection order and sortable items. `DirectoryCardControls.tsx` owns each item's drag handle and page menu.
