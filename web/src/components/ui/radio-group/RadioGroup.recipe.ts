import { radioGroupAnatomy } from "@ark-ui/react";
import { defineSlotRecipe } from "@pandacss/dev";

export const radioGroup = defineSlotRecipe({
  className: "radioGroup",
  slots: radioGroupAnatomy.keys(),
  base: {
    root: {
      colorPalette: "accent",
      display: "flex",
      flexDirection: {
        _vertical: "column",
        _horizontal: "row",
      },
    },
    itemControl: {
      background: "transparent",
      borderColor: "border.default",
      borderRadius: "full",
      borderWidth: "1px",
      display: "flex",
      alignItems: "center",
      justifyContent: "center",
      transitionDuration: "normal",
      transitionProperty: "background, border-color",
      transitionTimingFunction: "default",
      _before: {
        backgroundColor: "accent.default",
        borderRadius: "full",
        content: "''",
        display: "block",
        transform: "scale(0)",
        transitionDuration: "fast",
        transitionProperty: "transform",
        transitionTimingFunction: "default",
      },
      _hover: {
        background: "background.controlHover",
      },
      _checked: {
        borderColor: "accent.default",
        _hover: {
          background: "transparent",
        },
        _before: {
          transform: "scale(1)",
        },
      },
      _disabled: {
        borderColor: "border.disabled",
        color: "text.disabled",
        _hover: {
          bg: "initial",
          color: "text.disabled",
        },
      },
    },
    item: {
      alignItems: "center",
      cursor: "pointer",
      display: "flex",
      _disabled: {
        cursor: "not-allowed",
      },
    },
    itemText: {
      color: "text.default",
      fontWeight: "medium",
      _disabled: {
        color: "text.disabled",
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
          gap: {
            _vertical: "2",
            _horizontal: "4",
          },
        },
        item: {
          gap: "1.5",
        },
        itemControl: {
          width: "4",
          height: "4",
          _before: {
            width: "2",
            height: "2",
          },
        },
        itemText: {
          fontSize: "xs",
          lineHeight: "1.125rem",
        },
      },
      md: {
        root: {
          gap: {
            _vertical: "3",
            _horizontal: "5",
          },
        },
        item: {
          gap: "2",
        },
        itemControl: {
          width: "5",
          height: "5",
          _before: {
            width: "2.5",
            height: "2.5",
          },
        },
        itemText: {
          fontSize: "sm",
          lineHeight: "1.25rem",
        },
      },
      lg: {
        root: {
          gap: {
            _vertical: "4",
            _horizontal: "6",
          },
        },
        item: {
          gap: "2.5",
        },
        itemControl: {
          width: "6",
          height: "6",
          _before: {
            width: "3",
            height: "3",
          },
        },
        itemText: {
          fontSize: "md",
          lineHeight: "1.5rem",
        },
      },
    },
  },
});
