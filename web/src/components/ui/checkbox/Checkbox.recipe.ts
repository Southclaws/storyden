import { checkboxAnatomy } from "@ark-ui/react";
import { defineSlotRecipe } from "@pandacss/dev";

export const checkbox = defineSlotRecipe({
  className: "checkbox",
  slots: checkboxAnatomy.keys(),
  base: {
    root: {
      alignItems: "center",
      display: "flex",
    },
    label: {
      color: "text.default",
      fontWeight: "medium",
    },
    control: {
      alignItems: "center",
      borderColor: "border.default",
      borderWidth: "1px",
      color: "colorPalette.fg",
      cursor: "pointer",
      display: "flex",
      justifyContent: "center",
      transitionDuration: "normal",
      transitionProperty: "border-color, background",
      transitionTimingFunction: "default",
      _hover: {
        background: "control.hoverBackground",
      },
      _checked: {
        background: "colorPalette.default",
        borderColor: "colorPalette.default",
        _hover: {
          background: "colorPalette.default",
        },
      },
      _indeterminate: {
        background: "colorPalette.default",
        borderColor: "colorPalette.default",
        _hover: {
          background: "colorPalette.default",
        },
      },
      "&:has(+ :focus-visible)": {
        outlineOffset: "2px",
        outline: "2px solid",
        outlineColor: "border.default",
        _checked: {
          outlineColor: "colorPalette.default",
        },
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
          gap: "1.5",
        },
        control: {
          width: "4",
          height: "4",
          borderRadius: "xs",
          "& svg": {
            width: "3",
            height: "3",
          },
        },
        label: {
          textStyle: "xs",
        },
      },
      md: {
        root: {
          gap: "2",
        },
        control: {
          width: "5",
          height: "5",
          borderRadius: "xs",
          "& svg": {
            width: "3.5",
            height: "3.5",
          },
        },
        label: {
          textStyle: "sm",
        },
      },
      lg: {
        root: {
          gap: "2.5",
        },
        control: {
          width: "6",
          height: "6",
          borderRadius: "xs",
          "& svg": {
            width: "4",
            height: "4",
          },
        },
        label: {
          textStyle: "md",
        },
      },
    },
  },
});
