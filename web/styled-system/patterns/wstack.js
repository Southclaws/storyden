import { getPatternStyles, patternFns } from './runtime';
import { css } from '../css/index';

const wstackConfig = {transform(props) {
	return {
		display: "flex",
		flexDirection: "row",
		gap: "3",
		width: "full",
		justifyContent: "space-between",
		...props
	};
}}

export function wstackRaw(styles) {
  const s = getPatternStyles(wstackConfig, styles || {})
  return wstackConfig.transform(s, patternFns)
}

export const wstack = /* @__PURE__ */ Object.assign(function wstack(styles = {}) {
  return css(wstackRaw(styles))
}, { raw: wstackRaw })