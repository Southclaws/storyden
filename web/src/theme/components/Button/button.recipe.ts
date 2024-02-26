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
    color: "fg.default",
    alignItems: "center",
    appearance: "none",
    borderRadius: "lg",
    boxShadow: "xs",
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
    _hover: {
      boxShadow: "md",
    },
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
        backgroundColor: "bg.default",
        _hover: {
          background: "bg.subtle",
        },
        _active: {
          background: "bg.muted",
        },
        _disabled,
      },
      primary: {
        backgroundColor: "bg.accent",
        _hover: {
          backgroundColor: "bg.accent.subtle",
        },
        _active: {
          backgroundColor: "bg.accent.muted",
        },
        _disabled,
      },
      destructive: {
        backgroundColor: "bg.destructive",
        _hover: {
          backgroundColor: "bg.destructive.subtle", // TODO: Proper tokens for states on all colour tokens
        },
        _active: {
          backgroundColor: "bg.destructive.muted",
        },
        _disabled,
      },
      ghost: {
        _hover: {
          backgroundColor: "bg.muted",
        },
        _active: {
          backgroundColor: "bg.muted/80",
        },
      },
      blank: {
        backgroundColor: "none",
        px: "0",
        _disabled,
      },
    },

    size: {
      xs: {
        h: "6",
        minW: "8",
        textStyle: "xs",
        px: "1",
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
