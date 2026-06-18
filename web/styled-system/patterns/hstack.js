import { getPatternStyles, patternFns } from './runtime';
import { css } from '../css/index';

const hstackConfig = {transform(props) {
	const { justify, gap, ...rest } = props;
	return {
		display: "flex",
		alignItems: "center",
		justifyContent: justify,
		gap,
		flexDirection: "row",
		...rest
	};
},defaultValues:{gap:'8px'}}

export function hstackRaw(styles) {
  const s = getPatternStyles(hstackConfig, styles || {})
  return hstackConfig.transform(s, patternFns)
}

export const hstack = /* @__PURE__ */ Object.assign(function hstack(styles = {}) {
  return css(hstackRaw(styles))
}, { raw: hstackRaw })