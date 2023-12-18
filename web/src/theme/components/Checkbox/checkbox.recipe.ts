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
      color: "fg.default",
      fontWeight: "medium",
    },
    control: {
      alignItems: "center",
      borderColor: "border.default",
      borderWidth: "1px",
      color: "fg.muted",
      cursor: "pointer",
      display: "flex",
      justifyContent: "center",
      transitionDuration: "normal",
      transitionProperty: "border-color, background",
      transitionTimingFunction: "default",
      _hover: {
        background: "bg.subtle",
      },
      _checked: {
        background: "accent.100",
        borderColor: "border.accent",
        _hover: {
          background: "accent.50",
        },
      },
    },
  },
  defaultVariants: {
    size: "md",
  },
  variants: {
    size: {
      sm: {
        root: {
          gap: "2",
        },
        control: {
          width: "4",
          height: "4",
          borderRadius: "sm",
          "& svg": {
            width: "3",
            height: "3",
          },
        },
        label: {
          textStyle: "sm",
        },
      },
      md: {
        root: {
          gap: "3",
        },
        control: {
          width: "5",
          height: "5",
          borderRadius: "md",
          "& svg": {
            width: "3.5",
            height: "3.5",
          },
        },
        label: {
          textStyle: "md",
        },
      },
      lg: {
        root: {
          gap: "4",
        },
        control: {
          width: "6",
          height: "6",
          borderRadius: "lg",
          "& svg": {
            width: "4",
            height: "4",
          },
        },
        label: {
          textStyle: "lg",
        },
      },
    },
  },
});
