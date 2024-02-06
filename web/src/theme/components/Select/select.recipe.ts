import { selectAnatomy } from "@ark-ui/anatomy";
import { defineSlotRecipe } from "@pandacss/dev";

export const select = defineSlotRecipe({
  className: "select",
  slots: selectAnatomy.keys(),
  base: {
    root: {
      display: "flex",
      flexDirection: "column",
      gap: "1",
      width: "full",
    },
    content: {
      background: "bg.default",
      borderRadius: "md",
      boxShadow: "lg",
      display: "flex",
      flexDirection: "column",
      zIndex: "dropdown",
      _hidden: {
        display: "none",
      },
      _open: {
        animation: "fadeIn 0.25s ease-out",
      },
      _closed: {
        animation: "fadeOut 0.2s ease-out",
      },
      _focusVisible: {
        outlineOffset: "2px",
        outline: "2px solid",
        outlineColor: "border.outline",
      },
    },
    item: {
      alignItems: "center",
      borderRadius: "md",
      cursor: "pointer",
      display: "flex",
      justifyContent: "space-between",
      transitionDuration: "fast",
      transitionProperty: "background, color",
      transitionTimingFunction: "default",
      _hover: {
        background: "gray.300",
        color: "fg.default",
      },
      _highlighted: {
        background: "gray.300",
        color: "fg.default",
      },
      _selected: {
        color: "fg.default",
      },
      _disabled: {
        color: "fg.disabled",
        cursor: "not-allowed",
        _hover: {
          background: "transparent",
        },
      },
    },
    itemGroupLabel: {
      fontWeight: "semibold",
      textStyle: "sm",
    },
    itemIndicator: {
      color: "fg.default",
    },
    label: {
      color: "fg.default",
      fontWeight: "medium",
    },
    positioner: {
      zIndex: "dropdown",
    },
    trigger: {
      appearance: "none",
      alignItems: "center",
      cursor: "pointer",

      backgroundColor: "whiteAlpha.600",
      borderColor: "blackAlpha.50",
      borderRadius: "lg",
      boxShadow: "xs",

      color: "fg.default",
      display: "inline-flex",
      justifyContent: "space-between",
      outline: 0,
      position: "relative",
      transitionDuration: "normal",
      transitionProperty: "background, box-shadow, border-color",
      transitionTimingFunction: "default",
      width: "full",

      _placeholderShown: {
        color: "fg.subtle",
      },
      "& :where(svg)": {
        color: "fg.subtle",
      },
    },
  },
  defaultVariants: {
    size: "md",
  },
  variants: {
    size: {
      sm: {
        content: { p: "0.5", gap: "1" },
        item: { textStyle: "sm", px: "2", height: "9" },
        itemIndicator: {
          "& :where(svg)": {
            width: "4",
            height: "4",
          },
        },
        itemGroupLabel: {
          px: "2",
          py: "1.5",
        },
        label: { textStyle: "sm" },
        trigger: {
          px: "2.5",
          h: "9",
          minW: "9",
          fontSize: "sm",
          gap: "2",
          "& :where(svg)": {
            width: "4",
            height: "4",
          },
        },
      },
      md: {
        content: { p: "1", gap: "1" },
        item: { textStyle: "md", px: "2", height: "10" },
        itemIndicator: {
          "& :where(svg)": {
            width: "4",
            height: "4",
          },
        },
        itemGroupLabel: {
          px: "2",
          py: "1.5",
        },
        label: { textStyle: "sm" },
        trigger: {
          px: "3",
          h: "10",
          minW: "10",
          fontSize: "md",
          gap: "2",
          "& :where(svg)": {
            width: "4",
            height: "4",
          },
        },
      },
      lg: {
        content: { p: "1.5", gap: "1" },
        item: { textStyle: "md", px: "2", height: "11" },
        itemIndicator: {
          "& :where(svg)": {
            width: "5",
            height: "5",
          },
        },
        itemGroupLabel: {
          px: "2",
          py: "1.5",
        },
        label: { textStyle: "sm" },
        trigger: {
          px: "3.5",
          h: "11",
          minW: "11",
          fontSize: "md",
          gap: "2",
          "& :where(svg)": {
            width: "5",
            height: "5",
          },
        },
      },
    },
  },
});
