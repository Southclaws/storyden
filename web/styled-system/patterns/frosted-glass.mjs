import { mapObject } from '../helpers.mjs';
import { css } from '../css/index.mjs';

const FrostedGlassConfig = {
transform(props) {
  return {
    backgroundColor: "whiteAlpha.800",
    backdropBlur: "frosted",
    backdropFilter: "auto",
    boxShadow: "sm",
    borderRadius: "lg"
  };
}}

export const getFrostedGlassStyle = (styles = {}) => FrostedGlassConfig.transform(styles, { map: mapObject })

export const FrostedGlass = (styles) => css(getFrostedGlassStyle(styles))
FrostedGlass.raw = getFrostedGlassStyle