import Color from "colorjs.io";
import { readableColor } from "polished";

export const FALLBACK_COLOUR = "#27b981";

const ramp = [
  "50", // 1
  "100", // 2
  "200", // 3
  "300", // 4
  "400", // 5
  "500", // 6
  "600", // 7
  "700", // 8
  "800", // 9
  "900", // 10
];

const rampSize = 10;

export const flatClampL: [number, number] = [98.7, 51.8];
export const flatClampC: [number, number] = [2, 45];
export const flatContrast = 1.241;

export const darkClampL: [number, number] = [65, 20];
export const darkClampC: [number, number] = [10, 14];
export const darkContrast = 1.33;

export function getColourVariants(colour: string): Record<string, string> {
  const c = parseColourWithFallback(colour) as any;

  const hue = getHue(c);

  const rgb = c.to("srgb").toString({ format: "hex" });

  const textColour = readableColorWithFallback(rgb);

  const flatRamp = ramp.reduceRight((o, r, i) => {
    const [minL, maxL] = flatClampL;
    const [minC, maxC] = flatClampC;

    const L = minL + ((maxL - minL) / rampSize) * i * flatContrast;
    const C = minC + ((maxC - minC) / rampSize) * i;

    const fill = `oklch(${L}% ${C}% ${hue}deg)`;

    const text = readableColorWithFallback(
      parseColourWithFallback(fill).to("srgb").toString({ format: "hex" }),
    );

    return {
      [`--accent-colour-flat-fill-${r}`]: fill,
      [`--accent-colour-flat-text-${r}`]: text,
      ...o,
    };
  }, {});

  const darkRamp = ramp.reduceRight((o, r, i) => {
    const [minL, maxL] = darkClampL;
    const [minC, maxC] = darkClampC;

    const L = minL + ((maxL - minL) / rampSize) * i * darkContrast;
    const C = minC + ((maxC - minC) / rampSize) * i;

    const fill = `oklch(${L}% ${C}% ${hue}deg)`;

    const text = readableColorWithFallback(
      parseColourWithFallback(fill).to("srgb").toString({ format: "hex" }),
    );

    return {
      [`--accent-colour-dark-fill-${r}`]: fill,
      [`--accent-colour-dark-text-${r}`]: text,
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

function parseColourWithFallback(colour: string): any {
  try {
    return new (Color as any)(colour);
  } catch (e) {
    return new (Color as any)(FALLBACK_COLOUR);
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
    return readableColor(rgb, "#303030", "#E8ECEA", false);
  } catch (e) {
    return "black";
  }
}
