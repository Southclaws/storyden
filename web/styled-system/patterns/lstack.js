import { getPatternStyles, patternFns } from './runtime';
import { css } from '../css/index';

const lstackConfig = {transform() {
	return {
		display: "flex",
		gap: "3",
		flexDirection: "column",
		width: "full",
		alignItems: "start"
	};
}}

export function lstackRaw(styles) {
  const s = getPatternStyles(lstackConfig, styles || {})
  return lstackConfig.transform(s, patternFns)
}

export const lstack = /* @__PURE__ */ Object.assign(function lstack(styles = {}) {
  return css(lstackRaw(styles))
}, { raw: lstackRaw })