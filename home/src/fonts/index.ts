import localFont from "next/font/local";

// prettier-ignore
export const joie = localFont({
  src: [
    { path: "./static/JoieGrotesk-Bold.woff", weight: "700" },
    { path: "./static/JoieGrotesk-Bold.woff2", weight: "700" },
  ],
  preload: true,
  variable: "--font-joie",
});

// prettier-ignore
export const worksans = localFont({
  src: [
    { path: "./static/WorkSans-Black.woff2", weight: "900" },
    { path: "./static/WorkSans-BlackItalic.woff2", weight: "900", style: "italic" },
    { path: "./static/WorkSans-Bold.woff2", weight: "700" },
    { path: "./static/WorkSans-BoldItalic.woff2", weight: "700", style: "italic" },
    { path: "./static/WorkSans-ExtraBold.woff2", weight: "800" },
    { path: "./static/WorkSans-ExtraBoldItalic.woff2", weight: "800", style: "italic" },
    { path: "./static/WorkSans-ExtraLight.woff2", weight: "200" },
    { path: "./static/WorkSans-ExtraLightItalic.woff2", weight: "200", style: "italic" },
    { path: "./static/WorkSans-Italic.woff2", weight: "400", style: "italic" },
    { path: "./static/WorkSans-Light.woff2", weight: "300" },
    { path: "./static/WorkSans-LightItalic.woff2", weight: "300", style: "italic" },
    { path: "./static/WorkSans-Medium.woff2", weight: "500" },
    { path: "./static/WorkSans-MediumItalic.woff2", weight: "500", style: "italic" },
    { path: "./static/WorkSans-Regular.woff2", weight: "400" },
    { path: "./static/WorkSans-SemiBold.woff2", weight: "600" },
    { path: "./static/WorkSans-SemiBoldItalic.woff2", weight: "600", style: "italic" },
    { path: "./static/WorkSans-Thin.woff2", weight: "100" },
    { path: "./static/WorkSans-ThinItalic.woff2", weight: "100", style: "italic" },
  ],
  preload: true,
  variable: "--font-worksans",
});

// prettier-ignore
export const hedvig = localFont({
  src: "./static/hedvig-letters-serif-v2-latin-regular.woff2",
  preload: true,
  variable: "--font-hedvig",
});

// prettier-ignore
export const intelone = localFont({
  src: [
    { path: "./static/IntelOneMono-Bold.woff2", weight: "700" },
    { path: "./static/IntelOneMono-BoldItalic.woff2", weight: "700", style: "italic" },
    { path: "./static/IntelOneMono-Italic.woff2", weight: "400", style: "italic" },
    { path: "./static/IntelOneMono-Light.woff2", weight: "300" },
    { path: "./static/IntelOneMono-LightItalic.woff2", weight: "300", style: "italic" },
    { path: "./static/IntelOneMono-Medium.woff2", weight: "500" },
    { path: "./static/IntelOneMono-MediumItalic.woff2", weight: "500", style: "italic" },
    { path: "./static/IntelOneMono-Regular.woff2", weight: "400" },
  ],
  preload: true,
  variable: "--font-intelone",
});
