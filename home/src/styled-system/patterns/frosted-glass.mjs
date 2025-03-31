import { getPatternStyles, patternFns } from '../helpers.mjs';
import { css } from '../css/index.mjs';

const FrostedGlassConfig = {
transform() {
  return {
    backgroundColor: "bg.opaque",
    backdropBlur: "frosted",
    backdropFilter: "auto"
  };
}}

export const getFrostedGlassStyle = (styles = {}) => {
  const _styles = getPatternStyles(FrostedGlassConfig, styles)
  return FrostedGlassConfig.transform(_styles, patternFns)
}

export const FrostedGlass = (styles) => css(getFrostedGlassStyle(styles))
FrostedGlass.raw = getFrostedGlassStyle