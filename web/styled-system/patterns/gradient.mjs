import { mapObject } from '../helpers.mjs';
import { css } from '../css/index.mjs';

const GradientConfig = {
transform() {
  return {
    backgroundImage: `linear-gradient(90deg, var(--colors-bg-default), transparent)`
  };
}}

export const getGradientStyle = (styles = {}) => GradientConfig.transform(styles, { map: mapObject })

export const Gradient = (styles) => css(getGradientStyle(styles))
Gradient.raw = getGradientStyle