import { getPatternStyles, patternFns } from './runtime';
import { css } from '../css/index';

const FloatingConfig = {transform() {
	return {
		backgroundColor: "bg.opaque",
		backdropBlur: "frosted",
		backdropFilter: "auto",
		boxShadow: "sm"
	};
}}

export function FloatingRaw(styles) {
  const s = getPatternStyles(FloatingConfig, styles || {})
  return FloatingConfig.transform(s, patternFns)
}

export const Floating = /* @__PURE__ */ Object.assign(function Floating(styles = {}) {
  return css(FloatingRaw(styles))
}, { raw: FloatingRaw })