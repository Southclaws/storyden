import { range } from "lodash";
import { map } from "lodash/fp";

// Dynamic lightness for theming
const L = "80%";

const C = "0.15";

const lch = (hue: number) => `oklch(${L} ${C} ${hue})`;

const stops = map(lch)(range(0, 361, 10));

export const conicGradient = {
  value: `
conic-gradient(
    ${stops.join(",\n")}
)
`,
};
