import { getPatternStyles, patternFns } from './runtime';
import { css } from '../css/index';

const squareConfig = {transform(props) {
	const { size, ...rest } = props;
	return {
		display: "flex",
		alignItems: "center",
		justifyContent: "center",
		flex: "0 0 auto",
		width: size,
		height: size,
		...rest
	};
}}

export function squareRaw(styles) {
  const s = getPatternStyles(squareConfig, styles || {})
  return squareConfig.transform(s, patternFns)
}

export const square = /* @__PURE__ */ Object.assign(function square(styles = {}) {
  return css(squareRaw(styles))
}, { raw: squareRaw })