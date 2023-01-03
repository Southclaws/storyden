import { Theme } from "@chakra-ui/react";

export const Heading: Theme["components"]["Heading"] = {
  baseStyle: {
    fontFamily: "heading",
    fontWeight: "bold",
  },
  sizes: {
    "4xl": {
      fontSize: ["6xl", null, "7xl"],
      lineHeight: 1,
    },
    "3xl": {
      fontSize: ["5xl", null, "6xl"],
      lineHeight: 1,
    },
    "2xl": {
      fontSize: ["4xl", null, "5xl"],
      lineHeight: [1.2, null, 1],
    },
    xl: {
      fontSize: ["3xl", null, "4xl"],
      lineHeight: [1.33, null, 1.2],
    },
    lg: {
      fontSize: ["2xl", null, "3xl"],
      lineHeight: [1.33, null, 1.2],
    },
    md: {
      fontSize: "xl",
      lineHeight: 1.2,
    },
    sm: {
      fontSize: "md",
      lineHeight: 1.2,
    },
    xs: {
      fontSize: "sm",
      lineHeight: 1.2,
    },
  },
  variants: {
    h1: { fontSize: "2xl", as: "h1" },
    h2: { fontSize: "xl", as: "h2" },
    h3: { fontSize: "lg", as: "h3" },
    h4: { fontSize: "md", as: "h4" },
    h5: { fontSize: "sm", as: "h5" },
    h6: { fontSize: "sm", as: "h6" },
  },
  defaultProps: {
    size: "xl",
  },
};
