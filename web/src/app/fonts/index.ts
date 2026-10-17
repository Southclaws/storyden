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
      // Use the Latin range printed by scripts/subset-fonts.py.
      prop: "unicode-range",
      value:
        "U+0000, U+0020-007E, U+00A0-00AC, U+00AE-010F, U+0112-0125, U+0128-0131, U+0134-0137, U+0139-013E, U+0143-0148, U+014C-014F, U+0152-0165, U+0168-016F, U+0172-017E, U+01A0-01A1, U+01AF-01B0, U+01CD-01DC, U+01DE-01E3, U+01E6-01ED, U+01F0, U+01F4-01F5, U+01F8-0201, U+0204-0205, U+0208-0209, U+020C-020D, U+0210-0211, U+0214-0215, U+0218-021B, U+021E-021F, U+0226-0233, U+02BB-02BC, U+02C6, U+02DA, U+02DC, U+0300-0304, U+0306-030A, U+030C, U+030F, U+0313, U+0315, U+031B, U+0323, U+0326-0328, U+032C, U+0337-0338, U+0342-0343, U+0346-036F, U+037E, U+0385, U+0387, U+1DC0-1DF5, U+1DFC-1DFF, U+1E02-1E05, U+1E08-1E0D, U+1E10-1E11, U+1E14-1E17, U+1E1C-1E29, U+1E2E-1E33, U+1E36-1E39, U+1E3E-1E47, U+1E4C-1E5D, U+1E60-1E6D, U+1E78-1E93, U+1E97-1E99, U+1EA0-1EF9, U+1FC1, U+1FED-1FEF, U+1FFD, U+2000-200B, U+2010-2027, U+202F-2055, U+2057, U+205F, U+20AC, U+20DB-20DE, U+20E8, U+20F0, U+2122, U+212A-212B, U+2191, U+2193, U+2212, U+2260, U+FEFF",
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
        "U+0000, U+0020-007E, U+00A0-00AC, U+00AE-010F, U+0112-0125, U+0128-0131, U+0134-0137, U+0139-013E, U+0143-0148, U+014C-014F, U+0152-0165, U+0168-016F, U+0172-017E, U+01A0-01A1, U+01AF-01B0, U+01CD-01DC, U+01DE-01E3, U+01E6-01ED, U+01F0, U+01F4-01F5, U+01F8-0201, U+0204-0205, U+0208-0209, U+020C-020D, U+0210-0211, U+0214-0215, U+0218-021B, U+021E-021F, U+0226-0233, U+02BB-02BC, U+02C6, U+02DA, U+02DC, U+0300-0304, U+0306-030A, U+030C, U+030F, U+0313, U+0315, U+031B, U+0323, U+0326-0328, U+032C, U+0337-0338, U+0342-0343, U+0346-036F, U+037E, U+0385, U+0387, U+1DC0-1DF5, U+1DFC-1DFF, U+1E02-1E05, U+1E08-1E0D, U+1E10-1E11, U+1E14-1E17, U+1E1C-1E29, U+1E2E-1E33, U+1E36-1E39, U+1E3E-1E47, U+1E4C-1E5D, U+1E60-1E6D, U+1E78-1E93, U+1E97-1E99, U+1EA0-1EF9, U+1FC1, U+1FED-1FEF, U+1FFD, U+2000-200B, U+2010-2027, U+202F-2055, U+2057, U+205F, U+20AC, U+20DB-20DE, U+20E8, U+20F0, U+2122, U+212A-212B, U+2191, U+2193, U+2212, U+2260, U+FEFF",
    },
  ],
});
