import { defineSlotRecipe } from "@pandacss/dev";

export const inputGroup = defineSlotRecipe({
  className: "input-group",
  slots: ["root", "element"],
  base: {
    root: {
      position: "relative",
      width: "full",
      "& > input:not(:first-child)": {
        ps: "var(--input-group-element-width)!",
      },
      "& > input:not(:last-child)": {
        pe: "var(--input-group-element-width)!",
      },
    },
    element: {
      alignItems: "center",
      color: "text.subtle",
      display: "flex",
      fontSize: "var(--input-group-element-font-size)",
      height: "full",
      justifyContent: "center",
      lineHeight: "1",
      minWidth: "var(--input-group-element-width)",
      paddingInline: "var(--input-group-element-padding)",
      position: "absolute",
      whiteSpace: "nowrap",
      zIndex: "2",
      _icon: {
        boxSize: "var(--input-group-icon-size)",
        color: "text.muted",
      },
    },
  },
  defaultVariants: {
    size: "sm",
  },
  variants: {
    size: {
      sm: {
        root: {
          "--input-group-element-font-size": "fontSizes.xs",
          "--input-group-element-padding": "spacing.2",
          "--input-group-element-width": "sizes.8",
          "--input-group-icon-size": "sizes.4",
        },
      },
      md: {
        root: {
          "--input-group-element-font-size": "fontSizes.sm",
          "--input-group-element-padding": "spacing.2.5",
          "--input-group-element-width": "sizes.10",
          "--input-group-icon-size": "sizes.4.5",
        },
      },
      lg: {
        root: {
          "--input-group-element-font-size": "fontSizes.md",
          "--input-group-element-padding": "spacing.3",
          "--input-group-element-width": "sizes.11",
          "--input-group-icon-size": "sizes.5",
        },
      },
    },
  },
});
