import { defineRecipe } from "@pandacss/dev";

export const link = defineRecipe({
  className: "link",
  base: {
    alignItems: "center",
    appearance: "none",
    borderRadius: "lg",
    cursor: "pointer",
    display: "inline-flex",
    fontWeight: "semibold",
    minWidth: "0",
    justifyContent: "center",
    outline: "none",
    position: "relative",
    transitionDuration: "normal",
    transitionProperty: "background, border-color, color, box-shadow",
    transitionTimingFunction: "default",
    userSelect: "none",
    verticalAlign: "middle",
    whiteSpace: "nowrap",
  },
  defaultVariants: {
    kind: "neutral",
    size: "md",
  },
  variants: {
    kind: {
      neutral: {
        borderWidth: "1px",
        borderColor: "border.default",
        _hover: {
          background: "gray.100",
        },
        _focusVisible: {
          outlineOffset: "2px",
          outline: "2px solid",
          outlineColor: "border.outline",
        },
        _active: {
          backgroundColor: "gray.200",
        },
      },
      primary: {
        backgroundColor: "accent.500",
        color: "accent.text.500",
        _hover: {
          backgroundColor: "accent.400",
        },
        _focusVisible: {
          outlineOffset: "2px",
          outline: "2px solid",
          outlineColor: "border.outline",
        },
        _active: {
          backgroundColor: "accent.600",
        },
      },
      secondary: {
        backgroundColor: "gray.200",
        color: "accent.text.200",
        _hover: {
          backgroundColor: "gray.300",
        },
        _active: {
          backgroundColor: "gray.400",
        },
      },
    },

    size: {
      xs: {
        h: "8",
        minW: "8",
        textStyle: "xs",
        px: "3",
        gap: "2",
        "& svg": {
          fontSize: "md",
          width: "3",
          height: "3",
        },
      },
      sm: {
        h: "9",
        minW: "9",
        textStyle: "sm",
        px: "3.5",
        gap: "2",
        "& svg": {
          width: "3",
          height: "3",
        },
      },
      md: {
        h: "10",
        minW: "10",
        textStyle: "sm",
        px: "4",
        gap: "2",
        "& svg": {
          width: "4",
          height: "4",
        },
      },
      lg: {
        h: "11",
        minW: "11",
        textStyle: "md",
        px: "4.5",
        gap: "2",
        "& svg": {
          width: "4",
          height: "4",
        },
      },
      xl: {
        h: "12",
        minW: "12",
        textStyle: "md",
        px: "4",
        gap: "2.5",
        "& svg": {
          width: "4",
          height: "4",
        },
      },
      "2xl": {
        h: "16",
        minW: "16",
        textStyle: "lg",
        px: "7",
        gap: "3",
        "& svg": {
          width: "5",
          height: "5",
        },
      },
    },
  },
});
