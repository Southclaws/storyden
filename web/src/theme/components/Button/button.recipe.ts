import { defineRecipe } from "@pandacss/dev";

const _disabled = {
  background: "bg.disabled",
  borderColor: "border.disabled",
  color: "fg.disabled",
  cursor: "not-allowed",
  _hover: {
    background: "bg.disabled",
    borderColor: "border.disabled",
    color: "fg.disabled",
  },
};

export const button = defineRecipe({
  className: "button",
  base: {
    alignItems: "center",
    appearance: "none",
    borderRadius: "lg",
    cursor: "pointer",
    display: "inline-flex",
    fontWeight: "semibold",
    minWidth: "0",
    width: "min-content",
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
  compoundVariants: [
    {
      // Blank buttons override some size settings which only apply to buttons
      // with borders and fills etc.
      kind: "blank",
      css: {
        px: "0",
        color: "accent.text.50",
      },
    },
  ],
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
        _disabled,
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
        _disabled,
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
        _disabled,
      },
      destructive: {
        backgroundColor: "rose.600",
        color: "white",
        _hover: {
          backgroundColor: "rose.500",
        },
        _active: {
          backgroundColor: "rose.700",
        },
        _disabled,
      },
      blank: {
        backgroundColor: "none",
        px: "0",
        _disabled,
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
          width: "4",
          height: "4",
        },
      },
      sm: {
        h: "9",
        minW: "9",
        textStyle: "sm",
        px: "3.5",
        gap: "2",
        "& svg": {
          width: "4",
          height: "4",
        },
      },
      md: {
        h: "10",
        minW: "10",
        textStyle: "sm",
        px: "4",
        gap: "2",
        "& svg": {
          width: "5",
          height: "5",
        },
      },
      lg: {
        h: "11",
        minW: "11",
        textStyle: "md",
        px: "4.5",
        gap: "2",
        "& svg": {
          width: "5",
          height: "5",
        },
      },
      xl: {
        h: "12",
        minW: "12",
        textStyle: "md",
        px: "5",
        gap: "2.5",
        "& svg": {
          width: "5",
          height: "5",
        },
      },
      "2xl": {
        h: "16",
        minW: "16",
        textStyle: "lg",
        px: "7",
        gap: "3",
        "& svg": {
          width: "6",
          height: "6",
        },
      },
    },
  },
});
