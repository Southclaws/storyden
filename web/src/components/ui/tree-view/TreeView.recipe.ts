import { treeViewAnatomy } from "@ark-ui/react";
import { defineSlotRecipe } from "@pandacss/dev";

export const treeView = defineSlotRecipe({
  className: "treeView",
  slots: treeViewAnatomy.keys(),
  jsx: ["TreeView", "DatagraphNodeTree"],
  base: {
    root: {
      width: "full",
    },
    branch: {
      "&[data-depth='1'] > [data-part='branch-content']": {
        _before: {
          bg: "border.default",
          content: '""',
          height: "full",
          left: "3",
          position: "absolute",
          width: "1px",
          zIndex: "1",
        },
      },
    },
    branchContent: {
      position: "relative",
      textOverflow: "ellipsis",
    },
    branchControl: {
      alignItems: "center",
      borderRadius: "sm",
      color: "text.muted",
      display: "flex",
      fontWeight: "medium",
      gap: "1",
      ps: "calc((var(--depth)) * 22px)",
      py: "1",
      pr: "1",
      h: "8",
      fontSize: "sm",
      lineHeight: "1.25rem",
      transitionDuration: "normal",
      transitionProperty: "background, color",
      transitionTimingFunction: "default",
      "&[data-depth='1']": {
        ps: "calc((var(--depth)) * 22px)",
      },
      "&[data-depth='1'] > [data-part='branch-text'] ": {
        fontWeight: "medium",
      },
      _hover: {
        background: "interactive.emphasized.surface",
        color: "interactive.emphasized.content",
      },
      _selected: {
        backgroundColor: "interactive.selected.surface",
        color: "interactive.selected.content",
      },
    },
    branchIndicator: {
      color: "accent.solid",
      transformOrigin: "center",
      transitionDuration: "normal",
      transitionProperty: "transform",
      transitionTimingFunction: "default",
      "& svg": {
        fontSize: "md",
        width: "4",
        height: "4",
      },
      _open: {
        transform: "rotate(90deg)",
      },
    },
    item: {
      borderRadius: "sm",
      color: "text.subtle",
      cursor: "pointer",
      fontWeight: "medium",
      position: "relative",
      ps: "calc(((var(--depth)) * 22px) + 22px)",
      py: "1",
      fontSize: "sm",
      lineHeight: "1.25rem",
      transitionDuration: "normal",
      transitionProperty: "background, color",
      transitionTimingFunction: "default",
      "&[data-depth='1']": {
        ps: "calc(((var(--depth)) * 22px) + 22px)",
        fontWeight: "medium",
        color: "text.subtle",
        _selected: {
          _before: {
            bg: "transparent",
          },
        },
      },
      _hover: {
        backgroundColor: "interactive.emphasized.surface",
        color: "interactive.emphasized.content",
      },
      _selected: {
        backgroundColor: "interactive.selected.surface",
        color: "interactive.selected.content",
        _hover: {
          backgroundColor: "interactive.emphasized.surface",
          color: "interactive.emphasized.content",
        },
        _before: {
          content: '""',
          position: "absolute",
          left: "3",
          top: "0",
          width: "2px",
          height: "full",
          bg: "text.default",
          zIndex: "1",
        },
      },
    },
    itemText: {
      fontSize: "xs",
    },
    branchText: {
      fontSize: "xs",
    },
    tree: {
      display: "flex",
      flexDirection: "column",
      gap: "0",
    },
  },
  defaultVariants: {
    variant: "clamped",
  },
  variants: {
    variant: {
      clamped: {
        root: {
          overflowX: "hidden",
          overflowY: "clip",
        },
        branchText: {
          width: "full",
          textWrap: "nowrap",
          textOverflow: "ellipsis",
          overflowX: "hidden",
        },
        tree: {
          gap: "1",
        },
      },
      scrollable: {
        root: {
          overflowX: "scroll",
        },
        tree: {
          gap: "1",
        },
      },
    },
  },
});
