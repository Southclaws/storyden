import { mapObject } from '../helpers.mjs';
import { css } from '../css/index.mjs';

const FrostedGlassConfig = {
transform() {
  return {
    backgroundColor: "bg.opaque",
    backdropBlur: "frosted",
    backdropFilter: "auto"
  };
}}

export const getFrostedGlassStyle = (styles = {}) => FrostedGlassConfig.transform(styles, { map: mapObject })

export const FrostedGlass = (styles) => css(getFrostedGlassStyle(styles))
FrostedGlass.raw = getFrostedGlassStyle