import { comboboxAnatomy } from "@ark-ui/react";
import { defineSlotRecipe } from "@pandacss/dev";

export const combobox = defineSlotRecipe({
  className: "combobox",
  slots: comboboxAnatomy.keys(),
  base: {
    root: {
      display: "flex",
      flexDirection: "column",
      gap: "1.5",
      width: "full",
    },
    control: {
      display: "flex",
      position: "relative",
    },
    label: {
      color: "text.default",
      fontWeight: "medium",
    },
    trigger: {
      alignItems: "center",
      bottom: "0",
      color: "text.subtle",
      display: "flex",
      justifyContent: "center",
      p: "0!",
      position: "absolute",
      right: "0",
      top: "0",
    },
    content: {
      background: "background.overlay",
      borderColor: "border.strong",
      borderRadius: "sm",
      borderWidth: "thin",
      boxShadow: "overlay",
      color: "text.default",
      display: "flex",
      flexDirection: "column",
      maxHeight: "[min(24rem,calc(100vh-2rem))]",
      maxWidth: "[calc(100vw-2rem)]",
      zIndex: "dropdown",
      height: "fit-content",
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
        outlineColor: "border.default",
      },
    },
    item: {
      alignItems: "center",
      borderRadius: "xs",
      cursor: "pointer",
      display: "flex",
      justifyContent: "space-between",
      transitionDuration: "fast",
      transitionProperty: "background, color",
      transitionTimingFunction: "default",
      _hover: {
        background: "control.hoverBackground",
      },
      _highlighted: {
        background: "selection.background",
      },
      _disabled: {
        color: "text.disabled",
        cursor: "not-allowed",
        _hover: {
          background: "transparent",
        },
      },
    },
    list: {
      maxHeight: "64",
      overflowY: "auto",
    },
    itemGroupLabel: {
      fontWeight: "semibold",
      textStyle: "sm",
    },
    itemIndicator: {
      color: "colorPalette.default",
    },
  },
  defaultVariants: {
    size: "sm",
  },
  variants: {
    size: {
      sm: {
        content: { p: "0.5", gap: "1", borderRadius: "sm" },
        item: { textStyle: "xs", px: "2", height: "6", borderRadius: "xs" },
        itemIndicator: {
          "& :where(svg)": {
            width: "3.5",
            height: "3.5",
          },
        },
        itemGroupLabel: {
          px: "2",
          py: "1.5",
        },
        label: { textStyle: "xs" },
        input: { paddingRight: "6!" },
        trigger: {
          height: "full!",
          width: "6!",
          "& :where(svg)": {
            height: "3.5!",
            width: "3.5!",
          },
        },
      },
      md: {
        content: { p: "0.5", gap: "1" },
        item: { textStyle: "sm", px: "2", height: "8" },
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
        input: { paddingRight: "8!" },
        trigger: {
          height: "full!",
          width: "8!",
          "& :where(svg)": {
            height: "4!",
            width: "4!",
          },
        },
      },
      lg: {
        content: { p: "1", gap: "1" },
        item: { textStyle: "md", px: "2", height: "10" },
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
        label: { textStyle: "md" },
        input: { paddingRight: "10!" },
        trigger: {
          height: "full!",
          width: "10!",
          "& :where(svg)": {
            height: "5!",
            width: "5!",
          },
        },
      },
    },
  },
});
