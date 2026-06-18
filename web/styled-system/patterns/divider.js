import { getPatternStyles, patternFns } from './runtime';
import { css } from '../css/index';

const dividerConfig = {transform(props, { map }) {
	const { orientation, thickness, color, ...rest } = props;
	return {
		"--thickness": thickness,
		width: map(orientation, (v) => v === "vertical" ? void 0 : "100%"),
		height: map(orientation, (v) => v === "horizontal" ? void 0 : "100%"),
		borderBlockEndWidth: map(orientation, (v) => v === "horizontal" ? "var(--thickness)" : void 0),
		borderInlineEndWidth: map(orientation, (v) => v === "vertical" ? "var(--thickness)" : void 0),
		borderColor: color,
		...rest
	};
},defaultValues:{orientation:'horizontal',thickness:'1px'}}

export function dividerRaw(styles) {
  const s = getPatternStyles(dividerConfig, styles || {})
  return dividerConfig.transform(s, patternFns)
}

export const divider = /* @__PURE__ */ Object.assign(function divider(styles = {}) {
  return css(dividerRaw(styles))
}, { raw: dividerRaw })