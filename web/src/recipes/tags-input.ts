import { tagsInputAnatomy } from "@ark-ui/anatomy";
import { defineSlotRecipe } from "@pandacss/dev";

export const tagsInput = defineSlotRecipe({
  className: "tagsInput",
  slots: tagsInputAnatomy.keys(),
  jsx: ["TagsInput", "Combotags"],
  staticCss: [{ size: ["sm", "md"] }],
  base: {
    root: {
      colorPalette: "accent",
      display: "flex",
      flexDirection: "column",
      gap: "1.5",
      width: "full",
      minWidth: "0",
    },
    control: {
      alignItems: "center",
      borderColor: "border.default",
      borderRadius: "l2",
      borderWidth: "1px",
      display: "flex",
      flexWrap: "wrap",
      outline: 0,
      minWidth: "0",
      transitionDuration: "normal",
      transitionProperty: "border-color, box-shadow",
      transitionTimingFunction: "default",
      _focusWithin: {
        borderColor: "colorPalette.default",
        boxShadow: "0 0 0 1px var(--colors-color-palette-default)",
      },
    },
    input: {
      background: "transparent",
      color: "fg.default",
      outline: "none",
      _placeholder: {
        opacity: "full",
        color: "fg.subtle",
      },
    },
    item: {
      overflow: "hidden",
    },
    itemPreview: {
      alignItems: "center",
      borderColor: "border.default",
      borderRadius: "l1",
      borderWidth: "1px",
      color: "fg.default",
      display: "inline-flex",
      fontWeight: "medium",
      minWidth: "0",
      width: "full",
      _highlighted: {
        borderColor: "colorPalette.default",
        boxShadow: "0 0 0 1px var(--colors-color-palette-default)",
      },
      _hidden: {
        display: "none",
      },
    },
    itemText: {
      textOverflow: "ellipsis",
      textWrap: "nowrap",
      overflow: "hidden",
    },
    itemInput: {
      background: "transparent",
      color: "fg.default",
      outline: "none",
    },
    label: {
      color: "fg.default",
      fontWeight: "medium",
      textStyle: "sm",
    },
  },
  defaultVariants: {
    size: "md",
  },
  variants: {
    size: {
      sm: {
        root: {
          gap: "1.5",
        },
        control: {
          fontSize: "sm",
          gap: "1.5",
          minW: "10",
          px: "1",
          py: "1",
        },
        itemPreview: {
          gap: "1",
          h: "6",
          pe: "1",
          ps: "2",
          textStyle: "sm",
        },
      },
      md: {
        root: {
          gap: "1.5",
        },
        control: {
          fontSize: "md",
          gap: "1.5",
          minW: "10",
          px: "3",
          py: "7px", // TODO line break
        },
        itemPreview: {
          gap: "1",
          h: "6",
          pe: "1",
          ps: "2",
          textStyle: "sm",
        },
      },
    },
  },
});
