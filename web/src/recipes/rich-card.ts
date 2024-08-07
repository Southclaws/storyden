import { defineSlotRecipe } from "@pandacss/dev";

export const richCard = defineSlotRecipe({
  className: "rich-card",
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
      "--card-image-max-height": "50px",
      "--card-border-radius": "radii.lg",

      containerType: "inline-size",
      display: "grid",
      width: "full",
      gap: "0",
      overflow: "hidden",
      borderRadius: "var(--card-border-radius)",
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
      lineClamp: "1",
      _hover: {
        textDecoration: "underline",
      },
    },
    text: {
      lineClamp: "1",
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
          aspectRatio: "1",
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
          // NOTE: This solves a small rendering issue on Chrome where the image
          // ever so slightly overflows the border radius. The cause looks like
          // it's due to a combination of the overflow hidden and object-fit.
          borderBottomRadius: "calc(var(--card-border-radius) + 3px)",
        },
        mediaMissing: {
          gridRow: "1 / 2",
          gridColumn: "1 / 2",
          height: "full",
          // NOTE: This padding is to center the Empty component. It's not
          // really perfect but close enough! It's a bit of a hack.
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
          gridTemplateRows: "minmax(min-content, 1fr)",
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
          gridTemplateColumns: "1fr 2fr minmax(0, min-content)",
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
