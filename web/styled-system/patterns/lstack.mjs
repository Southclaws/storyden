import { getPatternStyles, patternFns } from '../helpers.mjs';
import { css } from '../css/index.mjs';

const lstackConfig = {
transform() {
  return {
    display: "flex",
    gap: "3",
    flexDirection: "column",
    width: "full",
    alignItems: "start"
  };
}}

export const getLstackStyle = (styles = {}) => {
  const _styles = getPatternStyles(lstackConfig, styles)
  return lstackConfig.transform(_styles, patternFns)
}

export const lstack = (styles) => css(getLstackStyle(styles))
lstack.raw = getLstackStyle