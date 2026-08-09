import { treeViewAnatomy } from "@ark-ui/react";
import { defineSlotRecipe } from "@pandacss/dev";

export const dragTree = defineSlotRecipe({
  className: "dragTree",
  slots: [
    ...treeViewAnatomy.keys(),
    "actions",
    "link",
    "row",
    "dropIndicator",
    "dropIndicatorLine",
  ],
  jsx: ["DragTree"],
  base: {
    root: {
      width: "full",
    },
    branchContent: {
      position: "relative",
      textOverflow: "ellipsis",
    },
    row: {
      alignItems: "center",
      borderRadius: "sm",
      color: "text.muted",
      display: "flex",
      fontWeight: "medium",
      gap: "1",
      h: "6",
      ps: "calc((max(var(--depth), 1) - 1) * 22px)",
      py: "0",
      pr: "1",
      fontSize: "sm",
      lineHeight: "1.25rem",
      transitionDuration: "normal",
      transitionProperty: "background, color",
      transitionTimingFunction: "default",
      "&[data-depth='1'] > [data-part='branch-text'] ": {
        fontWeight: "medium",
      },
      _hover: {
        background: "interactive.emphasized.surface",
        color: "interactive.emphasized.content",
      },
      '&[data-selected="true"]': {
        backgroundColor: "interactive.selected.surface",
        color: "interactive.selected.content",
      },
      '&[data-drag-over="true"]': {
        outlineColor: "border.strong",
        outlineOffset: "-2px",
        outlineStyle: "dashed",
        outlineWidth: "thin",
      },
    },
    branchControl: {
      alignItems: "center",
      borderRadius: "sm",
      display: "flex",
      flexShrink: "0",
      position: "relative",
      "&:is(button)": {
        cursor: "pointer",
        _after: {
          content: '\"\"',
          height: "6",
          left: "50%",
          position: "absolute",
          top: "50%",
          transform: "translate(-50%, -50%)",
          width: "6",
        },
        _focusVisible: {
          outline: "2px solid",
          outlineColor: "border.default",
          outlineOffset: "2px",
        },
        _hover: {
          color: "text.default",
        },
      },
    },
    branchIndicator: {
      color: "text.subtle",
      transformOrigin: "center",
      transitionDuration: "normal",
      transitionProperty: "color, transform",
      transitionTimingFunction: "default",
      "& svg": {
        fontSize: "md",
        width: "4",
        height: "4",
      },
      _open: {
        color: "text.muted",
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
    label: {
      srOnly: true,
    },
    branchText: {
      fontSize: "xs",
    },
    actions: {
      alignItems: "center",
      display: "flex",
      flexShrink: "0",
      gap: "2",
      minWidth: "min",
      opacity: "0",
      _groupFocusWithin: {
        opacity: "full",
      },
      _groupHover: {
        opacity: "full",
      },
    },
    link: {
      display: "block",
      flex: "1",
      minWidth: "0",
      overflow: "hidden",
      touchAction: "none",
      textOverflow: "ellipsis",
      userSelect: "none",
      WebkitTouchCallout: "none",
      whiteSpace: "nowrap",
      _hover: {
        textDecoration: "none",
      },
    },
    dropIndicator: {
      height: "1px",
      marginInlineStart: "calc(((var(--depth)) * 22px) + 22px)",
      position: "relative",
    },
    dropIndicatorLine: {
      background: "transparent",
      height: "3px",
      insetBlockStart: "-1px",
      insetInline: "0",
      opacity: "0",
      pointerEvents: "none",
      position: "absolute",
      transitionDuration: "normal",
      transitionProperty: "background, opacity",
      transitionTimingFunction: "default",
      '&[data-active="true"]': {
        background: "selection.background",
        opacity: "full",
      },
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
          gap: "0.5",
        },
      },
      scrollable: {
        root: {
          overflowX: "scroll",
        },
        tree: {
          gap: "0.5",
        },
      },
    },
  },
});
