import { defineSlotRecipe } from "@pandacss/dev";

export const card = defineSlotRecipe({
  className: "card",
  slots: [
    // Global bits (span the entire card)
    "root",
    "mediaBackdropContainer",
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
    "mediaMissing",

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
    mediaBackdropContainer: {
      contain: "size",
      zIndex: "1",
    },
    mediaBackdrop: {
      width: "full",
      height: "full",
      objectPosition: "center",
      objectFit: "cover",
      blur: "xl",
      filter: "auto",
      opacity: "0.2",
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
      contain: "size",
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
        mediaMissing: {
          gridRow: "1 / 2",
          gridColumn: "1 / 2",
          height: "full",
          paddingBottom: "3lh",
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
          gridTemplateColumns: "1fr 2fr 1fr",
          gridTemplateAreas: "var(--card-row-areas)",
        },
        mediaBackdropContainer: {
          gridRow: "1 / 1",
          gridColumn: "2 / 4",
        },
        mediaContainer: {
          gridArea: "media",
          maxHeight: "full",
        },
        mediaMissing: {
          display: "none",
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
    size: {
      default: {},
      small: {},
    },
  },
  compoundVariants: [
    {
      size: "small",
      shape: "row",
      css: {
        root: {
          gridTemplateColumns: "1fr 2fr minmax(0, 3lh)",
        },
        text: { display: "none" },
        title: {
          fontSize: "sm",
        },
        controlsOverlayContainer: {
          display: "flex",
          justifyContent: "end",
          alignItems: "start",
          padding: "2",
        },
      },
    },
  ],
  defaultVariants: {
    mediaDisplay: "with",
    shape: "box",
    size: "default",
  },
  jsx: ["Card"],
});
