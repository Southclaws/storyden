import { getPatternStyles, patternFns } from './runtime';
import { css } from '../css/index';

const centerConfig = {transform(props) {
	const { inline, ...rest } = props;
	return {
		display: inline ? "inline-flex" : "flex",
		alignItems: "center",
		justifyContent: "center",
		...rest
	};
}}

export function centerRaw(styles) {
  const s = getPatternStyles(centerConfig, styles || {})
  return centerConfig.transform(s, patternFns)
}

export const center = /* @__PURE__ */ Object.assign(function center(styles = {}) {
  return css(centerRaw(styles))
}, { raw: centerRaw })