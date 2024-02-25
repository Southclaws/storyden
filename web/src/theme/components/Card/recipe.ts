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
  ],
  base: {
    root: {
      "--text-lines": "2",

      containerType: "inline-size",
      display: "grid",
      width: "full",
      gap: "0",
      overflow: "hidden",
      borderRadius: "lg",
      boxShadow: "sm",
      minHeight: "0",
    },
    mediaBackdrop: {
      width: "full",
      height: "full",
      objectPosition: "center",
      objectFit: "cover",
      blur: "xl",
      filter: "auto",
      zIndex: "1",
      opacity: "0.1",
      contain: "size",
    },
    mediaContainer: {
      zIndex: "2",
    },
    media: {
      width: "full",
      height: "full",
      objectFit: "cover",
    },
    textArea: {
      zIndex: "2",
      padding: "2",
      minWidth: "0",
      overflow: "hidden",
    },
    title: {
      display: "block",
      overflow: "hidden",
      textWrap: "nowrap",
      textOverflow: "ellipsis",
      _hover: {
        textDecoration: "underline",
      },
    },
    text: {
      display: "block",
      lineClamp: "var(--text-lines)",
      textOverflow: "ellipsis",
    },
  },
  variants: {
    mediaDisplay: {
      with: {
        root: {
          "--card-image-display": "block",
          "--card-row-areas": `"text text media"`,
        },
      },
      without: {
        root: {
          "--card-media-display": "none",
          "--card-row-areas": `"text text text"`,
        },
      },
    },
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
        root: {
          gridTemplateRows: "1fr",
          gridTemplateColumns: "2fr 1fr 1fr",
          gridTemplateAreas: "var(--card-row-areas)",
        },
        mediaBackdrop: {
          gridRow: "1 / 1",
          gridColumn: "2 / 4",
        },
        mediaContainer: {
          gridArea: "media",
        },
        textArea: {
          gridArea: "text",
          background: "backgroundGradientH",
        },
        text: {
          _containerSmall: {
            display: "none",
          },
        },
      },
    },
  },
  defaultVariants: {
    mediaDisplay: "with",
    shape: "box",
  },
  jsx: ["Card"],
});
