import { getPatternStyles, patternFns } from './runtime';
import { css } from '../css/index';

const circleConfig = {transform(props) {
	const { size, ...rest } = props;
	return {
		display: "flex",
		alignItems: "center",
		justifyContent: "center",
		flex: "0 0 auto",
		width: size,
		height: size,
		borderRadius: "9999px",
		...rest
	};
}}

export function circleRaw(styles) {
  const s = getPatternStyles(circleConfig, styles || {})
  return circleConfig.transform(s, patternFns)
}

export const circle = /* @__PURE__ */ Object.assign(function circle(styles = {}) {
  return css(circleRaw(styles))
}, { raw: circleRaw })