import { getPatternStyles, patternFns } from '../helpers.mjs';
import { css } from '../css/index.mjs';

const LStackConfig = {
transform() {
  return {
    display: "flex",
    gap: "3",
    flexDirection: "column",
    width: "full",
    alignItems: "start"
  };
}}

export const getLStackStyle = (styles = {}) => {
  const _styles = getPatternStyles(LStackConfig, styles)
  return LStackConfig.transform(_styles, patternFns)
}

export const LStack = (styles) => css(getLStackStyle(styles))
LStack.raw = getLStackStyle