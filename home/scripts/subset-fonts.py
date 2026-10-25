# /// script
# requires-python = ">=3.10"
# dependencies = ["fonttools[woff]==4.67.0", "brotli==1.2.0"]
# ///
"""Regenerate the English site's local fonts: uv run home/scripts/subset-fonts.py.

Keep the original fonts in static/ as the source of truth. The subset includes
Latin, extended Latin, combining accents, punctuation, currencies and arrows.
Other scripts continue to use the site's system-font fallback.
"""

from pathlib import Path
from unicodedata import normalize

from fontTools import subset
from fontTools.ttLib import TTFont

FONT_DIR = Path(__file__).resolve().parents[1] / "src/fonts/static"
UNICODES = (
    "U+0000-024F,U+02BB-02BC,U+02C6,U+02DA,U+02DC,U+0300-036F,"
    "U+2000-206F,U+20A0-20CF,U+2100-214F,"
    "U+2190-21FF,U+2212,U+2215,U+FEFF,U+FFFD"
)


def generate():
    output = FONT_DIR / "subsets"
    output.mkdir(exist_ok=True)
    total_before = total_after = 0
    for source in sorted(FONT_DIR.iterdir()):
        if source.suffix != ".woff2" and source.name != "Gorton-Normal-120.otf":
            continue
        font = TTFont(source, recalcTimestamp=False)
        original_cmap = font.getBestCmap()
        characters = set(subset.parse_unicodes(UNICODES)) & original_cmap.keys()
        characters |= {
            cp
            for cp in original_cmap
            if all(ord(char) in characters for char in normalize("NFD", chr(cp)))
        }
        options = subset.Options()
        options.notdef_outline = True
        options.recalc_timestamp = False
        subsetter = subset.Subsetter(options=options)
        subsetter.populate(unicodes=characters)
        subsetter.subset(font)

        # Preserve mappings for retained composite glyphs. HarfBuzz composition
        # and FreeType auto-hinting use these even when the glyph is not requested
        # directly; see the established web/scripts/subset-fonts.py workflow.
        glyphs = set(font.getGlyphOrder())
        for table in font["cmap"].tables:
            if table.isUnicode():
                table.cmap.update({
                    cp: glyph
                    for cp, glyph in original_cmap.items()
                    if glyph in glyphs and (cp <= 0xFFFF or table.format in (12, 13))
                })
        font.flavor = "woff2"
        destination = output / f"{source.stem}.woff2"
        font.save(destination)
        font.close()
        with TTFont(destination) as generated:
            assert characters <= generated.getBestCmap().keys(), source.name
        total_before += source.stat().st_size
        total_after += destination.stat().st_size
    print(f"Font subsets: {total_before:,} -> {total_after:,} bytes")


if __name__ == "__main__":
    generate()
