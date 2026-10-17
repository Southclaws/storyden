# Inter web fonts

The full variable WOFF2 files in `static/` are the source fonts. Regenerate the
checked-in subsets after replacing them:

```sh
uv run web/scripts/subset-fonts.py
pnpm --dir web exec oxfmt src/app/fonts/subsets.css
```

Ordinary builds use the checked-in files and do not need Python or `uv`.

`next/font/local` supplies the Latin faces and preloads only regular Latin.
`subsets.css` supplies extended Latin, Cyrillic, Greek and remaining symbols.
CSS `unicode-range` lets the browser select files from the actual text, including
mixed-language content, without locale detection. Italics load only when used.
All source characters, variable axes, hinting and applicable OpenType features
are retained. Combining marks are included in each subset for decomposed text.
The Latin face also covers their precomposed forms, including Vietnamese
accents. Retained glyphs keep their Unicode mappings for shaping and Linux
automatic hinting, even when CSS assigns those characters to another face.

Kerning cannot cross font-file boundaries. Letter shapes are preserved, but
some punctuation beside another script can have slightly different spacing.
For example, `«Тест»` is about 0.8 px wider at a 16 px font size.

Keep the Latin Unicode ranges in `index.ts` aligned with the ranges printed by
the generator, which include the `LATIN` definition and its composed accents.
Other ranges and CSS are generated from the fonts' character maps.
The complete original fonts are not imported by the application. CJK, Thai and
other scripts not supplied by Inter continue to use the existing system fallback.

The TTF files are separate inputs for Open Graph image generation.
