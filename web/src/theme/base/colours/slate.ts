import { defineSemanticTokens, defineTokens } from "@pandacss/dev";

const tokens = defineTokens.colors({
  light: {
    "1": { value: "hsl(240, 20%, 99%)" },
    "2": { value: "hsl(240, 20%, 98%)" },
    "3": { value: "hsl(240, 11.1%, 94.7%)" },
    "4": { value: "hsl(240, 9.5%, 91.8%)" },
    "5": { value: "hsl(230, 10.7%, 89%)" },
    "6": { value: "hsl(240, 10.1%, 86.5%)" },
    "7": { value: "hsl(233.3, 9.9%, 82.2%)" },
    "8": { value: "hsl(230.8, 10.2%, 75.1%)" },
    "9": { value: "hsl(230.8, 5.9%, 57.1%)" },
    "10": { value: "hsl(226.2, 5.4%, 52.7%)" },
    "11": { value: "hsl(220, 5.9%, 40%)" },
    "12": { value: "hsl(210, 12.5%, 12.5%)" },
    a1: { value: "hsla(240, 100%, 16.7%, 0)" },
    a2: { value: "hsla(240, 100%, 16.7%, 0)" },
    a3: { value: "hsla(240, 100%, 10%, 0.1)" },
    a4: { value: "hsla(240, 100%, 8.8%, 0.1)" },
    a5: { value: "hsla(229.2, 100%, 9.8%, 0.1)" },
    a6: { value: "hsla(240, 100%, 9.2%, 0.1)" },
    a7: { value: "hsla(232.2, 100%, 9%, 0.2)" },
    a8: { value: "hsla(230, 100%, 9.4%, 0.3)" },
    a9: { value: "hsla(229.7, 100%, 5.7%, 0.5)" },
    a10: { value: "hsla(224.4, 100%, 5.3%, 0.5)" },
    a11: { value: "hsla(219, 100%, 3.9%, 0.6)" },
    a12: { value: "hsla(206.7, 100%, 1.8%, 0.9)" },
  },
  dark: {
    "1": { value: "hsl(240, 5.6%, 7.1%)" },
    "2": { value: "hsl(220, 5.9%, 10%)" },
    "3": { value: "hsl(225, 5.7%, 13.7%)" },
    "4": { value: "hsl(210, 7.1%, 16.5%)" },
    "5": { value: "hsl(214.3, 7.1%, 19.4%)" },
    "6": { value: "hsl(213.3, 7.7%, 22.9%)" },
    "7": { value: "hsl(212.7, 7.6%, 28.4%)" },
    "8": { value: "hsl(212, 7.7%, 38.2%)" },
    "9": { value: "hsl(218.6, 6.3%, 43.9%)" },
    "10": { value: "hsl(221.5, 5.2%, 49.2%)" },
    "11": { value: "hsl(216, 6.8%, 71%)" },
    "12": { value: "hsl(220, 9.1%, 93.5%)" },
    a1: { value: "hsla(220, 5.9%, 10%, 0.02)" },
    a2: { value: "hsla(220, 5.9%, 10%, 0.04)" },
    a3: { value: "hsla(220, 5.9%, 10%, 0.08)" },
    a4: { value: "hsla(220, 5.9%, 10%, 0.12)" },
    a5: { value: "hsla(220, 5.9%, 10%, 0.16)" },
    a6: { value: "hsla(220, 5.9%, 10%, 0.24)" },
    a7: { value: "hsla(220, 5.9%, 10%, 0.32)" },
    a8: { value: "hsla(220, 5.9%, 10%, 0.42)" },
    a9: { value: "hsla(220, 5.9%, 10%, 0.52)" },
    a10: { value: "hsla(220, 5.9%, 10%, 0.62)" },
    a11: { value: "hsla(220, 5.9%, 10%, 0.7)" },
    a12: { value: "hsla(220, 5.9%, 10%, 0.9)" },
  },
});
const semanticTokens = defineSemanticTokens.colors({
  "1": {
    value: {
      _osLight: "{colors.slate.light.1}",
      _osDark: "{colors.slate.dark.1}",
    },
  },
  "2": {
    value: {
      _osLight: "{colors.slate.light.2}",
      _osDark: "{colors.slate.dark.2}",
    },
  },
  "3": {
    value: {
      _osLight: "{colors.slate.light.3}",
      _osDark: "{colors.slate.dark.3}",
    },
  },
  "4": {
    value: {
      _osLight: "{colors.slate.light.4}",
      _osDark: "{colors.slate.dark.4}",
    },
  },
  "5": {
    value: {
      _osLight: "{colors.slate.light.5}",
      _osDark: "{colors.slate.dark.5}",
    },
  },
  "6": {
    value: {
      _osLight: "{colors.slate.light.6}",
      _osDark: "{colors.slate.dark.6}",
    },
  },
  "7": {
    value: {
      _osLight: "{colors.slate.light.7}",
      _osDark: "{colors.slate.dark.7}",
    },
  },
  "8": {
    value: {
      _osLight: "{colors.slate.light.8}",
      _osDark: "{colors.slate.dark.8}",
    },
  },
  "9": {
    value: {
      _osLight: "{colors.slate.light.9}",
      _osDark: "{colors.slate.dark.9}",
    },
  },
  "10": {
    value: {
      _osLight: "{colors.slate.light.10}",
      _osDark: "{colors.slate.dark.10}",
    },
  },
  "11": {
    value: {
      _osLight: "{colors.slate.light.11}",
      _osDark: "{colors.slate.dark.11}",
    },
  },
  "12": {
    value: {
      _osLight: "{colors.slate.light.12}",
      _osDark: "{colors.slate.dark.12}",
    },
  },
  a1: {
    value: {
      _osLight: "{colors.slate.light.a1}",
      _osDark: "{colors.slate.dark.a1}",
    },
  },
  a2: {
    value: {
      _osLight: "{colors.slate.light.a2}",
      _osDark: "{colors.slate.dark.a2}",
    },
  },
  a3: {
    value: {
      _osLight: "{colors.slate.light.a3}",
      _osDark: "{colors.slate.dark.a3}",
    },
  },
  a4: {
    value: {
      _osLight: "{colors.slate.light.a4}",
      _osDark: "{colors.slate.dark.a4}",
    },
  },
  a5: {
    value: {
      _osLight: "{colors.slate.light.a5}",
      _osDark: "{colors.slate.dark.a5}",
    },
  },
  a6: {
    value: {
      _osLight: "{colors.slate.light.a6}",
      _osDark: "{colors.slate.dark.a6}",
    },
  },
  a7: {
    value: {
      _osLight: "{colors.slate.light.a7}",
      _osDark: "{colors.slate.dark.a7}",
    },
  },
  a8: {
    value: {
      _osLight: "{colors.slate.light.a8}",
      _osDark: "{colors.slate.dark.a8}",
    },
  },
  a9: {
    value: {
      _osLight: "{colors.slate.light.a9}",
      _osDark: "{colors.slate.dark.a9}",
    },
  },
  a10: {
    value: {
      _osLight: "{colors.slate.light.a10}",
      _osDark: "{colors.slate.dark.a10}",
    },
  },
  a11: {
    value: {
      _osLight: "{colors.slate.light.a11}",
      _osDark: "{colors.slate.dark.a11}",
    },
  },
  a12: {
    value: {
      _osLight: "{colors.slate.light.a12}",
      _osDark: "{colors.slate.dark.a12}",
    },
  },
  default: {
    value: {
      _osLight: "{colors.slate.light.9}",
      _osDark: "{colors.slate.dark.9}",
    },
  },
  emphasized: {
    value: {
      _osLight: "{colors.slate.light.10}",
      _osDark: "{colors.slate.dark.10}",
    },
  },
  fg: { value: { _osLight: "white", _osDark: "white" } },
  text: {
    value: {
      _osLight: "{colors.slate.light.12}",
      _osDark: "{colors.slate.dark.12}",
    },
  },
});

export default {
  name: "slate",
  tokens,
  semanticTokens,
};
