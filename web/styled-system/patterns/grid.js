import { getPatternStyles, patternFns } from './runtime';
import { css } from '../css/index';

const gridConfig = {transform(props, { map, isCssUnit }) {
	const { columnGap, rowGap, gap, columns, minChildWidth, ...rest } = props;
	const getValue = (v) => isCssUnit(v) ? v : `token(sizes.${v}, ${v})`;
	return {
		display: "grid",
		gridTemplateColumns: columns != null ? map(columns, (v) => `repeat(${v}, minmax(0, 1fr))`) : minChildWidth != null ? map(minChildWidth, (v) => `repeat(auto-fit, minmax(${getValue(v)}, 1fr))`) : void 0,
		gap,
		columnGap,
		rowGap,
		...rest
	};
},defaultValues(props) {
	return { gap: props.columnGap || props.rowGap ? void 0 : "8px" };
}}

export function gridRaw(styles) {
  const s = getPatternStyles(gridConfig, styles || {})
  return gridConfig.transform(s, patternFns)
}

export const grid = /* @__PURE__ */ Object.assign(function grid(styles = {}) {
  return css(gridRaw(styles))
}, { raw: gridRaw })