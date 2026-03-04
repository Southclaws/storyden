import { defineSlotRecipe } from "@pandacss/dev";

export const richCard = defineSlotRecipe({
  className: "rich-card",
  slots: [
    "container",
    "root",
    "headerContainer",
    "menuContainer",
    "titleContainer",
    "title",
    "contentContainer",
    "mediaContainer",
    "footerContainer",
    "mediaBackdropContainer",
    "mediaBackdrop",
    "textArea",
    "text",
    "media",
    "mediaMissing",
  ],
  base: {
    container: {
      containerType: "inline-size",
      width: "full",
    },
    root: {
      "--card-border-radius": "radii.lg",
      "--card-backdrop-index": "1",
      "--card-media-index": "2",
      "--card-content-index": "4",

      position: "relative",
      display: "grid",
      width: "full",
      height: "full",
      minWidth: "0",
      minHeight: "0",
      gap: "0",
      overflow: "hidden",
      borderRadius: "var(--card-border-radius)",
      boxShadow: "sm",
      backgroundColor: "bg.default",
    },

    headerContainer: {
      minWidth: "0",
      minHeight: "0",
    },

    menuContainer: {
      display: "flex",
      flexDirection: "column",
    },

    titleContainer: {
      display: "flex",
      flexDirection: "row",
      gap: "1",
    },

    title: {
      lineClamp: "1",
      fontWeight: "bold",
    },

    contentContainer: {
      display: "flex",
      flexDirection: "column",
      justifyContent: "space-between",
      gap: "2",
      height: "full",
      minWidth: "0",
      overflow: "hidden",
    },

    mediaContainer: {
      contain: "size",
    },

    footerContainer: {
      display: "flex",
      justifyContent: "start",
      width: "full",
      minWidth: "0",
      gap: "2",
    },

    mediaBackdropContainer: {
      contain: "size",
      pointerEvents: "none",
    },

    //

    mediaBackdrop: {
      width: "full",
      height: "full",
      objectPosition: "center",
      objectFit: "cover",
      blur: "xl",
      filter: "auto",
      opacity: "0.1",
    },

    media: {
      width: "full",
      height: "full",
      objectFit: "cover",
      // This fixes pixelated images in Chrome... for some reason?
      overflowClipMargin: "unset",
    },

    textArea: {
      mb: "2",
    },

    text: {
      lineClamp: "1",
    },
  },
  variants: {
    backgroundColor: {
      default: {
        root: {
          backgroundColor: "bg.default",
        },
      },
      emphasized: {
        root: {
          backgroundColor: "bg.emphasized",
        },
      },
      accent: {
        root: {
          backgroundColor: "bg.accent",
        },
      },
    },
    shape: {
      row: {
        root: {
          gridTemplateColumns: `[edge-start] 0.5rem [content-start] 1fr [content-end] 0.5rem [media-start] minmax(0, 25cqw) [media-end] 0.5rem [edge-end]`,
          gridTemplateRows: `[edge-start] 0.5rem [header-start] min-content [header-end] 0 [title-start] min-content [title-end] 0 [content-start] 1fr [content-end] 0 [footer-start] min-content [footer-end] 0.5rem [edge-end]`,
        },

        container: {
          width: "full",
        },

        // -
        // Container slots.
        // -

        headerContainer: {
          gridColumn: "content-start / content-end",
          gridRow: "header-start / header-end",
        },

        menuContainer: {
          gridColumn: "media-start / media-end",
          gridRow: "header-start / header-end",
          alignSelf: "start",
          justifySelf: "end",
        },

        titleContainer: {
          gridColumn: "content-start / content-end",
          gridRow: "title-start / title-end",
        },

        contentContainer: {
          gridColumn: "content-start / content-start",
          gridRow: "content-start / content-end",
          height: "min",
        },

        mediaContainer: {
          gridColumn: "media-start / edge-end",
          gridRow: "edge-start / edge-end",

          // SEE: comment on aspectRatio in "responsive" variant.
          // aspectRatio: "1.777",
          minHeight: "0",
        },

        footerContainer: {
          gridColumn: "content-start / media-end",
          gridRow: "footer-start / footer-end",
          gap: "2",
        },

        mediaBackdropContainer: {
          gridColumn: "edge-start / edge-end",
          gridRow: "edge-start / edge-end",
          background: "backgroundGradientH",
        },

        // -
        // Non-container slots.
        // -

        mediaMissing: {
          display: "none",
        },

        mediaBackdrop: {
          // This blends the backdrop image horizontally so the backdrop doesn't
          // cover the entire row and potentially make it too attention-seeking.
          maskImage:
            "linear-gradient(90deg, transparent 0%, white 40%, white 100%)",
        },
      },
      responsive: {
        root: {
          gridTemplateColumns: `[edge-start] 0.5rem [content-start] 1fr [content-end] 0.5rem [media-start] minmax(0, 25cqw) [media-end] 0.5rem [edge-end]`,
          gridTemplateRows: `[edge-start] 0.5rem [header-start] min-content [header-end] 0 [title-start] min-content [title-end] 0 [content-start] 1fr [content-end] 0 [footer-start] min-content [footer-end] 0.5rem [edge-end]`,
          _containerSmall: {
            gridTemplateColumns: `[edge-start] 0.5rem [content-start] 1fr [content-end] 0.5rem [edge-end]`,
            gridTemplateRows: `[edge-start] 0.5rem [header-start] min-content [header-end] 0 [title-start] min-content [title-end] 0 [content-start] 1fr [content-end] 0 [media-start] auto [media-end] 0 [footer-start] min-content [footer-end] 0.5rem [edge-end]`,
          },
        },

        // -
        // Container slots.
        // -

        headerContainer: {
          gridColumn: "content-start / content-end",
          gridRow: "header-start / header-end",
        },

        menuContainer: {
          gridColumn: "media-start / media-end",
          gridRow: "header-start / header-end",
          alignSelf: "start",
          justifySelf: "end",
        },

        titleContainer: {
          gridColumn: "content-start / content-end",
          gridRow: "title-start / title-end",
        },

        contentContainer: {
          gridColumn: "content-start / content-start",
          gridRow: "content-start / content-end",
          height: "min",
        },

        mediaContainer: {
          gridColumn: "media-start / edge-end",
          gridRow: "edge-start / edge-end",
          _containerSmall: {
            gridColumn: "content-start / content-end",
            gridRow: "media-start / media-end",
            minHeight: "64",
            marginBottom: "2",
            aspectRatio: "auto",
            height: "unset",
          },
          // NOTE: A firefox bug prevents us from doing this. Because the grid
          // track for the media container cannot be auto in Firefox, it must be
          // a fixed width. If it's a fixed width (25cqw in this case), then the
          // aspect ratio forces a height to be contributed to the parent's size
          // calculations which we do not want as it makes the card taller than
          // the text content. Ultimately, we want the text content to be the
          // only contribution to the height calculations. It's unclear why the
          // contain: size does not fix this in Firefox like it does in Chrome.
          // aspectRatio: "1.777",
          minHeight: "0",
          height: "full",
        },

        footerContainer: {
          gridColumn: "content-start / media-end",
          _containerSmall: {
            gridColumn: "content-start / content-end",
          },
          gridRow: "footer-start / footer-end",
          gap: "2",
        },

        mediaBackdropContainer: {
          gridColumn: "edge-start / edge-end",
          gridRow: "edge-start / edge-end",
          background: "backgroundGradientH",
        },

        // -
        // Non-container slots.
        // -

        media: {
          _containerSmall: {
            borderRadius: "sm",
          },
        },

        mediaMissing: {
          display: "none",
        },
        text: {},
      },
      box: {
        root: {
          gridTemplateColumns: `[edge-start] 0.5rem [content-start] 1fr [content-end] 0.5rem [edge-end]`,
          gridTemplateRows: `[edge-start] 0.5rem min-content min-content auto auto min-content 0.5rem [edge-end]`,
          gridTemplateAreas: `
            ". .       . "
            ". header  . "
            ". title   . "
            ". content . "
            ". media   . "
            ". footer  . "
            ". .       . "
          `,
        },

        // -
        // Container slots.
        // -

        headerContainer: {
          gridColumn: "content-start / content-end",
          gridRow: "header-start / header-end",
        },

        menuContainer: {
          gridColumn: "media-start / media-end",
          gridRow: "header-start / header-end",
          alignSelf: "start",
          justifySelf: "end",
        },

        titleContainer: {
          gridColumn: "content-start / content-end",
          gridRow: "title-start / title-end",
        },

        contentContainer: {
          gridArea: "content",
        },

        mediaContainer: {
          gridArea: "media",
          marginBottom: "2",
          minHeight: "64",
        },

        footerContainer: {
          gridArea: "footer",
          gap: "2",
        },

        mediaBackdropContainer: {
          gridRow: "edge-start / edge-end",
          gridColumn: "edge-start / edge-end",
        },

        // -
        // Non-container slots.
        // -

        media: {
          objectPosition: "center",
          borderRadius: "sm",
        },

        mediaMissing: {
          gridRow: "1 / 2",
          gridColumn: "1 / 2",
          height: "full",
          // NOTE: This padding is to center the Empty component. It's not
          // really perfect but close enough! It's a bit of a hack.
          paddingBottom: "3lh",
        },

        text: {
          _containerSmall: {
            display: "none",
          },
        },
      },
      fill: {
        root: {
          gridTemplateColumns: `[edge-start] 0.5rem [content-start] 1fr [content-end] 0.5rem [edge-end]`,
          gridTemplateRows: `[edge-start] 0.5rem [header-start] min-content [header-end] 300px [content-start] min-content [footer-start] min-content [edge-end]`,
        },

        // -
        // Container slots.
        // -

        headerContainer: {
          gridColumn: "content-start / content-end",
          gridRow: "header-start / header-end",
          borderRadius: "sm",
          padding: "1",

          background: "bg.canvas/60",
          backdropBlur: "frosted",
          backdropGrayscale: "0.5",
          backdropFilter: "auto",
        },

        menuContainer: {
          // Unused in this variant
          display: "none",
        },

        titleContainer: {
          // Unused in this variant
          display: "none",
        },

        contentContainer: {
          gridRow: "content-start / footer-start",
          gridColumn: "edge-start / edge-end",
          padding: "2",
          paddingBottom: "0",
        },

        mediaContainer: {
          gridRow: "edge-start / edge-end",
          gridColumn: "edge-start / edge-end",
          minHeight: "64",
        },

        footerContainer: {
          gridRow: "footer-start / edge-end",
          gridColumn: "edge-start / edge-end",
          paddingTop: "0",
          padding: "2",
          gap: "2",
        },

        mediaBackdropContainer: {
          gridRow: "content-start / edge-end",
          gridColumn: "edge-start / edge-end",
          background: "bg.canvas/90",
          backdropBlur: "frosted",
          backdropGrayscale: "0.5",
          backdropFilter: "auto",
        },

        mediaBackdrop: {
          display: "none",
        },

        // -
        // Non-container slots.
        // -

        media: {
          objectPosition: "top",
          // NOTE: This solves a small rendering issue on Chrome where the image
          // ever so slightly overflows the border radius. The cause looks like
          // it's due to a combination of the overflow hidden and object-fit.
          borderBottomRadius: "calc(var(--card-border-radius) + 3px)",
        },
        text: {
          _containerSmall: {
            display: "none",
          },
        },
      },
    },
  },
  compoundVariants: [],
  defaultVariants: {
    shape: "row",
    backgroundColor: "default",
  },
  jsx: ["Card", "NodeCard"],
});
