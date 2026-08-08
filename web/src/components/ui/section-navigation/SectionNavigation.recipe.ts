import { defineSlotRecipe } from "@pandacss/dev";

import {
  selectBaseStyles,
  selectGhostStyles,
  selectLargeStyles,
} from "../select/Select.recipe";

export const sectionNavigation = defineSlotRecipe({
  className: "section-navigation",
  slots: [
    "root",
    "trigger",
    "triggerLabel",
    "triggerIndicator",
    "positioner",
    "content",
    "navigation",
    "section",
    "sectionLabel",
    "items",
    "item",
    "link",
    "linkIndicator",
  ],
  base: {
    root: {
      width: "full",
    },
    trigger: {
      ...selectBaseStyles.trigger,
      ...selectLargeStyles.trigger,
      ...selectGhostStyles.trigger,
      textAlign: "start",
    },
    triggerLabel: {
      minWidth: "0",
      overflow: "hidden",
      textOverflow: "ellipsis",
      whiteSpace: "nowrap",
    },
    triggerIndicator: {
      display: "grid",
      flexShrink: "0",
      placeItems: "center",
    },
    positioner: {
      width: "var(--reference-width)",
      zIndex: "popover",
    },
    content: {
      ...selectBaseStyles.content,
      ...selectLargeStyles.content,
      width: "full",
    },
    navigation: {
      display: "flex",
      minWidth: "0",
      flexDirection: "column",
      gap: "2",
    },
    section: {
      display: "flex",
      minWidth: "0",
      flexDirection: "column",
      gap: "1",
    },
    sectionLabel: {
      ...selectBaseStyles.itemGroupLabel,
      ...selectLargeStyles.itemGroupLabel,
      color: "text.subtle",
    },
    items: {
      display: "flex",
      minWidth: "0",
      flexDirection: "column",
      gap: "1",
      listStyle: "none",
      margin: "0",
      padding: "0",
    },
    item: {
      minWidth: "0",
    },
    link: {
      ...selectBaseStyles.item,
      ...selectLargeStyles.item,
      color: "text.muted",
      textDecoration: "none",
      _focusVisible: {
        outline: "2px solid",
        outlineColor: "border.default",
        outlineOffset: "-2px",
      },
      "&[aria-current=page]": {
        background: "selection.background",
        color: "text.default",
      },
    },
    linkIndicator: {
      color: "text.muted",
      display: "grid",
      flexShrink: "0",
      placeItems: "center",
    },
  },
});
