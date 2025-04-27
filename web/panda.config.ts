import {
  defineConfig,
  defineSemanticTokens,
  defineTokens,
} from "@pandacss/dev";
import { createPreset } from "@park-ui/panda-preset";
import amber from "@park-ui/panda-preset/colors/amber";
import blue from "@park-ui/panda-preset/colors/blue";
import bronze from "@park-ui/panda-preset/colors/bronze";
import brown from "@park-ui/panda-preset/colors/brown";
import crimson from "@park-ui/panda-preset/colors/crimson";
import cyan from "@park-ui/panda-preset/colors/cyan";
import gold from "@park-ui/panda-preset/colors/gold";
import grass from "@park-ui/panda-preset/colors/grass";
import green from "@park-ui/panda-preset/colors/green";
import indigo from "@park-ui/panda-preset/colors/indigo";
import iris from "@park-ui/panda-preset/colors/iris";
import jade from "@park-ui/panda-preset/colors/jade";
import lime from "@park-ui/panda-preset/colors/lime";
import mauve from "@park-ui/panda-preset/colors/mauve";
import mint from "@park-ui/panda-preset/colors/mint";
import neutral from "@park-ui/panda-preset/colors/neutral";
import olive from "@park-ui/panda-preset/colors/olive";
import orange from "@park-ui/panda-preset/colors/orange";
import pink from "@park-ui/panda-preset/colors/pink";
import plum from "@park-ui/panda-preset/colors/plum";
import purple from "@park-ui/panda-preset/colors/purple";
import red from "@park-ui/panda-preset/colors/red";
import ruby from "@park-ui/panda-preset/colors/ruby";
import sage from "@park-ui/panda-preset/colors/sage";
import sand from "@park-ui/panda-preset/colors/sand";
import sky from "@park-ui/panda-preset/colors/sky";
import slate from "@park-ui/panda-preset/colors/slate";
import teal from "@park-ui/panda-preset/colors/teal";
import tomato from "@park-ui/panda-preset/colors/tomato";
import violet from "@park-ui/panda-preset/colors/violet";
import yellow from "@park-ui/panda-preset/colors/yellow";
import { range } from "lodash";
import { map } from "lodash/fp";

import { admonition } from "@/recipes/admonition";
import { badge } from "@/recipes/badge";
import { button } from "@/recipes/button";
import { colorPicker } from "@/recipes/color-picker";
import { combobox } from "@/recipes/combobox";
import { fileUpload } from "@/recipes/file-upload";
import { headingInput } from "@/recipes/heading-input";
import { input } from "@/recipes/input";
import { menu } from "@/recipes/menu";
import { popover } from "@/recipes/popover";
import { radioGroup } from "@/recipes/radio-group";
import { richCard } from "@/recipes/rich-card";
import { select } from "@/recipes/select";
import { table } from "@/recipes/table";
import { tabs } from "@/recipes/tabs";
import { tagsInput } from "@/recipes/tags-input";
import { toggleGroup } from "@/recipes/toggle-group";
import { tooltip } from "@/recipes/tooltip";
import { treeView } from "@/recipes/tree-view";
import { typographyHeading } from "@/recipes/typography-heading";

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
  fonts: {
    body: { value: "{fonts.inter}" },
    heading: { value: "{fonts.interDisplay}" },
  },
  blurs: {
    frosted: { value: "10px" },
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
  sizes: {
    prose: { value: "65ch" },
  },
  colors: {
    bg: {
      site: {
        value: { base: "{colors.accent.50}" },
      },
      accent: {
        value: { base: "{colors.accent.500}" },
      },
      opaque: {
        value: { base: "{colors.white.a10}" },
      },
      destructive: {
        value: { base: "{colors.tomato.3}" },
      },
      error: {
        value: { base: "{colors.tomato.2}" },
      },
    },
    fg: {
      accent: {
        value: { base: "{colors.accent.100}" },
      },
      destructive: {
        value: { base: "{colors.tomato.8}" },
      },
      error: {
        value: { base: "{colors.tomato.8}" },
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
    scrollGutter: { value: "var(--spacing-2)" },
  },
});

export default defineConfig({
  presets: [
    "@pandacss/preset-base",
    "@park-ui/panda-preset",
    createPreset({
      // NOTE: This is just for Park-ui's preset, the actual accent colour is
      // set by the administrator and is a dynamic runtime value.
      accentColor: neutral,
      grayColor: neutral,
      radius: "lg",
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
    containerSmall: "@container (max-width: 560px)",
    containerMedium: "@container (min-width: 561px) and (max-width: 999px)",
    containerLarge: "@container (min-width: 1000px)",
  },

  patterns: {
    extend: {
      lstack: {
        description: "A VStack with full width aligned left.",
        jsxName: "LStack",
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
      wstack: {
        description: "A HStack with full width and spaced children.",
        jsxName: "WStack",
        transform() {
          return {
            display: "flex",
            flexDirection: "row",
            gap: "3",
            width: "full",
            justifyContent: "space-between",
          };
        },
      },
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
            boxShadow: "sm",
            borderRadius: "lg",
            backgroundColor: "bg.default",
            padding,
          };
        },
      },
      menuItemColorPalette: {
        description: `A color palette for menu items.`,
        properties: {},
        transform(props) {
          return {
            colorPalette: props["colorPalette"],
            background: "colorPalette.4",
            color: "colorPalette.9",
            _hover: {
              background: "colorPalette.5",
              "& :where(svg)": {
                color: "colorPalette.10",
              },
            },
            _highlighted: {
              background: "colorPalette.5",
            },
            "& :where(svg)": {
              color: "colorPalette.9",
            },
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
    extend: {
      recipes: {
        badge: badge,
        input: input,
        admonition: admonition,
        button: button,
        headingInput: headingInput,
        typographyHeading: typographyHeading,
        richCard: richCard,
      },
      slotRecipes: {
        select: select,
        colorPicker: colorPicker,
        combobox: combobox,
        menu: menu,
        fileUpload: fileUpload,
        popover: popover,
        table: table,
        tagsInput: tagsInput,
        tabs: tabs,
        radioGroup: radioGroup,
        treeView: treeView,
        toggleGroup: toggleGroup,
        tooltip: tooltip,
      },
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

        fonts: {
          inter: { value: "var(--font-inter)" },
          interDisplay: { value: "var(--font-inter-display)" },
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
          amber: amber.tokens.light,
          blue: blue.tokens.light,
          bronze: bronze.tokens.light,
          brown: brown.tokens.light,
          crimson: crimson.tokens.light,
          cyan: cyan.tokens.light,
          gold: gold.tokens.light,
          grass: grass.tokens.light,
          green: green.tokens.light,
          indigo: indigo.tokens.light,
          iris: iris.tokens.light,
          jade: jade.tokens.light,
          lime: lime.tokens.light,
          mauve: mauve.tokens.light,
          mint: mint.tokens.light,
          neutral: neutral.tokens.light,
          olive: olive.tokens.light,
          orange: orange.tokens.light,
          pink: pink.tokens.light,
          plum: plum.tokens.light,
          purple: purple.tokens.light,
          red: red.tokens.light,
          ruby: ruby.tokens.light,
          sage: sage.tokens.light,
          sand: sand.tokens.light,
          sky: sky.tokens.light,
          slate: slate.tokens.light,
          teal: teal.tokens.light,
          tomato: tomato.tokens.light,
          violet: violet.tokens.light,
          yellow: yellow.tokens.light,
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
