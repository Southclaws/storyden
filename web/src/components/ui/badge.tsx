import { ark } from "@ark-ui/react/factory";
import chroma from "chroma-js";

import { styled } from "@/styled-system/jsx";
import { badge } from "@/styled-system/recipes";
import type { ComponentProps } from "@/styled-system/types";

export type BadgeProps = ComponentProps<typeof Badge>;
export const Badge = styled(ark.div, badge);

export function badgeColours(hexColour: string) {
  const colour = chroma(hexColour);
  const hue = colour.lch()[2];

  const bg = chroma(0.95, 0.1, hue, "oklch").css();
  const border = chroma(0.85, 0.2, hue, "oklch").css();
  const fg = chroma(0.55, 0.2, hue, "oklch").css();

  return { bg, border, fg };
}

export function badgeColourPalette(colourStyles: {
  bg: string;
  border: string;
  fg: string;
}) {
  const cssVars = colourStyles
    ? ({
        "--colors-color-palette-fg": colourStyles.fg,
        "--colors-color-palette-border": colourStyles.border,
        "--colors-color-palette-bg": colourStyles.bg,
      } as React.CSSProperties)
    : undefined;

  return cssVars;
}
