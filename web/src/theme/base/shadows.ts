import { defineTokens } from "@pandacss/dev";

export const shadows = defineTokens.shadows({
  xs: {
    value: [
      {
        offsetX: 0,
        offsetY: 1,
        blur: 2,
        spread: 0,
        color: "{colors.gray.light.a5}",
      },
      {
        offsetX: 0,
        offsetY: 0,
        blur: 1,
        spread: 0,
        color: "{colors.gray.light.a7}",
      },
    ],
  },
  sm: {
    value: [
      {
        offsetX: 0,
        offsetY: 2,
        blur: 4,
        spread: 0,
        color: "{colors.gray.light.a3}",
      },
      {
        offsetX: 0,
        offsetY: 0,
        blur: 1,
        spread: 0,
        color: "{colors.gray.light.a7}",
      },
    ],
  },
  md: {
    value: [
      {
        offsetX: 0,
        offsetY: 4,
        blur: 8,
        spread: 0,
        color: "{colors.gray.light.a3}",
      },
      {
        offsetX: 0,
        offsetY: 0,
        blur: 1,
        spread: 0,
        color: "{colors.gray.light.a7}",
      },
    ],
  },
  lg: {
    value: [
      {
        offsetX: 0,
        offsetY: 8,
        blur: 16,
        spread: 0,
        color: "{colors.gray.light.a3}",
      },
      {
        offsetX: 0,
        offsetY: 0,
        blur: 1,
        spread: 0,
        color: "{colors.gray.light.a7}",
      },
    ],
  },
  xl: {
    value: [
      {
        offsetX: 0,
        offsetY: 16,
        blur: 24,
        spread: 0,
        color: "{colors.gray.light.a3}",
      },
      {
        offsetX: 0,
        offsetY: 0,
        blur: 1,
        spread: 0,
        color: "{colors.gray.light.a7}",
      },
    ],
  },
  "2xl": {
    value: [
      {
        offsetX: 0,
        offsetY: 24,
        blur: 40,
        spread: 0,
        color: "{colors.gray.light.a3}",
      },
      {
        offsetX: 0,
        offsetY: 0,
        blur: 1,
        spread: 0,
        color: "{colors.gray.light.a7}",
      },
    ],
  },
});