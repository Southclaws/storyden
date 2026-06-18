import { getPatternStyles, patternFns } from './runtime';
import { css } from '../css/index';

const menuItemColorPaletteConfig = {transform(props) {
	return {
		colorPalette: props["colorPalette"],
		background: "colorPalette.4",
		color: "colorPalette.9",
		_hover: {
			background: "colorPalette.5",
			"& :where(svg)": { color: "colorPalette.10" }
		},
		_highlighted: { background: "colorPalette.5" },
		"& :where(svg)": { color: "colorPalette.9" }
	};
}}

export function menuItemColorPaletteRaw(styles) {
  const s = getPatternStyles(menuItemColorPaletteConfig, styles || {})
  return menuItemColorPaletteConfig.transform(s, patternFns)
}

export const menuItemColorPalette = /* @__PURE__ */ Object.assign(function menuItemColorPalette(styles = {}) {
  return css(menuItemColorPaletteRaw(styles))
}, { raw: menuItemColorPaletteRaw })