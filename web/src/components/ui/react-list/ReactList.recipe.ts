import { defineSlotRecipe } from "@pandacss/dev";

export const reactList = defineSlotRecipe({
  className: "reactList",
  slots: ["root", "reaction", "count", "picker"],
  jsx: ["ReactList"],
  base: {
    root: {
      display: "flex",
      flexWrap: "wrap",
      gap: "1",
    },
    reaction: {
      borderRadius: "md",
      color: "text.muted",
      fontSize: "sm",
      fontVariantNumeric: "tabular-nums",
      gap: "0",
      height: "7",
      paddingX: "2",
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
        boxShadow: "inset 0 0 0 1px token(colors.interactive.selected.border)",
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
      marginLeft: "2",
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
      borderRadius: "md",
      color: "text.subtle",
      fontSize: "sm",
      height: "7",
      paddingX: "2",
    },
  },
});
