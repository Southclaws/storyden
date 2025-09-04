import { defineSemanticTokens, defineTokens } from "@pandacss/dev";

const tokens = defineTokens.colors({
  light: {
    "1": { value: "#fcfcfc" },
    "2": { value: "#f9f9f9" },
    "3": { value: "#f0f0f0" },
    "4": { value: "#e8e8e8" },
    "5": { value: "#e0e0e0" },
    "6": { value: "#d9d9d9" },
    "7": { value: "#cecece" },
    "8": { value: "#bbbbbb" },
    "9": { value: "#8d8d8d" },
    "10": { value: "#838383" },
    "11": { value: "#646464" },
    "12": { value: "#202020" },
    a1: { value: "#00000003" },
    a2: { value: "#00000006" },
    a3: { value: "#0000000f" },
    a4: { value: "#00000017" },
    a5: { value: "#0000001f" },
    a6: { value: "#00000026" },
    a7: { value: "#00000031" },
    a8: { value: "#00000044" },
    a9: { value: "#00000072" },
    a10: { value: "#0000007c" },
    a11: { value: "#0000009b" },
    a12: { value: "#000000df" },
  },
  dark: {
    "1": { value: "#111111" },
    "2": { value: "#191919" },
    "3": { value: "#222222" },
    "4": { value: "#2a2a2a" },
    "5": { value: "#313131" },
    "6": { value: "#3a3a3a" },
    "7": { value: "#484848" },
    "8": { value: "#606060" },
    "9": { value: "#6e6e6e" },
    "10": { value: "#7b7b7b" },
    "11": { value: "#b4b4b4" },
    "12": { value: "#eeeeee" },
    a1: { value: "#00000000" },
    a2: { value: "#ffffff09" },
    a3: { value: "#ffffff12" },
    a4: { value: "#ffffff1b" },
    a5: { value: "#ffffff22" },
    a6: { value: "#ffffff2c" },
    a7: { value: "#ffffff3b" },
    a8: { value: "#ffffff55" },
    a9: { value: "#ffffff64" },
    a10: { value: "#ffffff72" },
    a11: { value: "#ffffffaf" },
    a12: { value: "#ffffffed" },
  },
});

const semanticTokens = defineSemanticTokens.colors({
  "1": {
    value: {
      base: "{colors.neutral.light.1}",
      osDark: "{colors.neutral.dark.1}",
    },
  },
  "2": {
    value: {
      base: "{colors.neutral.light.2}",
      osDark: "{colors.neutral.dark.2}",
    },
  },
  "3": {
    value: {
      base: "{colors.neutral.light.3}",
      osDark: "{colors.neutral.dark.3}",
    },
  },
  "4": {
    value: {
      base: "{colors.neutral.light.4}",
      osDark: "{colors.neutral.dark.4}",
    },
  },
  "5": {
    value: {
      base: "{colors.neutral.light.5}",
      osDark: "{colors.neutral.dark.5}",
    },
  },
  "6": {
    value: {
      base: "{colors.neutral.light.6}",
      osDark: "{colors.neutral.dark.6}",
    },
  },
  "7": {
    value: {
      base: "{colors.neutral.light.7}",
      osDark: "{colors.neutral.dark.7}",
    },
  },
  "8": {
    value: {
      base: "{colors.neutral.light.8}",
      osDark: "{colors.neutral.dark.8}",
    },
  },
  "9": {
    value: {
      base: "{colors.neutral.light.9}",
      osDark: "{colors.neutral.dark.9}",
    },
  },
  "10": {
    value: {
      base: "{colors.neutral.light.10}",
      osDark: "{colors.neutral.dark.10}",
    },
  },
  "11": {
    value: {
      base: "{colors.neutral.light.11}",
      osDark: "{colors.neutral.dark.11}",
    },
  },
  "12": {
    value: {
      base: "{colors.neutral.light.12}",
      osDark: "{colors.neutral.dark.12}",
    },
  },
  a1: {
    value: {
      base: "{colors.neutral.light.a1}",
      osDark: "{colors.neutral.dark.a1}",
    },
  },
  a2: {
    value: {
      base: "{colors.neutral.light.a2}",
      osDark: "{colors.neutral.dark.a2}",
    },
  },
  a3: {
    value: {
      base: "{colors.neutral.light.a3}",
      osDark: "{colors.neutral.dark.a3}",
    },
  },
  a4: {
    value: {
      base: "{colors.neutral.light.a4}",
      osDark: "{colors.neutral.dark.a4}",
    },
  },
  a5: {
    value: {
      base: "{colors.neutral.light.a5}",
      osDark: "{colors.neutral.dark.a5}",
    },
  },
  a6: {
    value: {
      base: "{colors.neutral.light.a6}",
      osDark: "{colors.neutral.dark.a6}",
    },
  },
  a7: {
    value: {
      base: "{colors.neutral.light.a7}",
      osDark: "{colors.neutral.dark.a7}",
    },
  },
  a8: {
    value: {
      base: "{colors.neutral.light.a8}",
      osDark: "{colors.neutral.dark.a8}",
    },
  },
  a9: {
    value: {
      base: "{colors.neutral.light.a9}",
      osDark: "{colors.neutral.dark.a9}",
    },
  },
  a10: {
    value: {
      base: "{colors.neutral.light.a10}",
      osDark: "{colors.neutral.dark.a10}",
    },
  },
  a11: {
    value: {
      base: "{colors.neutral.light.a11}",
      osDark: "{colors.neutral.dark.a11}",
    },
  },
  a12: {
    value: {
      base: "{colors.neutral.light.a12}",
      osDark: "{colors.neutral.dark.a12}",
    },
  },

  default: {
    value: {
      base: "{colors.neutral.light.9}",
      osDark: "{colors.neutral.dark.9}"
    }
  },
  emphasized: {
    value: {
      base: "{colors.neutral.light.12}",
      osDark: "{colors.neutral.dark.12}"
    },
  },
  fg: {
    value: {
      base: "{colors.neutral.light.12}",
      osDark: "{colors.neutral.dark.12}"
    }
  },
  text: {
    value: {
      base: "{colors.neutral.light.12}",
      osDark: "{colors.neutral.dark.12}"
    }
  },
});

export default {
  name: "neutral",
  tokens,
  semanticTokens,
};
