import { defineTokens } from "@pandacss/dev";

import amber from "./amber";
import blue from "./blue";
import green from "./green";
import neutral from "./neutral";
import orange from "./orange";
import pink from "./pink";
import red from "./red";
import slate from "./slate";
import tomato from "./tomato";

export const colours = defineTokens.colors({
  current: { value: "currentColor" },
  black: {
    DEFAULT: { value: "#000000" },
    a1: { value: "rgba(0, 0, 0, 0.05)" },
    a2: { value: "rgba(0, 0, 0, 0.1)" },
    a3: { value: "rgba(0, 0, 0, 0.15)" },
    a4: { value: "rgba(0, 0, 0, 0.2)" },
    a5: { value: "rgba(0, 0, 0, 0.3)" },
    a6: { value: "rgba(0, 0, 0, 0.4)" },
    a7: { value: "rgba(0, 0, 0, 0.5)" },
    a8: { value: "rgba(0, 0, 0, 0.6)" },
    a9: { value: "rgba(0, 0, 0, 0.7)" },
    a10: { value: "rgba(0, 0, 0, 0.8)" },
    a11: { value: "rgba(0, 0, 0, 0.9)" },
    a12: { value: "rgba(0, 0, 0, 0.95)" },
  },
  white: {
    DEFAULT: { value: "#ffffff" },
    a1: { value: "rgba(255, 255, 255, 0.05)" },
    a2: { value: "rgba(255, 255, 255, 0.1)" },
    a3: { value: "rgba(255, 255, 255, 0.15)" },
    a4: { value: "rgba(255, 255, 255, 0.2)" },
    a5: { value: "rgba(255, 255, 255, 0.3)" },
    a6: { value: "rgba(255, 255, 255, 0.4)" },
    a7: { value: "rgba(255, 255, 255, 0.5)" },
    a8: { value: "rgba(255, 255, 255, 0.6)" },
    a9: { value: "rgba(255, 255, 255, 0.7)" },
    a10: { value: "rgba(255, 255, 255, 0.8)" },
    a11: { value: "rgba(255, 255, 255, 0.9)" },
    a12: { value: "rgba(255, 255, 255, 0.95)" },
  },

  amber: amber.tokens,
  blue: blue.tokens,
  green: green.tokens,
  orange: orange.tokens,
  pink: pink.tokens,
  red: red.tokens,
  slate: slate.tokens,
  tomato: tomato.tokens,

  gray: neutral.tokens,
  neutral: neutral.tokens,
  transparent: { value: "rgb(0 0 0 / 0)" },

  // Accent colours are dynamic CSS vars and loaded from theme.css at runtime.
  accent: {
    light: {
      "1": { value: "var(--accent-colour-flat-fill-50)" },
      "2": { value: "var(--accent-colour-flat-fill-100)" },
      "3": { value: "var(--accent-colour-flat-fill-200)" },
      "4": { value: "var(--accent-colour-flat-fill-300)" },
      "5": { value: "var(--accent-colour-flat-fill-400)" },
      "6": { value: "var(--accent-colour-flat-fill-500)" },
      "7": { value: "var(--accent-colour-flat-fill-600)" },
      "8": { value: "var(--accent-colour-flat-fill-700)" },
      "9": { value: "var(--accent-colour-flat-fill-800)" },
      "10": { value: "var(--accent-colour-flat-fill-900)" },
      text: {
        "1": { value: "var(--accent-colour-flat-text-50)" },
        "2": { value: "var(--accent-colour-flat-text-100)" },
        "3": { value: "var(--accent-colour-flat-text-200)" },
        "4": { value: "var(--accent-colour-flat-text-300)" },
        "5": { value: "var(--accent-colour-flat-text-400)" },
        "6": { value: "var(--accent-colour-flat-text-500)" },
        "7": { value: "var(--accent-colour-flat-text-600)" },
        "8": { value: "var(--accent-colour-flat-text-700)" },
        "9": { value: "var(--accent-colour-flat-text-800)" },
        "10": { value: "var(--accent-colour-flat-text-900)" },
      },
    },
    dark: {
      "1": { value: "var(--accent-colour-dark-fill-50)" },
      "2": { value: "var(--accent-colour-dark-fill-100)" },
      "3": { value: "var(--accent-colour-dark-fill-200)" },
      "4": { value: "var(--accent-colour-dark-fill-300)" },
      "5": { value: "var(--accent-colour-dark-fill-400)" },
      "6": { value: "var(--accent-colour-dark-fill-500)" },
      "7": { value: "var(--accent-colour-dark-fill-600)" },
      "8": { value: "var(--accent-colour-dark-fill-700)" },
      "9": { value: "var(--accent-colour-dark-fill-800)" },
      "10": { value: "var(--accent-colour-dark-fill-900)" },
      text: {
        "1": { value: "var(--accent-colour-dark-text-50)" },
        "2": { value: "var(--accent-colour-dark-text-100)" },
        "3": { value: "var(--accent-colour-dark-text-200)" },
        "4": { value: "var(--accent-colour-dark-text-300)" },
        "5": { value: "var(--accent-colour-dark-text-400)" },
        "6": { value: "var(--accent-colour-dark-text-500)" },
        "7": { value: "var(--accent-colour-dark-text-600)" },
        "8": { value: "var(--accent-colour-dark-text-700)" },
        "9": { value: "var(--accent-colour-dark-text-800)" },
        "10": { value: "var(--accent-colour-dark-text-900)" },
      },
    },
  },
});
