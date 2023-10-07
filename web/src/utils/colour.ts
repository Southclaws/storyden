import Color from "colorjs.io";
import { readableColor } from "polished";

export const FALLBACK_COLOUR = "#27b981";

export function getColourVariants(
  colour: string,
  contrast: number = 1,
): Record<string, string> {
  const c = parseColourWithFallback(colour);

  const hue = getHue(c);

  const rgb = c.to("srgb").toString({ format: "hex" });

  const textColour = readableColorWithFallback(rgb);

  const flatClampLightness: [number, number] = [71.8, 88.7];
  const flatClampChroma: [number, number] = [14, 3];

  const flatRamp = [1, 2, 3].reduce((o, i, _, a) => {
    const indices = a.length;
    const [minL, maxL] = flatClampLightness;
    const [minC, maxC] = flatClampChroma;

    const L = minL + ((maxL - minL) / indices) * i * 1.725 * contrast;
    const C = minC + ((maxC - minC) / indices) * i;

    const fill = `oklch(${L}% ${C}% ${hue}deg)`;

    const text = readableColorWithFallback(
      parseColourWithFallback(fill).to("srgb").toString({ format: "hex" }),
    );

    return {
      [`--accent-colour-flat-fill-${i}`]: fill,
      [`--accent-colour-flat-text-${i}`]: text,
      ...o,
    };
  }, {});

  const darkClampLightness: [number, number] = [93.7, 50.8];
  const darkClampChroma: [number, number] = [14, 7];

  const darkRamp = [1, 2, 3].reduce((o, i, _, a) => {
    const indices = a.length;
    const [minL, maxL] = darkClampLightness;
    const [minC, maxC] = darkClampChroma;

    const L = minL + ((maxL - minL) / indices) * i * 1.725 * contrast;
    const C = minC + ((maxC - minC) / indices) * i;

    const fill = `oklch(${L}% ${C}% ${hue}deg)`;

    const text = readableColorWithFallback(
      parseColourWithFallback(fill).to("srgb").toString({ format: "hex" }),
    );

    return {
      [`--accent-colour-dark-fill-${i}`]: fill,
      [`--accent-colour-dark-text-${i}`]: text,
      ...o,
    };
  }, {});

  return {
    "--text-colour": textColour,

    "--accent-colour": `oklch(80% 20% ${hue}deg)`,

    ...flatRamp,
    ...darkRamp,
  };
}

export function getColourAsHex(colour: string) {
  return parseColourWithFallback(colour).to("srgb").toString({ format: "hex" });
}

function parseColourWithFallback(colour: string) {
  try {
    return new Color(colour);
  } catch (e) {
    return new Color(FALLBACK_COLOUR);
  }
}

function getHue(c: Color) {
  const hue = c.oklch["h"];
  if (!hue) {
    return 0;
  }
  if (isNaN(hue)) {
    return 0;
  }
  return hue;
}

function readableColorWithFallback(rgb: string): string {
  try {
    return readableColor(rgb, "#E8ECEA", "#303030", true);
  } catch (e) {
    return "black";
  }
}
