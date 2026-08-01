import { colorPickerAnatomy } from "@ark-ui/react";
import { defineSlotRecipe } from "@pandacss/dev";

export const colorPicker = defineSlotRecipe({
  className: "colorPicker",
  slots: colorPickerAnatomy.keys(),
  base: {
    root: {
      display: "flex",
      flexDirection: "column",
      gap: "1.5",
    },
    label: {
      color: "fg.default",
      fontWeight: "medium",
      textStyle: "sm",
    },
    control: {
      alignItems: "center",
      display: "flex",
      flexDirection: "row",
      gap: "2",
    },
    trigger: {
      flexShrink: "0",
      overflow: "visible",
      p: "0.5!",
      "& [data-scope=color-picker][data-part=swatch]": {
        borderRadius: "xs",
        boxShadow: "none",
        height: "full",
        width: "full",
      },
    },
    content: {
      background: "bg.default",
      borderRadius: "md",
      boxShadow: "overlay",
      display: "flex",
      flexDirection: "column",
      maxWidth: "sm",
      p: "4",
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
    area: {
      height: "36",
      borderRadius: "sm",
      overflow: "hidden",
    },
    areaThumb: {
      borderRadius: "full",
      height: "2.5",
      width: "2.5",
      boxShadow: "white 0px 0px 0px 2px, black 0px 0px 2px 1px",
      outline: "none",
    },
    areaBackground: {
      height: "full",
    },
    channelSlider: {
      borderRadius: "sm",
    },
    channelSliderTrack: {
      height: "3",
      borderRadius: "sm",
    },
    swatchGroup: {
      display: "grid",
      gridTemplateColumns: "repeat(7, 1fr)",
      gap: "2",
      background: "bg.default",
    },
    swatch: {
      height: "6",
      width: "6",
      borderRadius: "sm",
      boxShadow:
        "0 0 0 1px var(--colors-border-emphasized), 0 0 0 2px var(--colors-bg-default) inset",
    },
    channelSliderThumb: {
      borderRadius: "full",
      height: "2.5",
      width: "2.5",
      boxShadow: "white 0px 0px 0px 2px, black 0px 0px 2px 1px",
      transform: "translate(-50%, -50%)",
      outline: "none",
    },
    transparencyGrid: {
      borderRadius: "sm",
    },
  },
});
