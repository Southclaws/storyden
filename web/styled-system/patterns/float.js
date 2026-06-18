import { getPatternStyles, patternFns } from './runtime';
import { css } from '../css/index';

const floatConfig = {transform(props, { map }) {
	const { offset, offsetX, offsetY, placement, ...rest } = props;
	return {
		display: "inline-flex",
		justifyContent: "center",
		alignItems: "center",
		position: "absolute",
		insetBlockStart: map(placement, (v) => {
			const [side] = v.split("-");
			return {
				top: offsetY,
				middle: "50%",
				bottom: "auto"
			}[side];
		}),
		insetBlockEnd: map(placement, (v) => {
			const [side] = v.split("-");
			return {
				top: "auto",
				middle: "50%",
				bottom: offsetY
			}[side];
		}),
		insetInlineStart: map(placement, (v) => {
			const [, align] = v.split("-");
			return {
				start: offsetX,
				center: "50%",
				end: "auto"
			}[align];
		}),
		insetInlineEnd: map(placement, (v) => {
			const [, align] = v.split("-");
			return {
				start: "auto",
				center: "50%",
				end: offsetX
			}[align];
		}),
		translate: map(placement, (v) => {
			const [side, align] = v.split("-");
			return `${{
				start: "-50%",
				center: "-50%",
				end: "50%"
			}[align]} ${{
				top: "-50%",
				middle: "-50%",
				bottom: "50%"
			}[side]}`;
		}),
		...rest
	};
},defaultValues(props) {
	const offset = props.offset || "0";
	return {
		offset,
		offsetX: offset,
		offsetY: offset,
		placement: "top-end"
	};
}}

export function floatRaw(styles) {
  const s = getPatternStyles(floatConfig, styles || {})
  return floatConfig.transform(s, patternFns)
}

export const float = /* @__PURE__ */ Object.assign(function float(styles = {}) {
  return css(floatRaw(styles))
}, { raw: floatRaw })