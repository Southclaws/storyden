import { getPatternStyles, patternFns } from './runtime';
import { css } from '../css/index';

const stackConfig = {transform(props) {
	const { align, justify, direction, gap, ...rest } = props;
	return {
		display: "flex",
		flexDirection: direction,
		alignItems: align,
		justifyContent: justify,
		gap,
		...rest
	};
},defaultValues:{direction:'column',gap:'8px'}}

export function stackRaw(styles) {
  const s = getPatternStyles(stackConfig, styles || {})
  return stackConfig.transform(s, patternFns)
}

export const stack = /* @__PURE__ */ Object.assign(function stack(styles = {}) {
  return css(stackRaw(styles))
}, { raw: stackRaw })