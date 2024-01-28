import {
  defineConfig,
  defineSemanticTokens,
  defineTokens,
} from "@pandacss/dev";

export default defineConfig({
  preflight: true,
  strictPropertyValues: true,
  strictTokens: false,
  include: ["./src/**/*.{js,jsx,ts,tsx}", "./pages/**/*.{js,jsx,ts,tsx}"],
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
          fg: {
            default: {
              value: {
                base: "{colors.blackAlpha.700}",
                _dark: "{colors.gray.50}",
              },
            },
          },
          bg: {
            canvas: { value: "{colors.gray.100}" },
            default: {
              value: { base: "{colors.white}", _dark: "{colors.gray.200}" },
            },
            opaque: {
              value: {
                base: "{colors.whiteAlpha.700}",
                _dark: "{colors.blackAlpha.700}",
              },
            },
          },
        },
        blurs: {
          frosted: { value: "8px" },
        },
      }),
      tokens: {
        fonts: {
          mona: {
            value:
              "var(--font-mona-sans), Roboto, 'Helvetica Neue', 'Arial Nova', 'Nimbus Sans', Arial, sans-serif",
          },
        },
        colors: defineTokens.colors({
          whiteAlpha: {
            50: { value: "rgba(255, 255, 255, 0.04)" },
            100: { value: "rgba(255, 255, 255, 0.06)" },
            200: { value: "rgba(255, 255, 255, 0.08)" },
            300: { value: "rgba(255, 255, 255, 0.16)" },
            400: { value: "rgba(255, 255, 255, 0.24)" },
            500: { value: "rgba(255, 255, 255, 0.36)" },
            600: { value: "rgba(255, 255, 255, 0.48)" },
            700: { value: "rgba(255, 255, 255, 0.64)" },
            800: { value: "rgba(255, 255, 255, 0.80)" },
            900: { value: "rgba(255, 255, 255, 0.92)" },
          },
          blackAlpha: {
            50: { value: "rgba(0, 0, 0, 0.04)" },
            100: { value: "rgba(0, 0, 0, 0.06)" },
            200: { value: "rgba(0, 0, 0, 0.08)" },
            300: { value: "rgba(0, 0, 0, 0.16)" },
            400: { value: "rgba(0, 0, 0, 0.24)" },
            500: { value: "rgba(0, 0, 0, 0.36)" },
            600: { value: "rgba(0, 0, 0, 0.48)" },
            700: { value: "rgba(0, 0, 0, 0.64)" },
            800: { value: "rgba(0, 0, 0, 0.80)" },
            900: { value: "rgba(0, 0, 0, 0.92)" },
          },
        }),
      },
    },
  },
  outdir: "styled-system",
  jsxFramework: "react",
});
