import {
  defineConfig,
  defineSemanticTokens,
  defineTokens,
} from "@pandacss/dev";
import { createPreset } from "@park-ui/panda-preset";
import { range } from "lodash";
import { map } from "lodash/fp";

import { admonition } from "src/theme/components/Admonition/admonition.recipe";
import { button } from "src/theme/components/Button/button.recipe";
import { richCard } from "src/theme/components/Card/recipe";
import { checkbox } from "src/theme/components/Checkbox/checkbox.recipe";
import { heading } from "src/theme/components/Heading/heading.recipe";
import { headingInput } from "src/theme/components/HeadingInput/recipe";
import { input } from "src/theme/components/Input/input.recipe";
import { link } from "src/theme/components/Link/link.recipe";
import { menu } from "src/theme/components/Menu/menu.recipe";
import { popover } from "src/theme/components/Popover/popover.recipe";
import { select } from "src/theme/components/Select/select.recipe";
import { skeleton } from "src/theme/components/Skeleton/skeleton.recipe";
import { tabs } from "src/theme/components/Tabs/tabs.recipe";

// TODO: Dark mode = 40%
const L = "80%";

const C = "0.15";

const lch = (hue: number) => `oklch(${L} ${C} ${hue})`;

const stops = map(lch)(range(0, 361, 10));

const conicGradient = `
conic-gradient(
    ${stops.join(",\n")}
);
`;

const semanticTokens = defineSemanticTokens({
  blurs: {
    frosted: { value: "8px" },
  },
  opacity: {
    0: { value: "0" },
    1: { value: "0.1" },
    2: { value: "0.2" },
    3: { value: "0.3" },
    4: { value: "0.4" },
    5: { value: "0.5" },
    6: { value: "0.6" },
    7: { value: "0.7" },
    8: { value: "0.8" },
    9: { value: "0.9" },
    full: { value: "1" },
  },
  borderWidths: {
    none: { value: "0" },
    hairline: { value: "0.5px" },
    thin: { value: "1px" },
    medium: { value: "3px" },
    thick: { value: "3px" },
  },
  colors: {
    bg: {
      site: {
        value: { base: "{colors.accent.50}", _osDark: "{colors.gray.12}" },
      },
      accent: {
        value: { base: "{colors.accent.500}", _osDark: "{colors.accent.900}" },
      },
      opaque: {
        value: { base: "{colors.white}", _osDark: "{colors.gray.11}" },
      },
    },
    fg: {
      accent: {
        value: { base: "{colors.accent.100}", _osDark: "{colors.accent.200}" },
      },
    },
    border: {
      default: { value: "{colors.blackAlpha.200}" },
      muted: { value: "{colors.gray.5}" },
      subtle: { value: "{colors.gray.3}" },
      disabled: { value: "{colors.gray.4}" },

      outline: { value: "{colors.blackAlpha.50}" },
      accent: { value: "{colors.bg.accent}" },
    },
    conicGradient: {
      value: conicGradient,
    },
    cardBackgroundGradient: {
      value: "linear-gradient(90deg, var(--colors-bg-default), transparent)",
    },
    backgroundGradientH: {
      value: "linear-gradient(90deg, var(--colors-bg-default), transparent)",
    },
    backgroundGradientV: {
      value: "linear-gradient(0deg, var(--colors-bg-default), transparent)",
    },
  },
  spacing: {
    safeBottom: { value: "env(safe-area-inset-bottom)" },
  },
});

