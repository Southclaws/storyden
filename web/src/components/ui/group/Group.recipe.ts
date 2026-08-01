import { defineRecipe } from "@pandacss/dev";

export const group = defineRecipe({
  className: "ui-group",
  jsx: ["Group", "ButtonGroup"],
  base: {
    display: "inline-flex",
    position: "relative",
    gap: "2",
    "& > *": {
      _focusVisible: {
        zIndex: 1,
      },
    },
  },
  defaultVariants: {
    orientation: "horizontal",
  },
  variants: {
    orientation: {
      horizontal: {
        flexDirection: "row",
      },
      vertical: {
        flexDirection: "column",
      },
    },
    attached: {
      true: {
        gap: "0",
      },
    },
    grow: {
      true: {
        display: "flex",
        "& > *": {
          flex: 1,
        },
      },
    },
  },
  compoundVariants: [
    {
      orientation: "horizontal",
      attached: true,
      css: {
        "& > *:first-child:not(:only-child)": {
          borderEndRadius: "0",
          marginEnd: "-1px",
        },
        "& > *:last-child:not(:only-child)": {
          borderStartRadius: "0",
        },
        "& > *:not(:first-child):not(:last-child)": {
          borderRadius: "0",
          marginEnd: "-1px",
        },
      },
    },
    {
      orientation: "vertical",
      attached: true,
      css: {
        "& > *:first-child:not(:only-child)": {
          borderBottomRadius: "0",
          marginBottom: "-1px",
        },
        "& > *:last-child:not(:only-child)": {
          borderTopRadius: "0",
        },
        "& > *:not(:first-child):not(:last-child)": {
          borderRadius: "0",
          marginBottom: "-1px",
        },
      },
    },
  ],
});
