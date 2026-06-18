import { getPatternStyles, patternFns } from './runtime';
import { css } from '../css/index';

const linkOverlayConfig = {transform(props) {
	return {
		_before: {
			content: "\"\"",
			position: "absolute",
			inset: "0",
			zIndex: "0",
			...props["_before"]
		},
		...props
	};
}}

export function linkOverlayRaw(styles) {
  const s = getPatternStyles(linkOverlayConfig, styles || {})
  return linkOverlayConfig.transform(s, patternFns)
}

export const linkOverlay = /* @__PURE__ */ Object.assign(function linkOverlay(styles = {}) {
  return css(linkOverlayRaw(styles))
}, { raw: linkOverlayRaw })