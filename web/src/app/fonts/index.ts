import localFont from "next/font/local";

export const inter = localFont({
  src: "static/InterVariable.woff2",
  weight: "100 900",
  style: "normal",
  display: "swap",
  preload: true,
  variable: "--font-inter",
  declarations: [{ prop: "font-family", value: "inter" }],
});

// Separate declarations let italics share the body family without preloading.
export const interItalic = localFont({
  src: "static/InterVariable-Italic.woff2",
  weight: "100 900",
  style: "italic",
  display: "swap",
  preload: false,
  adjustFontFallback: false,
  variable: "--font-inter-italic",
  declarations: [{ prop: "font-family", value: "inter" }],
});
