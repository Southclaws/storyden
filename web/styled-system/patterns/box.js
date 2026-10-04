import { getPatternStyles, patternFns } from './runtime';
import { memo } from '../helpers';
import { css } from '../css/index';

const boxConfig = {transform(props) {
	return props;
}}

export function boxRaw(styles) {
  const s = getPatternStyles(boxConfig, styles || {})
  return boxConfig.transform(s, patternFns)
}

export const box = /* @__PURE__ */ Object.assign(/* @__PURE__ */ memo(function box(styles = {}) {
  return css(boxRaw(styles))
}), { raw: boxRaw, propKeys: [] })