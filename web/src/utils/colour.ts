import chroma from "chroma-js";
import { readableColor } from "polished";

export const FALLBACK_COLOUR = "#27b981";

type ScaleStop = readonly [saturation: number, lightness: number];

// The steps follow the same functional progression as the static Radix
// palettes used elsewhere in Storyden: backgrounds (1-2), controls (3-5),
// borders (6-8), solid fills (9-10), and text (11-12).
const lightScale: readonly ScaleStop[] = [
  [50, 99.2],
  [78, 98.2],
  [90, 96.3],
  [81, 93.7],
  [75, 90.6],
  [69, 86.3],
  [62, 80.6],
  [60, 73.5],
  [51, 54.1],
  [45, 50.2],
  [45, 49],
  [50, 25.1],
];

const darkScale: readonly ScaleStop[] = [
  [23, 8.6],
  [25, 11],
  [36, 17.1],
  [39, 22],
  [38, 26.1],
  [35, 31.2],
  [33, 38.4],
  [33, 50.4],
  [51, 54.1],
  [55, 58.8],
  [100, 80.8],
  [77, 91.6],
];

export function getColourVariants(colour: string): Record<string, string> {
  const hue = getHue(colour);
  const roundedHue = Math.round(hue * 100) / 100;

  const light = createScale(roundedHue, lightScale, "light");
  const dark = createScale(roundedHue, darkScale, "dark");
  const solid = light[8] ?? `hsl(${roundedHue}deg 51% 54.1%)`;
  const contrast = getReadableTextColour(parseColourWithFallback(solid).hex());

  const variants: Record<string, string> = {};

  light.forEach((value, index) => {
    const step = index + 1;
    const darkValue = dark[index];

    if (!darkValue) return;

    variants[`--accent-colour-light-${step}`] = value;
    variants[`--accent-colour-dark-${step}`] = darkValue;
    variants[`--sd-color-accent-${step}`] =
      `light-dark(var(--accent-colour-light-${step}), var(--accent-colour-dark-${step}))`;
  });

  return {
    ...variants,
    "--accent-colour":
      "light-dark(var(--accent-colour-light-9), var(--accent-colour-dark-9))",
    "--accent-colour-contrast": contrast,
    "--sd-color-accent": "var(--accent-colour)",
    "--sd-color-accent-emphasized": "var(--sd-color-accent-10)",
    "--sd-color-accent-foreground": "var(--accent-colour-contrast)",
    "--sd-color-accent-text": "var(--sd-color-accent-11)",
    "--sd-color-focus-ring": "var(--sd-color-accent-8)",
  };
}

function getHue(colour: string): number {
  const explicitHsl = colour
    .trim()
    .match(
      /^hsla?\(\s*([-+]?(?:\d+(?:\.\d+)?|\.\d+))\s*(deg|grad|rad|turn)?(?:\s|,)/i,
    );

  if (explicitHsl) {
    const value = Number(explicitHsl[1]);
    const unit = explicitHsl[2]?.toLowerCase();
    const degrees =
      unit === "turn"
        ? value * 360
        : unit === "rad"
          ? (value * 180) / Math.PI
          : unit === "grad"
            ? value * 0.9
            : value;

    return ((degrees % 360) + 360) % 360;
  }

  try {
    const hue = chroma(colour).hsl()[0];

    return Number.isFinite(hue) ? hue : 0;
  } catch {
    return chroma(FALLBACK_COLOUR).hsl()[0];
  }
}

function createScale(
  hue: number,
  stops: readonly ScaleStop[],
  appearance: "light" | "dark",
): string[] {
  const scale = stops.map(
    ([saturation, lightness]) => `hsl(${hue}deg ${saturation}% ${lightness}%)`,
  );
  const background = scale[1];
  const direction = appearance === "light" ? -1 : 1;

  if (!background) return scale;

  // Radix reserves steps 11 and 12 for readable text. Fixed HSL lightness
  // cannot satisfy every hue, so move only those two stops until they meet
  // Storyden's WCAG contrast floor against the scale's subtle background.
  const textStops = [
    { index: 10, minimum: 4.5 },
    { index: 11, minimum: 7 },
  ] as const;

  textStops.forEach(({ index, minimum }) => {
    const stop = stops[index];
    const value = scale[index];

    if (!stop || !value) return;

    let lightness = stop[1];

    while (
      chroma.contrast(scale[index] ?? value, background) < minimum &&
      lightness > 0 &&
      lightness < 100
    ) {
      lightness += direction;
      scale[index] =
        `hsl(${hue}deg ${stop[0]}% ${Math.round(lightness * 10) / 10}%)`;
    }
  });

  return scale;
}

export function getColourAsHex(colour: string) {
  return parseColourWithFallback(colour).hex();
}

function parseColourWithFallback(colour: string) {
  try {
    return chroma(colour);
  } catch (e) {
    return chroma(FALLBACK_COLOUR);
  }
}

export function getReadableTextColour(rgb: string): string {
  try {
    return readableColor(rgb, "#303030", "#E8ECEA", true);
  } catch (e) {
    return "#303030";
  }
}

export function deriveColour(s: string): string {
  const bytes = new TextEncoder().encode(s);

  const hash = bytes.reduce((r, b) => {
    const s = b * 42;
    const x = ((r + 1) * s) % 360;
    return x;
  }, 69);

  const hue = hash;

  return chroma(0.7226, 0.12, hue, "oklch").hex();
}
