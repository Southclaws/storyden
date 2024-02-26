import { defineSlotRecipe } from "@pandacss/dev";

export const card = defineSlotRecipe({
  className: "card",
  slots: [
    // Global bits (span the entire card)
    "root",
    "mediaBackdrop",

    // Top level bits
    "contentContainer",
    "mediaContainer",

    // Content container bits
    "textArea",
    "footer",
    "title",
    "text",

    // Media container bits
    "media",

    // Overlay bits
    "controlsOverlayContainer",
    "controls",
  ],
  base: {
    root: {
      "--card-text-lines": "2",
      "--card-image-max-height": "50px",

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
      opacity: "0.2",
      contain: "size",
    },

    contentContainer: {
      display: "flex",
      flexDirection: "column",
      justifyContent: "space-between",
      gap: "2",
      height: "full",
      zIndex: "2",
      padding: "2",
      minWidth: "0",
      overflow: "hidden",
    },
    mediaContainer: {
      zIndex: "2",
    },

    media: {
      width: "full",
      height: "full",
      objectFit: "cover",
    },
    textArea: {},
    footer: {
      display: "flex",
      justify: "start",
      width: "full",
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
      lineClamp: "var(--card-text-lines)",
      textOverflow: "ellipsis",
    },

    // The overlay is used to position the controls such as buttons etc. The
    // container itself should not capture interactions, but the controls do.
    controlsOverlayContainer: {
      zIndex: "3",
      pointerEvents: "none",
    },
    controls: {
      pointerEvents: "auto",
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
        contentContainer: {
          gridRow: "2 / 3",
          gridColumn: "1 / 2",
          backdropBlur: "frosted",
          backdropGrayscale: "0.5",
          backdropFilter: "auto",
          backgroundColor: "bg.opaque/90",
        },
        text: {
          _containerSmall: {
            display: "none",
          },
        },
        controlsOverlayContainer: {
          gridRow: "1 / 1",
          gridColumn: "1 / 4",
          display: "flex",
          justifyContent: "end",
          alignItems: "start",
          padding: "2",
        },
      },
      row: {
        root: {
          gridTemplateRows: "minmax(3lh, 1fr)",
          gridTemplateColumns: "2fr 1fr 1fr",
          gridTemplateAreas: "var(--card-row-areas)",
        },
        mediaBackdrop: {
          gridRow: "1 / 1",
          gridColumn: "2 / 4",
        },
        mediaContainer: {
          gridArea: "media",
          maxHeight: "full",
          contain: "size",
        },
        contentContainer: {
          gridArea: "text",
          background: "backgroundGradientH",
        },
        text: {
          _containerSmall: {
            display: "none",
          },
        },
        controlsOverlayContainer: {
          gridRow: "1 / 1",
          gridColumn: "1 / 4",
          display: "flex",
          justifyContent: "end",
          alignItems: "end",
          padding: "2",
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
