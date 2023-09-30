import Color from "colorjs.io";
import { readableColor } from "polished";

export const FALLBACK_COLOUR = "#27b981";

type Colours = {
  "--text-colour": string;

  "--accent-colour": string;
  "--accent-colour-muted": string;
  "--accent-colour-subtle": string;

  // For browsers without OKLCH
  "--accent-colour-fallback": string;
  "--accent-colour-muted-fallback": string;
  "--accent-colour-subtle-fallback": string;
};

export function getColourVariants(colour: string): Colours {
  const c = parseColourWithFallback(colour);

  const hue = getHue(c);

  const rgb = c.to("srgb").toString({ format: "hex" });

  const textColour = readableColorWithFallback(rgb);

  return {
    "--text-colour": textColour,

    "--accent-colour": `oklch(80% 0.2 ${hue}deg)`,
    "--accent-colour-muted": `oklch(90% 0.1 ${hue}deg)`,
    "--accent-colour-subtle": `oklch(100% 0.02 ${hue}deg)`,

    "--accent-colour-fallback": `hsl(${hue} 100% 43%)`,
    "--accent-colour-muted-fallback": `hsl(${hue} 24% 63%)`,
    "--accent-colour-subtle-fallback": `hsl(${hue} 100% 1%)`,
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
    return readableColor(rgb, "#E8ECEA", "#303030", false);
  } catch (e) {
    return "black";
  }
}
