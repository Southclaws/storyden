import { getPatternStyles, patternFns } from './runtime';
import { memo } from '../helpers';
import { css } from '../css/index';

const flexConfig = {transform(props) {
	const { direction, align, justify, wrap: wrap2, basis, grow, shrink, ...rest } = props;
	return {
		display: "flex",
		flexDirection: direction,
		alignItems: align,
		justifyContent: justify,
		flexWrap: wrap2,
		flexBasis: basis,
		flexGrow: grow,
		flexShrink: shrink,
		...rest
	};
}}

export function flexRaw(styles) {
  const s = getPatternStyles(flexConfig, styles || {})
  return flexConfig.transform(s, patternFns)
}

export const flex = /* @__PURE__ */ Object.assign(/* @__PURE__ */ memo(function flex(styles = {}) {
  return css(flexRaw(styles))
}), { raw: flexRaw, propKeys: ["align","basis","direction","grow","justify","shrink","wrap"] })