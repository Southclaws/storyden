import {
  defineConfig,
  defineSemanticTokens,
  defineTokens,
} from "@pandacss/dev";
import { range } from "lodash";
import { map } from "lodash/fp";

import { admonition } from "src/theme/components/Admonition/admonition.recipe";
import { button } from "src/theme/components/Button/button.recipe";
import { checkbox } from "src/theme/components/Checkbox/checkbox.recipe";
import { heading } from "src/theme/components/Heading/heading.recipe";
import { input } from "src/theme/components/Input/input.recipe";
import { link } from "src/theme/components/Link/link.recipe";
import { menu } from "src/theme/components/Menu/menu.recipe";
import { popover } from "src/theme/components/Popover/popover.recipe";
import { skeleton } from "src/theme/components/Skeleton/skeleton.recipe";
import { tabs } from "src/theme/components/Tabs/tabs.recipe";
import { titleInput } from "src/theme/components/TitleInput/titleInput.recipe";

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

export default defineConfig({
  preflight: true,
  strictTokens: true,
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
  },

  patterns: {
    extend: {
      FrostedGlass: {
        description: `A frosted glass effect for overlays, modals, menus, etc. This is most prominently used on the `,
        properties: {},
        transform(props) {
          return {
            backgroundColor: "whiteAlpha.800",
            backdropBlur: "frosted",
            backdropFilter: "auto",
            boxShadow: "sm",
            borderRadius: "lg",
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
      titleInput: titleInput,
      heading: heading,
      button: button,
      link: link,
      menu: menu,
      tabs: tabs,
      checkbox: checkbox,
      popover: popover,
      skeleton: skeleton,
    },
    extend: {
      semanticTokens: defineSemanticTokens({
        blurs: {
          frosted: { value: "8px" },
        },
        colors: {
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
          bg: {
            canvas: { value: "{colors.gray.100}" },
            default: {
              value: { base: "{colors.white}", _dark: "{colors.gray.200}" },
            },
            subtle: {
              value: { base: "{colors.gray.200}", _dark: "{colors.gray.300}" },
            },
            muted: {
              value: { base: "{colors.gray.300}", _dark: "{colors.gray.400}" },
            },
            emphasized: {
              value: { base: "{colors.gray.400}", _dark: "{colors.gray.500}" },
            },
            disabled: {
              value: { base: "{colors.gray.300}", _dark: "{colors.gray.400}" },
            },
            destructive: {
              value: { base: "{colors.red.300}", _dark: "{colors.red.400}" },
            },
          },
          fg: {
            default: { value: "{colors.gray.900}" },
            muted: { value: "{colors.gray.600}" },
            subtle: { value: "{colors.gray.500}" },
            disabled: { value: "{colors.gray.400}" },
            destructive: {
              value: { base: "{colors.red.500}", _dark: "{colors.red.400}" },
            },
          },
          border: {
            default: { value: "{colors.blackAlpha.200}" },
            muted: { value: "{colors.gray.500}" },
            subtle: { value: "{colors.gray.300}" },
            disabled: { value: "{colors.gray.400}" },

            outline: { value: "{colors.blackAlpha.50}" },
            accent: { value: "{colors.accent.default}" },
          },
          conicGradient: {
            value: conicGradient,
          },
        },
        spacing: {
          safeBottom: { value: "env(safe-area-inset-bottom)" },
        },
      }),
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

        fontSizes: {
          sm: { value: "1rem" },
          md: { value: "1.125rem" },
          lg: { value: "1.2rem" },
          xl: { value: "1.44rem" },
          "2xl": { value: "1.728rem" },
          "3xl": { value: "2.074rem" },
          "4xl": { value: "2.488rem" },
          heading: {
            1: { value: "var(--font-size-h1)" },
            2: { value: "var(--font-size-h2)" },
            3: { value: "var(--font-size-h3)" },
            4: { value: "var(--font-size-h4)" },
            5: { value: "var(--font-size-h5)" },
            6: { value: "var(--font-size-h6)" },
            variable: {
              1: { value: "var(--font-size-h1-variable)" },
              2: { value: "var(--font-size-h2-variable)" },
              3: { value: "var(--font-size-h3-variable)" },
              4: { value: "var(--font-size-h4-variable)" },
              5: { value: "var(--font-size-h5-variable)" },
              6: { value: "var(--font-size-h6-variable)" },
            },
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
