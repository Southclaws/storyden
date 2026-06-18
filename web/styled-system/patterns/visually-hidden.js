import { getPatternStyles, patternFns } from './runtime';
import { css } from '../css/index';

const visuallyHiddenConfig = {transform(props) {
	return {
		srOnly: true,
		...props
	};
}}

export function visuallyHiddenRaw(styles) {
  const s = getPatternStyles(visuallyHiddenConfig, styles || {})
  return visuallyHiddenConfig.transform(s, patternFns)
}

export const visuallyHidden = /* @__PURE__ */ Object.assign(function visuallyHidden(styles = {}) {
  return css(visuallyHiddenRaw(styles))
}, { raw: visuallyHiddenRaw })