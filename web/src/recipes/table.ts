import { defineSlotRecipe } from "@pandacss/dev";

export const table = defineSlotRecipe({
  className: "table",
  slots: ["root", "body", "cell", "footer", "head", "header", "row", "caption"],
  base: {
    root: {
      captionSide: "bottom",
      width: "full",
    },
    body: {
      "& tr:last-child": {
        borderBottomWidth: "0",
      },
    },
    caption: {
      color: "fg.subtle",
    },
    cell: {
      verticalAlign: "middle",
    },
    footer: {
      fontWeight: "medium",
      borderTopWidth: "1px",
      "& tr:last-child": {
        borderBottomWidth: "0",
      },
    },
    header: {
      color: "fg.muted",
      fontWeight: "medium",
      textAlign: "left",
      verticalAlign: "middle",
    },
    row: {
      borderBottomWidth: "1px",
      transitionDuration: "normal",
      transitionProperty: "background, color",
      transitionTimingFunction: "default",
    },
  },
  defaultVariants: {
    size: "md",
    variant: "plain",
  },
  variants: {
    variant: {
      dense: {
        root: {
          borderWidth: "1px",
          borderRadius: "lg",
          overflow: "hidden",
          overflowX: "scroll",
        },
        head: {
          bg: "bg.subtle",
        },
        row: {
          "& td, & th": {
            borderRightWidth: "1px",
            _last: {
              borderRightWidth: "0",
            },
          },
        },
      },
      plain: {
        row: {
          _hover: {
            bg: "bg.subtle",
          },
          _selected: {
            bg: "bg.muted",
          },
        },
      },
    },
    size: {
      sm: {
        root: {
          textStyle: "sm",
        },
        caption: {
          mt: "1",
        },
        cell: {
          height: "auto",
          p: "2",
        },
        header: {
          height: "auto",
          p: "2",
        },
      },
      md: {
        root: {
          textStyle: "sm",
        },
        caption: {
          mt: "4",
        },
        cell: {
          height: "14",
          px: "4",
        },
        header: {
          height: "11",
          px: "4",
        },
      },
    },
  },
});
