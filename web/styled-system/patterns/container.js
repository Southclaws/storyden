import { getPatternStyles, patternFns } from './runtime';
import { css } from '../css/index';

const containerConfig = {transform(props) {
	return {
		position: "relative",
		maxWidth: "8xl",
		mx: "auto",
		px: {
			base: "4",
			md: "6",
			lg: "8"
		},
		...props
	};
}}

export function containerRaw(styles) {
  const s = getPatternStyles(containerConfig, styles || {})
  return containerConfig.transform(s, patternFns)
}

export const container = /* @__PURE__ */ Object.assign(function container(styles = {}) {
  return css(containerRaw(styles))
}, { raw: containerRaw })