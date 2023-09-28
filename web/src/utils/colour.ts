import Color from "colorjs.io";
import { readableColor } from "polished";

export const FALLBACK_COLOUR = "#27b981";

type Colours = {
  "--text-colour": string;

  "--accent-colour": string;
  "--accent-colour-muted": string;

  // For browsers without OKLCH
  "--accent-colour-fallback": string;
  "--accent-colour-muted-fallback": string;
};

export function getColourVariants(colour: string): Colours {
  const c = parseColourWithFallback(colour);

  const hue = c.oklch["h"];

  const rgb = c.to("srgb").toString({ format: "hex" });

  console.log({ colour, c, hue, rgb });

  const textColour = readableColorWithFallback(rgb);

  return {
    "--text-colour": textColour,

    "--accent-colour": `oklch(80% 0.2 ${hue}deg)`,
    "--accent-colour-muted": `oklch(90% 0.1 ${hue}deg)`,

    "--accent-colour-fallback": `hsl(${hue} 100% 43%)`,
    "--accent-colour-muted-fallback": `hsl(${hue} 24% 63%)`,
  };
}

export function getColourAsHex(colour: string) {
  return parseColourWithFallback(colour).to("srgb").toString({ format: "hex" });
}

function parseColourWithFallback(colour: string) {
  try {
    return new Color(colour);
  } catch (e) {
    console.log(e);
    return new Color(FALLBACK_COLOUR);
  }
}

function readableColorWithFallback(rgb: string): string {
  try {
    return readableColor(rgb, "#E8ECEA", "#303030", false);
  } catch (e) {
    console.log(e);
    return "black";
  }
}
