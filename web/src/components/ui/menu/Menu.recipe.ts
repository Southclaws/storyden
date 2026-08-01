import { menuAnatomy } from "@ark-ui/react";
import { defineSlotRecipe } from "@pandacss/dev";

const itemStyle = {
  alignItems: "center",
  borderRadius: "xs",
  cursor: "pointer",
  display: "flex",
  fontWeight: "medium",
  gap: "2",
  textStyle: "sm",
  transitionDuration: "fast",
  transitionProperty: "background, color",
  transitionTimingFunction: "default",
  color: "text.default",
  _hover: {
    background: "control.hoverBackground",
    "& :where(svg)": {
      color: "text.default",
    },
  },
  _highlighted: {
    background: "selection.background",
  },
  "& :where(svg)": {
    color: "text.muted",
  },
  _disabled: {
    color: "text.disabled",
    cursor: "not-allowed",
    _hover: {
      color: "text.disabled",
      background: "none",
    },
  },
};

export const menu = defineSlotRecipe({
  className: "menu",
  jsx: ["Menu"],
  staticCss: [{ size: ["xs", "sm", "md"] }],
  slots: menuAnatomy.keys(),
  base: {
    itemGroupLabel: {
      fontWeight: "semibold",
      textStyle: "sm",
    },
    content: {
      maxHeight: "var(--available-height)",
      maxW: "20rem",
      overflowY: "scroll",
      background: "background.overlay",
      backdropBlur: "frosted",
      backdropFilter: "auto",
      borderColor: "border.strong",
      borderRadius: "sm",
      borderWidth: "thin",
      boxShadow: "overlay",
      color: "text.default",
      display: "flex",
      flexDirection: "column",
      outline: "none",
      width: "calc(100% + 2rem)",
      zIndex: "popover",
      _hidden: {
        display: "none",
      },
      _open: {
        animation: "fadeIn 0.25s ease-out",
      },
      _closed: {
        animation: "fadeOut 0.2s ease-out",
      },
    },
    itemGroup: {
      display: "flex",
      flexDirection: "column",
    },
    positioner: {
      zIndex: "popover",
    },
    item: itemStyle,
    triggerItem: itemStyle,
    trigger: {
      flexShrink: "0",
    },
  },
  defaultVariants: {
    size: "xs",
  },
  variants: {
    size: {
      xs: {
        itemGroup: {
          gap: "1",
        },
        itemGroupLabel: {
          py: "1",
          px: "1",
          mx: "1",
          textStyle: "xs",
        },
        content: {
          py: "1",
          gap: "1",
        },
        item: {
          h: "6",
          px: "1",
          mx: "1",
          textStyle: "xs",
          "& :where(svg)": {
            width: "4",
            height: "4",
          },
        },
        optionItem: {
          h: "8",
          px: "1.5",
          mx: "1",
          "& :where(svg)": {
            width: "4",
            height: "4",
          },
        },
        triggerItem: {
          h: "6",
          px: "1",
          mx: "1",
          textStyle: "xs",
          "& :where(svg)": {
            width: "4",
            height: "4",
          },
        },
      },
      sm: {
        itemGroup: {
          gap: "1",
        },
        itemGroupLabel: {
          py: "2",
          px: "2",
          mx: "1",
        },
        content: {
          py: "1",
          gap: "1",
        },
        item: {
          h: "9",
          px: "2",
          mx: "1",
          "& :where(svg)": {
            width: "4",
            height: "4",
          },
        },
        optionItem: {
          h: "9",
          px: "2",
          mx: "1",
          "& :where(svg)": {
            width: "4",
            height: "4",
          },
        },
        triggerItem: {
          h: "9",
          px: "2",
          mx: "1.5",
          "& :where(svg)": {
            width: "4",
            height: "4",
          },
        },
      },
      md: {
        itemGroup: {
          gap: "1",
        },
        itemGroupLabel: {
          py: "2.5",
          px: "2.5",
          mx: "1",
        },
        content: {
          py: "1",
          gap: "1",
        },
        item: {
          h: "10",
          px: "2.5",
          mx: "1",
          "& :where(svg)": {
            width: "4",
            height: "4",
          },
        },
        optionItem: {
          h: "10",
          px: "2.5",
          mx: "1",
          "& :where(svg)": {
            width: "4",
            height: "4",
          },
        },
        triggerItem: {
          h: "10",
          px: "2.5",
          mx: "1.5",
          "& :where(svg)": {
            width: "4",
            height: "4",
          },
        },
      },
      lg: {
        itemGroup: {
          gap: "1",
        },
        itemGroupLabel: {
          py: "2.5",
          px: "2.5",
          mx: "1",
        },
        content: {
          py: "1",
          gap: "1",
        },
        item: {
          h: "11",
          px: "2.5",
          mx: "1",
          "& :where(svg)": {
            width: "5",
            height: "5",
          },
        },
        optionItem: {
          h: "11",
          px: "2.5",
          mx: "1",
          "& :where(svg)": {
            width: "5",
            height: "5",
          },
        },
        triggerItem: {
          h: "11",
          px: "2.5",
          mx: "1.5",
          "& :where(svg)": {
            width: "5",
            height: "5",
          },
        },
      },
    },
  },
});
