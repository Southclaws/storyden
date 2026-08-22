import { selectAnatomy } from "@ark-ui/react";
import { defineSlotRecipe } from "@pandacss/dev";

import { overlayContentStyles } from "@/theme/semantic/overlay";

export const selectBaseStyles = {
  content: {
    ...overlayContentStyles,
    borderRadius: "sm",
    display: "flex",
    flexDirection: "column",
    maxHeight: "[min(24rem,calc(100vh-2rem))]",
    maxWidth: "[calc(100vw-2rem)]",
    overflowY: "auto",
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
      outlineColor: "border.default",
    },
  },
  item: {
    alignItems: "center",
    borderRadius: "xs",
    cursor: "pointer",
    display: "flex",
    gap: "2",
    justifyContent: "space-between",
    transitionDuration: "fast",
    transitionProperty: "background, color",
    transitionTimingFunction: "default",
    _hover: {
      background: "background.controlHover",
      color: "text.default",
    },
    _highlighted: {
      background: "selection.background",
      color: "text.default",
    },
    _selected: {
      color: "text.default",
    },
    _disabled: {
      color: "text.disabled",
      cursor: "not-allowed",
      _hover: {
        background: "transparent",
        color: "text.disabled",
      },
    },
  },
  itemGroupLabel: {
    fontWeight: "semibold",
    fontSize: "sm",
    lineHeight: "1.25rem",
  },
  trigger: {
    appearance: "none",
    alignItems: "center",
    borderColor: "border.default",
    borderRadius: "sm",
    cursor: "pointer",
    color: "text.default",
    display: "inline-flex",
    justifyContent: "space-between",
    outline: 0,
    position: "relative",
    transitionDuration: "normal",
    transitionProperty: "background, box-shadow, border-color",
    transitionTimingFunction: "default",
    width: "full",
    _placeholderShown: {
      color: "text.muted",
    },
    _disabled: {
      color: "text.disabled",
      cursor: "not-allowed",
      "& :where(svg)": {
        color: "text.disabled",
      },
    },
    "& :where(svg)": {
      color: "text.muted",
    },
  },
} as const;

export const selectGhostStyles = {
  trigger: {
    _hover: {
      background: "gray.a3",
    },
    _focus: {
      background: "gray.a3",
    },
  },
} as const;

export const selectLargeStyles = {
  content: { p: "1", gap: "1" },
  item: { fontSize: "md", lineHeight: "1.5rem", px: "2", height: "10" },
  itemGroupLabel: {
    px: "2",
    py: "1.5",
  },
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
} as const;

export const select = defineSlotRecipe({
  className: "select",
  slots: selectAnatomy.keys(),
  base: {
    root: {
      colorPalette: "accent",
      display: "flex",
      flexDirection: "column",
      gap: "1.5",
      width: "full",
    },
    control: {
      display: "flex",
      width: "full",
    },
    content: selectBaseStyles.content,
    item: selectBaseStyles.item,
    itemGroupLabel: selectBaseStyles.itemGroupLabel,
    itemIndicator: {
      color: "colorPalette.solid",
    },
    label: {
      color: "text.default",
      fontWeight: "medium",
    },
    trigger: selectBaseStyles.trigger,
  },
  defaultVariants: {
    size: "sm",
    variant: "outline",
  },
  variants: {
    variant: {
      outline: {
        trigger: {
          borderWidth: "1px",
          _focus: {
            borderColor: "colorPalette.solid",
            boxShadow: "0 0 0 1px var(--colors-color-palette-solid)",
          },
        },
      },
      ghost: selectGhostStyles,
    },
    size: {
      sm: {
        content: { p: "0.5", gap: "1", borderRadius: "sm" },
        item: {
          fontSize: "xs",
          lineHeight: "1.125rem",
          px: "2",
          height: "6",
          borderRadius: "sm",
        },
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
        label: { fontSize: "xs", lineHeight: "1.125rem" },
        trigger: {
          px: "2.5",
          h: "6",
          minW: "9",
          fontSize: "xs",
          borderRadius: "sm",
          gap: "2",
          "& :where(svg)": {
            width: "4",
            height: "4",
          },
        },
      },
      md: {
        content: { p: "0.5", gap: "1" },
        item: { fontSize: "sm", lineHeight: "1.25rem", px: "2", height: "9" },
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
        label: { fontSize: "sm", lineHeight: "1.25rem" },
        trigger: {
          px: "2.5",
          h: "8",
          minW: "9",
          fontSize: "sm",
          gap: "2",
          "& :where(svg)": {
            width: "4",
            height: "4",
          },
        },
      },
      lg: {
        content: selectLargeStyles.content,
        item: selectLargeStyles.item,
        itemIndicator: {
          "& :where(svg)": {
            width: "4",
            height: "4",
          },
        },
        itemGroupLabel: selectLargeStyles.itemGroupLabel,
        label: { fontSize: "sm", lineHeight: "1.25rem" },
        trigger: selectLargeStyles.trigger,
      },
    },
  },
});
