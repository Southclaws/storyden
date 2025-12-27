import { getPatternStyles, patternFns } from '../helpers.mjs';
import { css } from '../css/index.mjs';

const menuItemColorPaletteConfig = {
transform(props21) {
  return {
    colorPalette: props21["colorPalette"],
    background: "colorPalette.4",
    color: "colorPalette.9",
    _hover: {
      background: "colorPalette.5",
      "& :where(svg)": {
        color: "colorPalette.10"
      }
    },
    _highlighted: {
      background: "colorPalette.5"
    },
    "& :where(svg)": {
      color: "colorPalette.9"
    }
  };
}}

export const getMenuItemColorPaletteStyle = (styles = {}) => {
  const _styles = getPatternStyles(menuItemColorPaletteConfig, styles)
  return menuItemColorPaletteConfig.transform(_styles, patternFns)
}

export const menuItemColorPalette = (styles) => css(getMenuItemColorPaletteStyle(styles))
menuItemColorPalette.raw = getMenuItemColorPaletteStyle