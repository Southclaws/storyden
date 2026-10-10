import { getPatternStyles, patternFns } from './runtime';
import { memo } from '../helpers';
import { css } from '../css/index';

const frostedGlassConfig = {transform() {
	return {
		backgroundColor: "bg.opaque",
		backdropBlur: "frosted",
		backdropFilter: "auto"
	};
}}

export function frostedGlassRaw(styles) {
  const s = getPatternStyles(frostedGlassConfig, styles || {})
  return frostedGlassConfig.transform(s, patternFns)
}

export const frostedGlass = /* @__PURE__ */ Object.assign(/* @__PURE__ */ memo(function frostedGlass(styles = {}) {
  return css(frostedGlassRaw(styles))
}), { raw: frostedGlassRaw, propKeys: [] })