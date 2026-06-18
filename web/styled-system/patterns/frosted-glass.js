import { getPatternStyles, patternFns } from './runtime';
import { css } from '../css/index';

const FrostedGlassConfig = {transform() {
	return {
		backgroundColor: "bg.opaque",
		backdropBlur: "frosted",
		backdropFilter: "auto"
	};
}}

export function FrostedGlassRaw(styles) {
  const s = getPatternStyles(FrostedGlassConfig, styles || {})
  return FrostedGlassConfig.transform(s, patternFns)
}

export const FrostedGlass = /* @__PURE__ */ Object.assign(function FrostedGlass(styles = {}) {
  return css(FrostedGlassRaw(styles))
}, { raw: FrostedGlassRaw })