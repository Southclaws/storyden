import { getPatternStyles, patternFns } from '../helpers.mjs';
import { css } from '../css/index.mjs';

const wstackConfig = {
transform() {
  return {
    display: "flex",
    flexDirection: "row",
    gap: "3",
    width: "full",
    justifyContent: "space-between"
  };
}}

export const getWstackStyle = (styles = {}) => {
  const _styles = getPatternStyles(wstackConfig, styles)
  return wstackConfig.transform(_styles, patternFns)
}

export const wstack = (styles) => css(getWstackStyle(styles))
wstack.raw = getWstackStyle