export default defineConfig({
  presets: [
    "@pandacss/preset-base",
    "@park-ui/panda-preset",
    createPreset({
      // NOTE: This is just for Park-ui's preset, the actual accent colour is
      // set by the administrator and is a dynamic runtime value.
      accentColor: "neutral",
      additionalColors: ["*"],
    }),
  ],
  preflight: true,
  strictTokens: true,
  strictPropertyValues: true,
  validation: "error",
  include: ["./src/**/*.tsx"],
  jsxFramework: "react",
  exclude: [],

  conditions: {
    checked:
      "&:is(:checked, [data-checked], [aria-checked=true], [data-state=checked])",
    indeterminate:
      "&:is(:indeterminate, [data-indeterminate], [aria-checked=mixed], [data-state=indeterminate])",
    closed: "&:is([data-state=closed])",
    open: "&:is([open], [data-state=open])",
    hidden: "&:is([hidden])",
    current: "&:is([data-current])",
    today: "&:is([data-today])",
    placeholderShown: "&:is(:placeholder-shown, [data-placeholder-shown])",
    collapsed:
      '&:is([aria-collapsed=true], [data-collapsed], [data-state="collapsed"])',
    containerSmall: "@container (max-width: 300px)",
  },

  patterns: {
    extend: {
      LStack: {
        description: "A VStack with full width aligned left.",
        transform() {
          return {
            display: "flex",
            gap: "3",
            flexDirection: "column",
            width: "full",
            alignItems: "start",
          };
        },
      },
      FrostedGlass: {
        description: `A frosted glass effect for overlays, modals, menus, etc. This is most prominently used on the navigation overlays and menus.`,
        properties: {},
        transform() {
          return {
            backgroundColor: "bg.opaque/60",
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
            backgroundColor: "bg.opaque/80",
            backdropBlur: "frosted",
            backdropFilter: "auto",
            borderRadius: "lg",
            boxShadow: "sm",
          };
        },
      },
      CardBox: {
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
    },
  },

  // NOTE: The theme references some CSS variables defined in global.css, this
  // is in order to provide some level of customisability for hosts who want to
  // override CSS with custom rules. Panda is primarily just there to wire it
  // all together and define the semantic tokens.
  //
  // NOTE: There may be some parts of global.css that reference tokens generated
  // by Panda, this is best avoided but it's some leftovers from the early days.
  theme: {
    recipes: {
      admonition: admonition,
      input: input,
      headingInput: headingInput,
      heading: heading,
      button: button,
      link: link,
      // menu: menu,
      // tabs: tabs,
      // select: select,
      // checkbox: checkbox,
      // popover: popover,
      skeleton: skeleton,
      richCard: richCard, // TODO: RENAME
    },
    extend: {
      semanticTokens,
      tokens: defineTokens({
        zIndex: {
          hide: { value: -1 },
          base: { value: 0 },
          docked: { value: 10 },
          dropdown: { value: 1000 },
          sticky: { value: 1100 },
          banner: { value: 1200 },
          overlay: { value: 1300 },
          modal: { value: 1400 },
          popover: { value: 1500 },
          skipLink: { value: 1600 },
          toast: { value: 1700 },
          tooltip: { value: 1800 },
        },
        radii: {
          none: { value: "0" },
          xs: { value: "0.125rem" },
          sm: { value: "0.25rem" },
          md: { value: "0.375rem" },
          lg: { value: "0.5rem" },
          xl: { value: "0.75rem" },
          "2xl": { value: "1rem" },
          "3xl": { value: "1.5rem" },
          full: { value: "9999px" },
        },

        // NOTE: Font sizes are specified in global.css in order to make use of
        // CSS features not available (or, not as readable) in Panda's config.
        fontSizes: {
          sm: { value: "var(--global-font-size-sm)" },
          md: { value: "var(--global-font-size-md)" },
          lg: { value: "var(--global-font-size-lg)" },
          xl: { value: "var(--global-font-size-xl)" },
          "2xl": { value: "var(--global-font-size-2xl)" },
          "3xl": { value: "var(--global-font-size-3xl)" },
          "4xl": { value: "var(--global-font-size-4xl)" },
          heading: {
            1: { value: "var(--global-font-size-h1)" },
            2: { value: "var(--global-font-size-h2)" },
            3: { value: "var(--global-font-size-h3)" },
            4: { value: "var(--global-font-size-h4)" },
            5: { value: "var(--global-font-size-h5)" },
            6: { value: "var(--global-font-size-h6)" },
            variable: {
              1: { value: "var(--global-font-size-h1-variable)" },
              2: { value: "var(--global-font-size-h2-variable)" },
              3: { value: "var(--global-font-size-h3-variable)" },
              4: { value: "var(--global-font-size-h4-variable)" },
              5: { value: "var(--global-font-size-h5-variable)" },
              6: { value: "var(--global-font-size-h6-variable)" },
            },
          },
        },
        colors: defineTokens.colors({
          accent: {
            50: { value: "var(--accent-colour-flat-fill-50)" },
            100: { value: "var(--accent-colour-flat-fill-100)" },
            200: { value: "var(--accent-colour-flat-fill-200)" },
            300: { value: "var(--accent-colour-flat-fill-300)" },
            400: { value: "var(--accent-colour-flat-fill-400)" },
            DEFAULT: { value: "var(--accent-colour-flat-fill-500)" },
            500: { value: "var(--accent-colour-flat-fill-500)" },
            600: { value: "var(--accent-colour-flat-fill-600)" },
            700: { value: "var(--accent-colour-flat-fill-700)" },
            800: { value: "var(--accent-colour-flat-fill-800)" },
            900: { value: "var(--accent-colour-flat-fill-900)" },
            text: {
              50: { value: "var(--accent-colour-flat-text-50)" },
              100: { value: "var(--accent-colour-flat-text-100)" },
              200: { value: "var(--accent-colour-flat-text-200)" },
              300: { value: "var(--accent-colour-flat-text-300)" },
              400: { value: "var(--accent-colour-flat-text-400)" },
              DEFAULT: { value: "var(--accent-colour-flat-text-500)" },
              500: { value: "var(--accent-colour-flat-text-500)" },
              600: { value: "var(--accent-colour-flat-text-600)" },
              700: { value: "var(--accent-colour-flat-text-700)" },
              800: { value: "var(--accent-colour-flat-text-800)" },
              900: { value: "var(--accent-colour-flat-text-900)" },
            },
            dark: {
              50: { value: "var(--accent-colour-dark-fill-50)" },
              100: { value: "var(--accent-colour-dark-fill-100)" },
              200: { value: "var(--accent-colour-dark-fill-200)" },
              300: { value: "var(--accent-colour-dark-fill-300)" },
              400: { value: "var(--accent-colour-dark-fill-400)" },
              DEFAULT: { value: "var(--accent-colour-dark-fill-500)" },
              500: { value: "var(--accent-colour-dark-fill-500)" },
              600: { value: "var(--accent-colour-dark-fill-600)" },
              700: { value: "var(--accent-colour-dark-fill-700)" },
              800: { value: "var(--accent-colour-dark-fill-800)" },
              900: { value: "var(--accent-colour-dark-fill-900)" },
              text: {
                50: { value: "var(--accent-colour-dark-text-50)" },
                100: { value: "var(--accent-colour-dark-text-100)" },
                200: { value: "var(--accent-colour-dark-text-200)" },
                300: { value: "var(--accent-colour-dark-text-300)" },
                400: { value: "var(--accent-colour-dark-text-400)" },
                DEFAULT: { value: "var(--accent-colour-dark-text-500)" },
                500: { value: "var(--accent-colour-dark-text-500)" },
                600: { value: "var(--accent-colour-dark-text-600)" },
                700: { value: "var(--accent-colour-dark-text-700)" },
                800: { value: "var(--accent-colour-dark-text-800)" },
                900: { value: "var(--accent-colour-dark-text-900)" },
              },
            },
          },
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
      }),
    },
    keyframes: {
      shimmer: {
        "100%": { transform: "translateX(100%)" },
      },
    },
  },

  outdir: "styled-system",
});
