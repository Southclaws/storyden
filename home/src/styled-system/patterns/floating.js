import { getPatternStyles, patternFns } from './runtime';
import { memo } from '../helpers';
import { css } from '../css/index';

const floatingConfig = {transform() {
	return {
		backgroundColor: "bg.opaque",
		backdropBlur: "frosted",
		backdropFilter: "auto",
		borderRadius: "lg",
		boxShadow: "sm"
	};
}}

export function floatingRaw(styles) {
  const s = getPatternStyles(floatingConfig, styles || {})
  return floatingConfig.transform(s, patternFns)
}

export const floating = /* @__PURE__ */ Object.assign(/* @__PURE__ */ memo(function floating(styles = {}) {
  return css(floatingRaw(styles))
}), { raw: floatingRaw, propKeys: [] })