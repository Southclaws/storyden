import Color from "colorjs.io";
import { readableColor } from "polished";

type Colours = {
  "--text-colour": string;

  "--accent-colour": string;
  "--accent-colour-muted": string;

  // For browsers without OKLCH
  "--accent-colour-fallback": string;
  "--accent-colour-muted-fallback": string;
};

export function getColourVariants(colour: string): Colours {
  const c = new Color(colour);

  const hue = c.oklch["h"];

  const rgb = c.to("srgb").toString({ format: "rgb" });

  const textColour = readableColor(rgb, "#E8ECEA", "#303030", false);

  return {
    "--text-colour": textColour,

    "--accent-colour": `oklch(80% 0.2 ${hue}deg)`,
    "--accent-colour-muted": `oklch(90% 0.1 ${hue}deg)`,

    "--accent-colour-fallback": `hsl(${hue} 100% 43%)`,
    "--accent-colour-muted-fallback": `hsl(${hue} 24% 63%)`,
  };
}
