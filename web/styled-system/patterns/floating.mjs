import { getPatternStyles, patternFns } from '../helpers.mjs';
import { css } from '../css/index.mjs';

const FloatingConfig = {
transform() {
  return {
    backgroundColor: "bg.opaque",
    backdropBlur: "frosted",
    backdropFilter: "auto",
    boxShadow: "sm"
  };
}}

export const getFloatingStyle = (styles = {}) => {
  const _styles = getPatternStyles(FloatingConfig, styles)
  return FloatingConfig.transform(_styles, patternFns)
}

export const Floating = (styles) => css(getFloatingStyle(styles))
Floating.raw = getFloatingStyle