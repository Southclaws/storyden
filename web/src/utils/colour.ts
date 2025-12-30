import chroma from "chroma-js";
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

export const darkClampL: [number, number] = [35, 20];
export const darkClampC: [number, number] = [8, 0.5];
export const darkContrast = 1.33;

export function getColourVariants(colour: string): Record<string, string> {
  const c = parseColourWithFallback(colour);

  const hue = c.oklch()[2];

  const rgb = c.hex();

  const textColour = readableColorWithFallback(rgb);

  const flatRamp = ramp.reduceRight((o, r, i) => {
    const [minL, maxL] = flatClampL;
    const [minC, maxC] = flatClampC;

    const L = minL + ((maxL - minL) / rampSize) * i * flatContrast;
    const C = minC + ((maxC - minC) / rampSize) * i;

    const fill = `oklch(${L}% ${C}% ${hue}deg)`;

    const text = readableColorWithFallback(parseColourWithFallback(fill).hex());

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

    const text = readableColorWithFallback(parseColourWithFallback(fill).hex());

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
  return parseColourWithFallback(colour).hex();
}

function parseColourWithFallback(colour: string) {
  try {
    return chroma(colour);
  } catch (e) {
    return chroma(FALLBACK_COLOUR);
  }
}

function getHue(c) {
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
