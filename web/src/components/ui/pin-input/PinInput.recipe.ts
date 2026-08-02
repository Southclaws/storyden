import { pinInputAnatomy } from "@ark-ui/react";
import { defineSlotRecipe } from "@pandacss/dev";

export const pinInput = defineSlotRecipe({
  className: "pinInput",
  slots: pinInputAnatomy.keys(),
  base: {
    root: {
      display: "flex",
      flexDirection: "column",
      gap: "1.5",
    },
    control: {
      display: "flex",
    },
    label: {
      color: "text.default",
      fontWeight: "medium",
    },
    input: {
      px: "0!",
      textAlign: "center",
    },
  },
  defaultVariants: {
    size: "sm",
  },
  variants: {
    size: {
      sm: {
        control: {
          gap: "1.5",
        },
        label: {
          fontSize: "xs",
          lineHeight: "1.125rem",
        },
        input: {
          minWidth: "6!",
          width: "6",
        },
      },
      md: {
        control: {
          gap: "2",
        },
        label: {
          fontSize: "sm",
          lineHeight: "1.25rem",
        },
        input: {
          minWidth: "8!",
          width: "8",
        },
      },
      lg: {
        control: {
          gap: "2.5",
        },
        label: {
          fontSize: "md",
          lineHeight: "1.5rem",
        },
        input: {
          minWidth: "10!",
          width: "10",
        },
      },
    },
  },
});
