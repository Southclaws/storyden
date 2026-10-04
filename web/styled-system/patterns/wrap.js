import { getPatternStyles, patternFns } from './runtime';
import { memo } from '../helpers';
import { css } from '../css/index';

const wrapConfig = {transform(props) {
	const { columnGap, rowGap, gap = columnGap || rowGap ? void 0 : "8px", align, justify, ...rest } = props;
	return {
		display: "flex",
		flexWrap: "wrap",
		alignItems: align,
		justifyContent: justify,
		gap,
		columnGap,
		rowGap,
		...rest
	};
}}

export function wrapRaw(styles) {
  const s = getPatternStyles(wrapConfig, styles || {})
  return wrapConfig.transform(s, patternFns)
}

export const wrap = /* @__PURE__ */ Object.assign(/* @__PURE__ */ memo(function wrap(styles = {}) {
  return css(wrapRaw(styles))
}), { raw: wrapRaw, propKeys: ["align","columnGap","gap","justify","rowGap"] })