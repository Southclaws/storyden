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
    },
    count: {
      marginLeft: "2",
      display: "inline-block",
      "--slide-distance": "10px",
      '&[data-direction="-1"]': { "--slide-distance": "-10px" },
      '&[data-animated="true"]': {
        animation: "slideFadeIn 250ms cubic-bezier(0.2, 0.8, 0.2, 1)",
      },
      "@media (prefers-reduced-motion: reduce)": { animation: "none" },
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
