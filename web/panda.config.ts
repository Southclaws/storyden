import { defineConfig } from "@pandacss/dev";

import { admonition } from "@/components/ui/admonition/Admonition.recipe";
import { alert } from "@/components/ui/alert/Alert.recipe";
import { badge, badgeColorPalettes } from "@/components/ui/badge/Badge.recipe";
import { blockEditor } from "@/components/ui/block-editor/BlockEditor.recipe";
import { button } from "@/components/ui/button/Button.recipe";
import { cardBox } from "@/components/ui/card-box/CardBox.recipe";
import { checkbox } from "@/components/ui/checkbox/Checkbox.recipe";
import { clipboard } from "@/components/ui/clipboard/Clipboard.recipe";
import { colorPicker } from "@/components/ui/color-picker/ColorPicker.recipe";
import { combobox } from "@/components/ui/combobox/Combobox.recipe";
import { datePicker } from "@/components/ui/date-picker/DatePicker.recipe";
import { fileUpload } from "@/components/ui/file-upload/FileUpload.recipe";
import { group } from "@/components/ui/group/Group.recipe";
import { headingInput } from "@/components/ui/heading-input/HeadingInput.recipe";
import { inputGroup } from "@/components/ui/input-group/InputGroup.recipe";
import { input } from "@/components/ui/input/Input.recipe";
import { menu } from "@/components/ui/menu/Menu.recipe";
import { multiSelectPicker } from "@/components/ui/multi-select-picker/MultiSelectPicker.recipe";
import { numberInput } from "@/components/ui/number-input/NumberInput.recipe";
import { pageHeader } from "@/components/ui/page-header/PageHeader.recipe";
import { pinInput } from "@/components/ui/pin-input/PinInput.recipe";
import { popover } from "@/components/ui/popover/Popover.recipe";
import { progress } from "@/components/ui/progress/Progress.recipe";
import { radioGroup } from "@/components/ui/radio-group/RadioGroup.recipe";
import { reactList } from "@/components/ui/react-list/ReactList.recipe";
import { sectionNavigation } from "@/components/ui/section-navigation/SectionNavigation.recipe";
import { select } from "@/components/ui/select/Select.recipe";
import { slider } from "@/components/ui/slider/Slider.recipe";
import { cardGrid } from "@/components/ui/surface/CardGrid.recipe";
import { cardRows } from "@/components/ui/surface/CardRows.recipe";
import { richCard } from "@/components/ui/surface/RichCard.recipe";
import { switchRecipe } from "@/components/ui/switch/Switch.recipe";
import { table } from "@/components/ui/table/Table.recipe";
import { tabs } from "@/components/ui/tabs/Tabs.recipe";
import { text } from "@/components/ui/text/Text.recipe";
import { textarea } from "@/components/ui/textarea/Textarea.recipe";
import { toggleGroup } from "@/components/ui/toggle-group/ToggleGroup.recipe";
import { tooltip } from "@/components/ui/tooltip/Tooltip.recipe";
import { dragTree } from "@/components/ui/tree-view/DragTree.recipe";
import { treeView } from "@/components/ui/tree-view/TreeView.recipe";
import { tokens } from "@/theme/base";
import { semanticTokens, textStyles } from "@/theme/semantic";

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
  staticCss: {
    css: [
      {
        properties: {
          colorPalette: [...badgeColorPalettes],
        },
      },
    ],
    recipes: {
      button: [
        {
          intent: ["success", "warning", "destructive"],
          size: ["sm", "md", "lg"],
          variant: ["solid", "outline", "ghost", "subtle", "plain"],
        },
      ],
      checkbox: [{ size: ["sm", "md", "lg"] }],
      combobox: [{ size: ["sm", "md", "lg"] }],
      input: [
        {
          size: ["sm", "md", "lg"],
          variant: ["outline", "ghost", "inset"],
        },
      ],
      textarea: [
        {
          size: ["sm", "md", "lg"],
          variant: ["outline", "ghost", "inset"],
        },
      ],
      inputGroup: [{ size: ["sm", "md", "lg"] }],
      multiSelectPicker: [{ size: ["sm", "md", "lg"] }],
      numberInput: [{ size: ["sm", "md", "lg"] }],
      pinInput: [{ size: ["sm", "md", "lg"] }],
      progress: [
        {
          shape: ["circle", "horizontal"],
          size: ["sm", "md", "lg"],
        },
      ],
      radioGroup: [{ size: ["sm", "md", "lg"] }],
      select: [{ size: ["sm", "md", "lg"], variant: ["outline", "ghost"] }],
      slider: [{ size: ["sm", "md", "lg"] }],
      switchRecipe: [{ size: ["sm", "md", "lg"] }],
      text: [{ variant: ["body", "supporting", "metadata"] }],
      toggleGroup: [
        { size: ["sm", "md", "lg"], variant: ["outline", "ghost"] },
      ],
    },
  },

  conditions: {
    target: "&:target",
    checked:
      "&:is(:checked, [data-checked], [aria-checked=true], [data-state=checked])",
    indeterminate:
      "&:is(:indeterminate, [data-indeterminate], [aria-checked=mixed], [data-state=indeterminate])",
    closed: "&:is([data-state=closed])",
    open: "&:is([open], [data-state=open])",
    detailsOpen: "details[open] &",
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
        transform(props) {
          return {
            display: "flex",
            gap: "3",
            flexDirection: "column",
            width: "full",
            alignItems: "start",
            ...props,
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
        cardBox: cardBox,
        cardRows: cardRows,
        group: group,
        input: input,
        multiSelectPicker: multiSelectPicker,
        text: text,
        textarea: textarea,
        admonition: admonition,
        headingInput: headingInput,
        richCard: richCard,
      },
      slotRecipes: {
        alert: alert,
        blockEditor: blockEditor,
        cardGrid: cardGrid,
        clipboard: clipboard,
        numberInput: numberInput,
        inputGroup: inputGroup,
        datePicker: datePicker,
        dragTree: dragTree,
        select: select,
        sectionNavigation: sectionNavigation,
        colorPicker: colorPicker,
        combobox: combobox,
        menu: menu,
        fileUpload: fileUpload,
        popover: popover,
        progress: progress,
        pageHeader: pageHeader,
        table: table,
        slider: slider,
        pinInput: pinInput,
        tabs: tabs,
        radioGroup: radioGroup,
        reactList: reactList,
        switchRecipe: switchRecipe,
        treeView: treeView,
        toggleGroup: toggleGroup,
        tooltip: tooltip,
      },
      semanticTokens,
      tokens: tokens,
      keyframes: {
        fadeIn: {
          from: { opacity: "0" },
          to: { opacity: "1" },
        },
        fadeOut: {
          from: { opacity: "1" },
          to: { opacity: "0" },
        },
        shimmer: {
          "100%": { transform: "translateX(100%)" },
        },
        targetPulse: {
          "0%, 100%": { backgroundColor: "transparent" },
          "50%": {
            backgroundColor: "var(--colors-interactive-emphasized-surface)",
          },
        },
      },
      textStyles,
    },
  },

  outdir: "styled-system",
});
