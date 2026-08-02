import { defineSlotRecipe } from "@pandacss/dev";

import { select } from "../select/Select.recipe";

const selectBase = select.base;
const selectLarge = select.variants?.size?.lg;
const selectGhost = select.variants?.variant?.ghost;

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
      ...selectBase?.trigger,
      ...selectLarge?.trigger,
      ...selectGhost?.trigger,
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
      ...selectBase?.content,
      ...selectLarge?.content,
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
      ...selectBase?.itemGroupLabel,
      ...selectLarge?.itemGroupLabel,
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
      ...selectBase?.item,
      ...selectLarge?.item,
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
