import { getPatternStyles, patternFns } from './runtime';
import { css } from '../css/index';

const bleedConfig = {transform(props, { map, isCssUnit, isCssVar }) {
	const { inline, block, ...rest } = props;
	const valueFn = (v) => isCssUnit(v) || isCssVar(v) ? v : `token(spacing.${v}, ${v})`;
	return {
		"--bleed-x": map(inline, valueFn),
		"--bleed-y": map(block, valueFn),
		marginInline: "calc(var(--bleed-x, 0) * -1)",
		marginBlock: "calc(var(--bleed-y, 0) * -1)",
		...rest
	};
},defaultValues:{inline:'0',block:'0'}}

export function bleedRaw(styles) {
  const s = getPatternStyles(bleedConfig, styles || {})
  return bleedConfig.transform(s, patternFns)
}

export const bleed = /* @__PURE__ */ Object.assign(function bleed(styles = {}) {
  return css(bleedRaw(styles))
}, { raw: bleedRaw })