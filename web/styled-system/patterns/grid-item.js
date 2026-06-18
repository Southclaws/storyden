import { getPatternStyles, patternFns } from './runtime';
import { css } from '../css/index';

const gridItemConfig = {transform(props, { map }) {
	const { colSpan, rowSpan, colStart, rowStart, colEnd, rowEnd, ...rest } = props;
	const spanFn = (v) => v === "auto" ? v : `span ${v}`;
	return {
		gridColumn: colSpan != null ? map(colSpan, spanFn) : void 0,
		gridRow: rowSpan != null ? map(rowSpan, spanFn) : void 0,
		gridColumnStart: colStart,
		gridColumnEnd: colEnd,
		gridRowStart: rowStart,
		gridRowEnd: rowEnd,
		...rest
	};
}}

export function gridItemRaw(styles) {
  const s = getPatternStyles(gridItemConfig, styles || {})
  return gridItemConfig.transform(s, patternFns)
}

export const gridItem = /* @__PURE__ */ Object.assign(function gridItem(styles = {}) {
  return css(gridItemRaw(styles))
}, { raw: gridItemRaw })