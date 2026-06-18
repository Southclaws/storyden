import { getPatternStyles, patternFns } from './runtime';
import { css } from '../css/index';

const vstackConfig = {transform(props) {
	const { justify, gap, ...rest } = props;
	return {
		display: "flex",
		alignItems: "center",
		justifyContent: justify,
		gap,
		flexDirection: "column",
		...rest
	};
},defaultValues:{gap:'8px'}}

export function vstackRaw(styles) {
  const s = getPatternStyles(vstackConfig, styles || {})
  return vstackConfig.transform(s, patternFns)
}

export const vstack = /* @__PURE__ */ Object.assign(function vstack(styles = {}) {
  return css(vstackRaw(styles))
}, { raw: vstackRaw })