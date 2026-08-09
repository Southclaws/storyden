import { getContrast } from "polished";
import { describe, expect, it } from "vitest";

import amber from "./amber";
import blue from "./blue";
import green from "./green";
import neutral from "./neutral";
import orange from "./orange";
import pink from "./pink";
import red from "./red";
import slate from "./slate";
import tomato from "./tomato";

type Theme = "light" | "dark";
type Condition = "_osLight" | "_osDark";

type ColourPalette = {
  name: string;
  tokens: Record<Theme, Record<string, { value: string }>>;
  semanticTokens: Record<
    "solid" | "solidForeground" | "text",
    { value: Record<Condition, string> }
  >;
};

const palettes = [
  amber,
  blue,
  green,
  neutral,
  orange,
  pink,
  red,
  slate,
  tomato,
] as unknown as ColourPalette[];

const canvases = {
  light: neutral.tokens.light["1"].value,
  dark: slate.tokens.dark["1"].value,
};

function resolveSemanticColour(
  palette: ColourPalette,
  token: "solid" | "solidForeground" | "text",
  theme: Theme,
) {
  const condition = theme === "light" ? "_osLight" : "_osDark";
  const value = palette.semanticTokens[token].value[condition];

  if (value === "black" || value === "white") {
    return value;
  }

  const reference = value.match(/\{colors\.[^.]+\.(light|dark)\.([^.}]+)\}/);

  if (!reference) {
    throw new Error(`Cannot resolve ${palette.name}.${token} from ${value}`);
  }

  const [, referencedTheme, step] = reference;

  if (!referencedTheme || !step) {
    throw new Error(`Invalid colour reference ${value}`);
  }

  return palette.tokens[referencedTheme as Theme][step]?.value;
}

describe.each(palettes)("$name colour palette", (palette) => {
  it.each(["light", "dark"] as const)(
    "keeps %s semantic text readable on canvas and tinted controls",
    (theme) => {
      const text = resolveSemanticColour(palette, "text", theme);
      const tint = palette.tokens[theme]["3"]?.value;

      expect(getContrast(text!, canvases[theme])).toBeGreaterThanOrEqual(4.5);
      expect(getContrast(text!, tint!)).toBeGreaterThanOrEqual(4.5);
    },
  );

  it.each(["light", "dark"] as const)(
    "keeps %s solid foreground readable",
    (theme) => {
      const background = resolveSemanticColour(palette, "solid", theme);
      const foreground = resolveSemanticColour(
        palette,
        "solidForeground",
        theme,
      );

      expect(getContrast(foreground!, background!)).toBeGreaterThanOrEqual(4.5);
    },
  );
});

describe("destructive solid colours", () => {
  it.each(["light", "dark"] as const)(
    "uses a light foreground in %s mode",
    (theme) => {
      expect(resolveSemanticColour(red, "solidForeground", theme)).toBe(
        "white",
      );
    },
  );
});
