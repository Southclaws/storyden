import { defineConfig, defineTextStyles } from "@pandacss/dev";

import { admonition } from "@/components/ui/admonition/Admonition.recipe";
import { alert } from "@/components/ui/alert/Alert.recipe";
import { badge } from "@/components/ui/badge/Badge.recipe";
import { button } from "@/components/ui/button/Button.recipe";
import { checkbox } from "@/components/ui/checkbox/Checkbox.recipe";
import { clipboard } from "@/components/ui/clipboard/Clipboard.recipe";
import { colorPicker } from "@/components/ui/color-picker/ColorPicker.recipe";
import { combobox } from "@/components/ui/combobox/Combobox.recipe";
import { datePicker } from "@/components/ui/date-picker/DatePicker.recipe";
import { fileUpload } from "@/components/ui/file-upload/FileUpload.recipe";
import { headingInput } from "@/components/ui/heading-input/HeadingInput.recipe";
import { inputGroup } from "@/components/ui/input-group/InputGroup.recipe";
import { input } from "@/components/ui/input/Input.recipe";
import { menu } from "@/components/ui/menu/Menu.recipe";
import { numberInput } from "@/components/ui/number-input/NumberInput.recipe";
import { pinInput } from "@/components/ui/pin-input/PinInput.recipe";
import { popover } from "@/components/ui/popover/Popover.recipe";
import { progress } from "@/components/ui/progress/Progress.recipe";
import { radioGroup } from "@/components/ui/radio-group/RadioGroup.recipe";
import { richCard } from "@/components/ui/rich-card/RichCard.recipe";
import { select } from "@/components/ui/select/Select.recipe";
import { slider } from "@/components/ui/slider/Slider.recipe";
import { switchRecipe } from "@/components/ui/switch/Switch.recipe";
import { table } from "@/components/ui/table/Table.recipe";
import { tabs } from "@/components/ui/tabs/Tabs.recipe";
import { text } from "@/components/ui/text/Text.recipe";
import { toggleGroup } from "@/components/ui/toggle-group/ToggleGroup.recipe";
import { tooltip } from "@/components/ui/tooltip/Tooltip.recipe";
import { treeView } from "@/components/ui/tree-view/TreeView.recipe";
import { typographyHeading } from "@/components/ui/typography-heading/TypographyHeading.recipe";
import { tokens } from "@/theme/base";
import { semanticTokens } from "@/theme/semantic";

export default defineConfig({
  presets: ["@pandacss/preset-base"],
  preflight: true,
  lightningcss: true,
  strictTokens: true,
  strictPropertyValues: true,
  validation: "error",
  include: ["./src/**/*.tsx"],
  jsxFramework: "react",
  exclude: [],

  conditions: {
    target: "&:target",
    checked:
      "&:is(:checked, [data-checked], [aria-checked=true], [data-state=checked])",
    indeterminate:
      "&:is(:indeterminate, [data-indeterminate], [aria-checked=mixed], [data-state=indeterminate])",
    closed: "&:is([data-state=closed])",
    open: "&:is([open], [data-state=open])",
    on: "&:is([data-state=on])",
    off: "&:is([data-state=off])",
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
        transform(props) {
          return {
            display: "flex",
            flexDirection: "row",
            gap: "3",
            width: "full",
            justifyContent: "space-between",
            ...props,
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
          const { kind, display, ...rest } = props;

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
            ...rest,
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
      breakpoints: {
        sm: "640px",
        md: "768px",
        lg: "1024px",
        xl: "1280px",
        "2xl": "1536px",
      },
      recipes: {
        badge: badge,
        checkbox: checkbox,
        button: button,
        input: input,
        switchRecipe: switchRecipe,
        text: text,
        admonition: admonition,
        headingInput: headingInput,
        typographyHeading: typographyHeading,
        richCard: richCard,
      },
      slotRecipes: {
        alert: alert,
        clipboard: clipboard,
        numberInput: numberInput,
        inputGroup: inputGroup,
        datePicker: datePicker,
        select: select,
        colorPicker: colorPicker,
        combobox: combobox,
        menu: menu,
        fileUpload: fileUpload,
        popover: popover,
        progress: progress,
        table: table,
        slider: slider,
        pinInput: pinInput,
        tabs: tabs,
        radioGroup: radioGroup,
        treeView: treeView,
        toggleGroup: toggleGroup,
        tooltip: tooltip,
      },
      semanticTokens,
      tokens: tokens,
      keyframes: {
        shimmer: {
          "100%": { transform: "translateX(100%)" },
        },
        targetPulse: {
          "0%, 100%": { backgroundColor: "transparent" },
          "50%": { backgroundColor: "var(--colors-bg-emphasized)" },
        },
      },
      textStyles: defineTextStyles({
        xs: { value: { fontSize: "xs", lineHeight: "1.125rem" } },
        sm: { value: { fontSize: "sm", lineHeight: "1.25rem" } },
        md: { value: { fontSize: "md", lineHeight: "1.5rem" } },
        lg: { value: { fontSize: "lg", lineHeight: "1.75rem" } },
        xl: { value: { fontSize: "xl", lineHeight: "1.875rem" } },
        "2xl": { value: { fontSize: "2xl", lineHeight: "2rem" } },
        "3xl": { value: { fontSize: "3xl", lineHeight: "2.375rem" } },
        "4xl": {
          value: {
            fontSize: "4xl",
            lineHeight: "2.75rem",
            letterSpacing: "-0.02em",
          },
        },
        "5xl": {
          value: {
            fontSize: "5xl",
            lineHeight: "3.75rem",
            letterSpacing: "-0.02em",
          },
        },
        "6xl": {
          value: {
            fontSize: "6xl",
            lineHeight: "4.5rem",
            letterSpacing: "-0.02em",
          },
        },
        "7xl": {
          value: {
            fontSize: "7xl",
            lineHeight: "5.75rem",
            letterSpacing: "-0.02em",
          },
        },
      }),
    },
  },

  outdir: "styled-system",
});
