import { getContrast } from "polished";
import { describe, expect, test } from "vitest";

import { getColourVariants, getReadableTextColour } from "./colour";

describe("getColourVariants", () => {
  test("builds the documented twelve-step accent contract from hue", () => {
    const variants = getColourVariants("hsl(250, 15%, 80%)");

    expect({
      light1: variants["--accent-colour-light-1"],
      light9: variants["--accent-colour-light-9"],
      light11: variants["--accent-colour-light-11"],
      light12: variants["--accent-colour-light-12"],
      dark1: variants["--accent-colour-dark-1"],
      dark9: variants["--accent-colour-dark-9"],
      dark11: variants["--accent-colour-dark-11"],
      dark12: variants["--accent-colour-dark-12"],
      accent: variants["--sd-color-accent"],
      emphasized: variants["--sd-color-accent-emphasized"],
      text: variants["--sd-color-accent-text"],
      focus: variants["--sd-color-focus-ring"],
    }).toEqual({
      light1: "hsl(250deg 50% 99.2%)",
      light9: "hsl(250deg 51% 54.1%)",
      light11: "hsl(250deg 45% 49%)",
      light12: "hsl(250deg 50% 25.1%)",
      dark1: "hsl(250deg 23% 8.6%)",
      dark9: "hsl(250deg 51% 54.1%)",
      dark11: "hsl(250deg 100% 80.8%)",
      dark12: "hsl(250deg 77% 91.6%)",
      accent: "var(--accent-colour)",
      emphasized: "var(--sd-color-accent-10)",
      text: "var(--sd-color-accent-11)",
      focus: "var(--sd-color-accent-8)",
    });

    expect(
      Object.keys(variants).filter((key) =>
        key.match(/^--sd-color-accent-\d+$/),
      ),
    ).toHaveLength(12);
  });

  test("ignores input saturation and lightness", () => {
    expect(getColourVariants("hsl(250, 100%, 20%)")).toEqual(
      getColourVariants("hsl(250, 15%, 80%)"),
    );
  });

  test("extracts hue from arbitrary API colours and falls back safely", () => {
    expect(
      getColourVariants("#6855b8")["--accent-colour-light-9"],
    ).toBe("hsl(251.52deg 51% 54.1%)");
    expect(
      getColourVariants("not-a-colour")["--accent-colour-light-9"],
    ).toBe("hsl(156.99deg 51% 54.1%)");
  });

  test.each([0, 30, 60, 120, 180, 240, 300])(
    "keeps hue %d text steps readable in both colour schemes",
    (hue) => {
      const variants = getColourVariants(`hsl(${hue}, 15%, 80%)`);

      expect(
        getContrast(
          variants["--accent-colour-light-11"]!,
          variants["--accent-colour-light-2"]!,
        ),
      ).toBeGreaterThanOrEqual(4.5);
      expect(
        getContrast(
          variants["--accent-colour-light-12"]!,
          variants["--accent-colour-light-2"]!,
        ),
      ).toBeGreaterThanOrEqual(7);
      expect(
        getContrast(
          variants["--accent-colour-dark-11"]!,
          variants["--accent-colour-dark-2"]!,
        ),
      ).toBeGreaterThanOrEqual(4.5);
      expect(
        getContrast(
          variants["--accent-colour-dark-12"]!,
          variants["--accent-colour-dark-2"]!,
        ),
      ).toBeGreaterThanOrEqual(7);
      expect(
        getContrast(
          variants["--accent-colour-light-9"]!,
          variants["--accent-colour-contrast"]!,
        ),
      ).toBeGreaterThanOrEqual(4.5);
    },
  );
});

describe("getReadableTextColour", () => {
  test.each(["#b2ffd8", "#fbdfff", "#ffd3d8", "#202020", "#777777"])(
    "returns an accessible foreground for %s",
    (background) => {
      const foreground = getReadableTextColour(background);

      expect(getContrast(background, foreground)).toBeGreaterThanOrEqual(4.5);
    },
  );

  test("falls back safely for invalid colours", () => {
    expect(getReadableTextColour("not-a-colour")).toBe("#303030");
  });
});
