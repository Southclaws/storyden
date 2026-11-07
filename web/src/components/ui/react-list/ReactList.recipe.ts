import { defineSlotRecipe } from "@pandacss/dev";

import { button } from "../button/Button.recipe";

const reactionButtonStyles = {
  ...button.base,
  ...button.variants!.size.sm,
  ...button.variants!.variant.subtle,
  borderRadius: "md",
  fontSize: "sm",
  height: "7",
};

export const reactList = defineSlotRecipe({
  className: "reaction",
  slots: ["root", "trigger", "count", "picker"],
  jsx: ["ReactList"],
  base: {
    root: {
      display: "flex",
      flexWrap: "wrap",
      gap: "1",
    },
    trigger: {
      ...reactionButtonStyles,
      color: "text.muted",
      fontVariantNumeric: "tabular-nums",
      _disabled: {
        background: "colorPalette.3/80",
        color: "text.muted",
        cursor: "default",
        _hover: {
          background: "colorPalette.3/80",
          color: "text.muted",
        },
      },
      '&[aria-pressed="true"]': {
        background: "interactive.selected.surface",
        color: "interactive.selected.content",
        boxShadow:
          "[inset 0 0 0 1px token(colors.interactive.selected.border)]",
        _hover: {
          background: "interactive.emphasized.surface",
          color: "interactive.emphasized.content",
        },
        _disabled: {
          background: "interactive.selected.surface",
          color: "interactive.selected.content",
          _hover: {
            background: "interactive.selected.surface",
            color: "interactive.selected.content",
          },
        },
      },
    },
    count: {
      display: "inline-block",
      "--slide-distance": "10px",
      '&[data-direction="-1"]': { "--slide-distance": "-10px" },
      "@media (prefers-reduced-motion: no-preference)": {
        '&[data-animated="true"]': {
          animation: "slideFadeIn 250ms cubic-bezier(0.2, 0.8, 0.2, 1)",
        },
      },
    },
    picker: {
      ...reactionButtonStyles,
      color: "text.subtle",
    },
  },
});
