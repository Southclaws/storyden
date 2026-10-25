import localFont from "next/font/local";

// oxfmt-ignore
export const joie = localFont({
  src: [
    { path: "./static/subsets/JoieGrotesk-Bold.woff2", weight: "700" },
  ],
  preload: true,
  variable: "--font-joie",
});

// oxfmt-ignore
export const worksans = localFont({
  src: [
    { path: "./static/subsets/WorkSans-Black.woff2", weight: "900" },
    { path: "./static/subsets/WorkSans-BlackItalic.woff2", weight: "900", style: "italic" },
    { path: "./static/subsets/WorkSans-Bold.woff2", weight: "700" },
    { path: "./static/subsets/WorkSans-BoldItalic.woff2", weight: "700", style: "italic" },
    { path: "./static/subsets/WorkSans-ExtraBold.woff2", weight: "800" },
    { path: "./static/subsets/WorkSans-ExtraBoldItalic.woff2", weight: "800", style: "italic" },
    { path: "./static/subsets/WorkSans-ExtraLight.woff2", weight: "200" },
    { path: "./static/subsets/WorkSans-ExtraLightItalic.woff2", weight: "200", style: "italic" },
    { path: "./static/subsets/WorkSans-Italic.woff2", weight: "400", style: "italic" },
    { path: "./static/subsets/WorkSans-Light.woff2", weight: "300" },
    { path: "./static/subsets/WorkSans-LightItalic.woff2", weight: "300", style: "italic" },
    { path: "./static/subsets/WorkSans-Medium.woff2", weight: "500" },
    { path: "./static/subsets/WorkSans-MediumItalic.woff2", weight: "500", style: "italic" },
    { path: "./static/subsets/WorkSans-Regular.woff2", weight: "400" },
    { path: "./static/subsets/WorkSans-SemiBold.woff2", weight: "600" },
    { path: "./static/subsets/WorkSans-SemiBoldItalic.woff2", weight: "600", style: "italic" },
    { path: "./static/subsets/WorkSans-Thin.woff2", weight: "100" },
    { path: "./static/subsets/WorkSans-ThinItalic.woff2", weight: "100", style: "italic" },
  ],
  preload: false,
  variable: "--font-worksans",
});

// oxfmt-ignore
export const hedvig = localFont({
  src: "./static/subsets/hedvig-letters-serif-v2-latin-regular.woff2",
  preload: false,
  variable: "--font-hedvig",
});

// oxfmt-ignore
export const intelone = localFont({
  src: [
    { path: "./static/subsets/IntelOneMono-Bold.woff2", weight: "700" },
    { path: "./static/subsets/IntelOneMono-BoldItalic.woff2", weight: "700", style: "italic" },
    { path: "./static/subsets/IntelOneMono-Italic.woff2", weight: "400", style: "italic" },
    { path: "./static/subsets/IntelOneMono-Light.woff2", weight: "300" },
    { path: "./static/subsets/IntelOneMono-LightItalic.woff2", weight: "300", style: "italic" },
    { path: "./static/subsets/IntelOneMono-Medium.woff2", weight: "500" },
    { path: "./static/subsets/IntelOneMono-MediumItalic.woff2", weight: "500", style: "italic" },
    { path: "./static/subsets/IntelOneMono-Regular.woff2", weight: "400" },
  ],
  preload: false,
  variable: "--font-intelone",
});

// oxfmt-ignore
export const gorton = localFont({
  src: [
    { path: "./static/subsets/Gorton-Normal-120.woff2", weight: "400" },
  ],
  preload: false,
  variable: "--font-gorton",
});

// The product typeface is confined to faithful UI demonstrations.
export const product = localFont({
  src: [
    { path: "./static/subsets/Inter-Regular.woff2", weight: "400" },
    { path: "./static/subsets/Inter-SemiBold.woff2", weight: "600" },
  ],
  preload: false,
  variable: "--font-product",
});
