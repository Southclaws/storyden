import {
  defineConfig,
  defineSemanticTokens,
  defineTokens,
} from "@pandacss/dev";

const tokens = defineTokens({
  fonts: {
    joie: { value: "var(--font-joie)" },
    worksans: { value: "var(--font-worksans)" },
    hedvig: { value: "var(--font-hedvig)" },
    intelone: { value: "var(--font-intelone)" },
  },
  colors: {
    Mono: {
      ink: {
        type: "color",
        value: "#131b1a",
      },
      slush: {
        type: "color",
        value: "#f2f0ef",
      },
    },
    Primary: {
      forest: {
        type: "color",
        value: "#307343",
      },
      saddle: {
        type: "color",
        value: "#854627",
      },
      campfire: {
        type: "color",
        value: "#d68e4d",
      },
      moonlit: {
        type: "color",
        value: "#104059",
      },
    },
    Shades: {
      iron: {
        type: "color",
        value: "#303e47",
      },
      slate: {
        type: "color",
        value: "#212429",
      },
      newspaper: {
        type: "color",
        value: "#d8dbcd",
      },
      stone: {
        type: "color",
        value: "#8cada4",
      },
    },
  },
});

export default defineConfig({
  preflight: false,
  strictPropertyValues: true,
  strictTokens: false,
  layers: {
    base: "panda_base",
    tokens: "panda_tokens",
    recipes: "panda_recipes",
    utilities: "panda_utilities",
  },
  include: ["./src/**/*.{js,jsx,ts,tsx}"],
  exclude: [],
  patterns: {
    extend: {
      FrostedGlass: {
        description: `A frosted glass effect for overlays, modals, menus, etc. This is most prominently used on the navigation overlays and menus.`,
        properties: {},
        transform() {
          return {
            backgroundColor: "bg.opaque",
            backdropBlur: "frosted",
            backdropFilter: "auto",
          };
        },
      },
      Floating: {
        description: `Floating overlay elements.`,
        properties: {},
        transform() {
          return {
            backgroundColor: "bg.opaque",
            backdropBlur: "frosted",
            backdropFilter: "auto",
            borderRadius: "lg",
            boxShadow: "sm",
          };
        },
      },
      Card: {
        description: `A card component that can be used to display content in a container with a border and a shadow.`,
        properties: {
          kind: {
            type: "enum",
            value: ["edge", "default"],
          },
          display: {
            type: "property",
            value: "display",
          },
        },
        transform(props) {
          const { kind, display } = props;

          const padding = kind === "edge" ? "0" : "2";

          return {
            display,
            flexDirection: "column",
            gap: "1",
            width: "full",
            overflow: "hidden",
            boxShadow: "sm",
            borderRadius: "lg",
            backgroundColor: "bg.default",
            padding,
          };
        },
      },

      linkButton: {
        description: "Link button",
        transform: (props) => ({
          backgroundColor: "white",
          alignItems: "center",
          appearance: "none",
          borderRadius: "lg",
          boxShadow: "xs",
          cursor: "pointer",
          display: "inline-flex",
          fontWeight: "semibold",
          minWidth: "0",
          justifyContent: "center",
          outline: "none",
          position: "relative",
          transitionDuration: "normal",
          transitionProperty: "background, border-color, color, box-shadow",
          transitionTimingFunction: "default",
          userSelect: "none",
          verticalAlign: "middle",
          whiteSpace: "nowrap",
          _hover: {
            background: "gray.100",
            boxShadow: "md",
          },
          _focusVisible: {
            outlineOffset: "2px",
            outline: "2px solid",
            outlineColor: "border.outline",
          },
          _active: {
            backgroundColor: "gray.200",
          },
          h: "11",
          minW: "11",
          textStyle: "md",
          px: "5",
          gap: "2",
          "& svg": {
            width: "4",
            height: "4",
          },
          ...props,
        }),
      },
    },
  },

  theme: {
    extend: {
      semanticTokens: defineSemanticTokens({
        colors: {
          // fg: {
          //   default: {
          //     value: {
          //       base: "{colors.blackAlpha.700}",
          //       _dark: "{colors.gray.50}",
          //     },
          //   },
          // },
          // bg: {
          //   canvas: { value: "{colors.gray.100}" },
          //   default: {
          //     value: { base: "{colors.white}", _dark: "{colors.gray.200}" },
          //   },
          //   opaque: {
          //     value: {
          //       base: "{colors.whiteAlpha.700}",
          //       _dark: "{colors.blackAlpha.700}",
          //     },
          //   },
          // },
        },
        blurs: {
          frosted: { value: "8px" },
        },
      }),
      tokens: tokens,
    },
  },
  outdir: "src/styled-system",
  jsxFramework: "react",
});
