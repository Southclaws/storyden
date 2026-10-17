import localFont from "next/font/local";

import "./subsets.css";

export const inter = localFont({
  src: "static/subsets/InterVariable-latin.woff2",
  weight: "100 900",
  style: "normal",
  display: "swap",
  preload: true,
  variable: "--font-inter",
  declarations: [
    { prop: "font-family", value: "inter" },
    {
      // Keep this aligned with LATIN in scripts/subset-fonts.py.
      prop: "unicode-range",
      value:
        "U+0000-00FF,U+0131,U+0152-0153,U+02BB-02BC,U+02C6,U+02DA,U+02DC,U+0300-036F,U+2000-206F,U+20AC,U+2122,U+2191,U+2193,U+2212,U+2215,U+FEFF,U+FFFD",
    },
  ],
});

// Separate declarations let italics share the body family without preloading.
export const interItalic = localFont({
  src: "static/subsets/InterVariable-Italic-latin.woff2",
  weight: "100 900",
  style: "italic",
  display: "swap",
  preload: false,
  adjustFontFallback: false,
  variable: "--font-inter-italic",
  declarations: [
    { prop: "font-family", value: "inter" },
    {
      prop: "unicode-range",
      value:
        "U+0000-00FF,U+0131,U+0152-0153,U+02BB-02BC,U+02C6,U+02DA,U+02DC,U+0300-036F,U+2000-206F,U+20AC,U+2122,U+2191,U+2193,U+2212,U+2215,U+FEFF,U+FFFD",
    },
  ],
});
