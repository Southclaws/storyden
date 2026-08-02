import { datePickerAnatomy } from "@ark-ui/react";
import { defineSlotRecipe } from "@pandacss/dev";

export const datePicker = defineSlotRecipe({
  className: "datePicker",
  jsx: ["DatePicker"],
  slots: [...datePickerAnatomy.keys()],
  base: {
    root: {
      display: "flex",
      flexDirection: "column",
      gap: "1.5",
    },
    content: {
      background: "background.overlay",
      borderColor: "border.strong",
      borderRadius: "md",
      borderWidth: "thin",
      boxShadow: "overlay",
      color: "text.default",
      display: "flex",
      flexDirection: "column",
      gap: "3",
      p: "4",
      width: "344px",
      zIndex: "dropdown",
      _open: {
        animation: "fadeIn 0.25s ease-out",
      },
      _closed: {
        animation: "fadeOut 0.2s ease-out",
      },
      _hidden: {
        display: "none",
      },
    },
    control: {
      display: "flex",
      flexDirection: "row",
      gap: "2",
    },
    label: {
      color: "text.default",
      fontWeight: "medium",
      fontSize: "sm",
      lineHeight: "1.25rem",
    },
    tableHeader: {
      color: "text.muted",
      fontWeight: "semibold",
      height: "10",
      fontSize: "sm",
      lineHeight: "1.25rem",
    },
    viewControl: {
      display: "flex",
      gap: "2",
      justifyContent: "space-between",
    },
    table: {
      width: "full",
      borderCollapse: "separate",
      borderSpacing: "1",
      m: "-1",
    },
    tableCell: {
      textAlign: "center",
    },
    tableCellTrigger: {
      width: "100%",
      _today: {
        _before: {
          content: "'−'",
          color: "colorPalette.default",
          position: "absolute",
          marginTop: "6",
        },
      },
      "&[data-in-range]": {
        background: "selection.background",
      },
      _selected: {
        _before: {
          color: "colorPalette.fg",
        },
      },
    },
    view: {
      display: "flex",
      flexDirection: "column",
      gap: "3",
      _hidden: {
        display: "none",
      },
    },
  },
});
