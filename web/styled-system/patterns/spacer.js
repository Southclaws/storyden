import { getPatternStyles, patternFns } from './runtime';
import { css } from '../css/index';

const spacerConfig = {transform(props, { map, isCssUnit, isCssVar }) {
	const { size, ...rest } = props;
	return {
		alignSelf: "stretch",
		justifySelf: "stretch",
		flex: map(size, (v) => {
			if (v == null) return "1";
			return `0 0 ${isCssUnit(v) || isCssVar(v) ? v : `token(spacing.${v}, ${v})`}`;
		}),
		...rest
	};
}}

export function spacerRaw(styles) {
  const s = getPatternStyles(spacerConfig, styles || {})
  return spacerConfig.transform(s, patternFns)
}

export const spacer = /* @__PURE__ */ Object.assign(function spacer(styles = {}) {
  return css(spacerRaw(styles))
}, { raw: spacerRaw })