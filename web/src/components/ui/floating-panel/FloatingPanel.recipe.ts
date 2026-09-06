import { floatingPanelAnatomy } from "@ark-ui/react";
import { defineSlotRecipe } from "@pandacss/dev";

import { overlayContentStyles } from "@/theme/semantic/overlay";

export const floatingPanel = defineSlotRecipe({
  className: "floating-panel",
  slots: floatingPanelAnatomy.keys(),
  base: {
    positioner: {
      height: "var(--height)",
      width: "var(--width)",
      zIndex: "var(--z-index-popover) !important",
    },
    content: {
      ...overlayContentStyles,
      borderRadius: "panel",
      display: "flex",
      flexDirection: "column",
      height: "full",
      minHeight: "0",
      overflow: "hidden",
      width: "full",
      _open: {
        animation: "fadeIn 0.18s ease-out",
      },
      _closed: {
        animation: "fadeOut 0.12s ease-out",
      },
    },
    dragTrigger: {
      cursor: "move",
      touchAction: "none",
      userSelect: "none",
      width: "full",
    },
    header: {
      alignItems: "center",
      borderBottomColor: "border.default",
      borderBottomWidth: "thin",
      display: "flex",
      gap: "3",
      justifyContent: "space-between",
      minHeight: "10",
      paddingInline: "3",
    },
    title: {
      fontSize: "sm",
      fontWeight: "semibold",
      lineHeight: "tight",
    },
    control: {
      alignItems: "center",
      display: "flex",
      gap: "1",
    },
    body: {
      display: "flex",
      flex: "1",
      flexDirection: "column",
      minHeight: "0",
      overflow: "hidden",
    },
    closeTrigger: {
      cursor: "pointer",
    },
    resizeTrigger: {
      position: "absolute",
      touchAction: "none",
      "&[data-axis=n], &[data-axis=s]": {
        cursor: "ns-resize",
        height: "2",
        insetInline: "2",
      },
      "&[data-axis=e], &[data-axis=w]": {
        cursor: "ew-resize",
        insetBlock: "2",
        width: "2",
      },
      "&[data-axis=n]": { top: "-1" },
      "&[data-axis=s]": { bottom: "-1" },
      "&[data-axis=e]": { right: "-1" },
      "&[data-axis=w]": { left: "-1" },
      "&[data-axis=se]": {
        bottom: "-1",
        cursor: "nwse-resize",
        height: "4",
        right: "-1",
        width: "4",
      },
    },
  },
});
