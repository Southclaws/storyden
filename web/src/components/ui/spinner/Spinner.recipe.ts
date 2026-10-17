import { defineRecipe } from "@pandacss/dev";

export const spinner = defineRecipe({
  className: "spinner",
  jsx: ["Spinner"],
  base: {
    animation: "spin",
    borderColor: "background.controlTrack",
    borderRadius: "full",
    borderStyle: "solid",
    borderTopColor: "text.muted",
    flexShrink: "0",
    pointerEvents: "none",
  },
  variants: {
    size: {
      sm: { w: "4", h: "4", borderWidth: "[2px]" },
      md: { w: "5", h: "5", borderWidth: "[3px]" },
      lg: { w: "6", h: "6", borderWidth: "[5px]" },
    },
  },
  defaultVariants: { size: "md" },
});
