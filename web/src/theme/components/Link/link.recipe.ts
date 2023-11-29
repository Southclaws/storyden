import { defineRecipe } from "@pandacss/dev";

export const link = defineRecipe({
  className: "link",
  base: {
    alignItems: "center",
    appearance: "none",
    borderRadius: "lg",
    boxShadow: "xs",
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
        // TODO: Add back when chakra is removed.
        // borderWidth: "1px",
        // borderColor: "border.default",
        background: "blackAlpha.100",
        _hover: {
          background: "gray.100",
          boxShadow: "md",
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
          boxShadow: "md",
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
          boxShadow: "md",
        },
        _active: {
          backgroundColor: "gray.400",
        },
      },
      ghost: {
        _hover: {
          backgroundColor: "gray.200",
          boxShadow: "md",
        },
        _active: {
          backgroundColor: "gray.300",
        },
      },
    },

    size: {
      xs: {
        h: "6",
        minW: "8",
        textStyle: "xs",
        px: "2",
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
        px: "2",
        gap: "2",
        "& svg": {
          width: "5",
          height: "5",
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
