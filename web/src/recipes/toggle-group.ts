import { toggleGroupAnatomy } from "@ark-ui/react";
import { defineSlotRecipe } from "@pandacss/dev";

export const toggleGroup = defineSlotRecipe({
  className: "toggleGroup",
  slots: toggleGroupAnatomy.keys(),
  base: {
    root: {
      display: "flex",
      overflow: "hidden",
      position: "relative",
      _vertical: {
        flexDirection: "column",
      },
    },
    item: {
      alignItems: "center",
      appearance: "none",
      cursor: "pointer",
      color: "fg.muted",
      display: "inline-flex",
      fontWeight: "normal",
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
      _on: {
        background: "bg.emphasized",
        color: "fg.emphasized",
        _hover: {
          background: "bg.emphasized",
          color: "fg.emphasized",
        },
      },
      _hover: {
        background: "bg.subtle",
      },
      _disabled: {
        borderColor: "border.disabled",
        color: "fg.disabled",
        cursor: "not-allowed",
        _hover: {
          background: "transparent",
          borderColor: "border.disabled",
          color: "fg.disabled",
        },
      },
    },
  },
  defaultVariants: {
    size: "md",
    variant: "outline",
  },
  variants: {
    variant: {
      outline: {
        root: {
          borderWidth: "1px",
          borderRadius: "l2",
          borderColor: "border.default",
          _horizontal: {
            divideX: "1px",
          },
          _vertical: {
            divideY: "1px",
          },
        },
        item: {
          borderColor: "border.default",
          _focusVisible: {
            zIndex: 1,
            outlineOffset: "-2px",
            outline: "2px solid",
            outlineColor: "border.outline",
          },
        },
      },
      ghost: {
        root: {
          gap: "1",
        },
        item: {
          borderRadius: "l2",
          _focusVisible: {
            outlineOffset: "2px",
            outline: "2px solid",
            outlineColor: "border.outline",
          },
        },
      },
    },
    size: {
      xs: {
        item: {
          px: "1",
          h: "6",
          minW: "6",
          textStyle: "xs",
          gap: "1",
          "& svg": {
            width: "4",
            height: "4",
          },
        },
      },
      sm: {
        item: {
          px: "2",
          h: "8",
          minW: "8",
          textStyle: "sm",
          gap: "2",
          "& svg": {
            width: "4",
            height: "4",
          },
        },
      },
      md: {
        item: {
          px: "2",
          h: "10",
          minW: "10",
          textStyle: "sm",
          gap: "2",
          "& svg": {
            width: "5",
            height: "5",
          },
        },
      },
      lg: {
        item: {
          px: "2",
          h: "11",
          minW: "11",
          textStyle: "md",
          gap: "2",
          "& svg": {
            width: "5",
            height: "5",
          },
        },
      },
    },
  },
  compoundVariants: [
    {
      variant: "outline",
      size: "xs",
      css: {
        item: {
          h: "{sizes.5.5}",
        },
      },
    },
    {
      variant: "outline",
      size: "sm",
      css: {
        item: {
          h: "7.5",
        },
      },
    },
    {
      variant: "outline",
      size: "md",
      css: {
        item: {
          h: "9.5",
        },
      },
    },
    {
      variant: "outline",
      size: "lg",
      css: {
        item: {
          h: "10.5",
        },
      },
    },
  ],
});
