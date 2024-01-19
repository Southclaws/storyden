import { defineSlotRecipe } from "@pandacss/dev";

export const card = defineSlotRecipe({
  className: "card",
  slots: [
    "root",
    "textArea",
    "title",
    "text",
    "media",
    "mediaContainer",
    "mediaBackdrop",
    "mediaBackdropContainer",
  ],
  base: {
    root: {
      containerType: "inline-size",
      display: "grid",
      w: "full",
      h: "full",
      gap: "0",
      overflow: "hidden",
      borderRadius: "lg",
      boxShadow: "sm",
    },
    mediaBackdrop: {
      objectPosition: "center",
      objectFit: "cover",
      blur: "xl",
      filter: "auto",
    },
    mediaBackdropContainer: {
      width: "full",
      height: "full",
      zIndex: "0",
    },
    mediaContainer: {
      zIndex: "1",
    },
    media: {
      width: "full",
      height: "full",
      objectFit: "cover",
    },
    textArea: {
      zIndex: "2",
      padding: "2",
      backgroundColor: "bg.opaque",
      backdropBlur: "frosted",
      backdropFilter: "auto",
    },
    text: {
      lineClamp: "2",
    },
    title: {
      lineClamp: "1",
    },
  },
  variants: {
    shape: {
      box: {
        root: {
          gridTemplateRows: "1fr auto",
          gridTemplateColumns: "1fr",
          aspectRatio: "square",
        },
        mediaBackdropContainer: {
          gridRow: "1 / 3",
          gridColumn: "1 / 2",
        },
        mediaContainer: {
          gridRow: "1 / 3",
          gridColumn: "1 / 2",
        },
        media: {
          objectPosition: "top",
        },
        textArea: {
          gridRow: "2 / 3",
          gridColumn: "1 / 2",
        },
        text: {
          _containerSmall: {
            display: "none",
          },
        },
      },
      row: {
        // root: {
        //   gridTemplateRows: "5lh",
        //   gridTemplateColumns: "1fr 5lh",
        // },
        // mediaBackdrop: {
        //   gridRow: "1",
        //   gridColumn: "1 / 3",
        // },
        // media: {
        //   gridRow: "1",
        //   gridColumn: "2 / 3",
        // },
        // textArea: {
        //   gridRow: "1",
        //   gridColumn: "2 / 3",
        // },
      },
    },
  },
  defaultVariants: {
    shape: "box",
  },
});
