import { ark } from "@ark-ui/react/factory";
import chroma from "chroma-js";

import { styled } from "@/styled-system/jsx";
import { badge } from "@/styled-system/recipes";
import type { ComponentProps } from "@/styled-system/types";
import { getReadableTextColour } from "@/utils/colour";

export const StyledBadge = styled(ark.div, badge);
export type StyledBadgeProps = ComponentProps<typeof StyledBadge>;

export function badgeColours(hexColour: string) {
  const colour = chroma(hexColour);
  const hue = colour.lch()[2];

  const bg = chroma(0.95, 0.1, hue, "oklch").css();
  const border = chroma(0.85, 0.2, hue, "oklch").css();
  const fg = getReadableTextColour(bg);

  return { bg, border, fg };
}

export function badgeColourPalette(colourStyles: {
  bg: string;
  border: string;
  fg: string;
}) {
  const cssVars = colourStyles
    ? ({
        "--colors-color-palette-3": colourStyles.bg,
        "--colors-color-palette-6": colourStyles.border,
        "--colors-color-palette-12": colourStyles.fg,
      } as React.CSSProperties)
    : undefined;

  return cssVars;
}
