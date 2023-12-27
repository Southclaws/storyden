import { menuAnatomy } from "@ark-ui/anatomy";
import { defineSlotRecipe } from "@pandacss/dev";

const itemStyle = {
  alignItems: "center",
  borderRadius: "sm",
  cursor: "pointer",
  display: "flex",
  fontWeight: "medium",
  textStyle: "sm",
  transitionDuration: "fast",
  transitionProperty: "background, color",
  transitionTimingFunction: "default",
  _hover: {
    background: "bg.subtle",
    "& :where(svg)": { color: "fg.default" },
  },
  _highlighted: { background: "bg.subtle" },
  "& :where(svg)": { color: "fg.muted" },
};

export const menu = defineSlotRecipe({
  className: "menu",
  slots: menuAnatomy.keys(),
  base: {
    itemGroupLabel: {
      fontWeight: "semibold",
      textStyle: "sm",
      color: "fg.subtle",
    },
    content: {
      background: "bg.default",
      borderRadius: "md",
      boxShadow: "lg",
      display: "flex",
      flexDirection: "column",
      outline: "none",
      width: "calc(100% + 2rem)",
      _hidden: {
        display: "none",
      },
      _open: {
        animation: "fadeIn 0.25s ease-out",
      },
      _closed: {
        display: "none",
        animation: "fadeOut 0.2s ease-out",
      },
    },
    itemGroup: {
      display: "flex",
      flexDirection: "column",
      gap: "1",
    },
    positioner: {
      zIndex: "dropdown",
    },
    item: itemStyle,
    optionItem: itemStyle,
    triggerItem: itemStyle,
  },
  defaultVariants: {
    size: "md",
  },
  variants: {
    size: {
      sm: {
        itemGroupLabel: { py: "1.5", px: "1.5", mx: "1" },
        content: { py: "1", gap: "1" },
        item: {
          h: "8",
          px: "1.5",
          mx: "1",
          "& :where(svg)": { width: "4", height: "4" },
        },
        optionItem: {
          h: "8",
          px: "1.5",
          mx: "1",
          "& :where(svg)": { width: "4", height: "4" },
        },
        triggerItem: {
          h: "8",
          px: "1.5",
          mx: "1",
          "& :where(svg)": { width: "4", height: "4" },
        },
      },
      md: {
        itemGroupLabel: { py: "2.5", px: "2.5", mx: "1" },
        content: { py: "1", gap: "1" },
        item: {
          h: "10",
          px: "2.5",
          mx: "1",
          "& :where(svg)": { width: "4", height: "4" },
        },
        optionItem: {
          h: "10",
          px: "2.5",
          mx: "1",
          "& :where(svg)": { width: "4", height: "4" },
        },
        triggerItem: {
          h: "10",
          px: "2.5",
          mx: "1.5",
          "& :where(svg)": { width: "4", height: "4" },
        },
      },
      lg: {
        itemGroupLabel: { py: "2.5", px: "2.5", mx: "1" },
        content: { py: "1", gap: "1" },
        item: {
          h: "11",
          px: "2.5",
          mx: "1",
          "& :where(svg)": { width: "5", height: "5" },
        },
        optionItem: {
          h: "11",
          px: "2.5",
          mx: "1",
          "& :where(svg)": { width: "5", height: "5" },
        },
        triggerItem: {
          h: "11",
          px: "2.5",
          mx: "1.5",
          "& :where(svg)": { width: "5", height: "5" },
        },
      },
    },
  },
});
