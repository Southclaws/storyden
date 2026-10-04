import { getPatternStyles, patternFns } from './runtime';
import { memo } from '../helpers';
import { css } from '../css/index';

const lstackConfig = {transform(props) {
	return {
		display: "flex",
		gap: "3",
		flexDirection: "column",
		width: "full",
		alignItems: "start",
		...props
	};
}}

export function lstackRaw(styles) {
  const s = getPatternStyles(lstackConfig, styles || {})
  return lstackConfig.transform(s, patternFns)
}

export const lstack = /* @__PURE__ */ Object.assign(/* @__PURE__ */ memo(function lstack(styles = {}) {
  return css(lstackRaw(styles))
}), { raw: lstackRaw, propKeys: [] })