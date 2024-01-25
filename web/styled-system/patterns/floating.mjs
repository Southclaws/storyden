import { mapObject } from '../helpers.mjs';
import { css } from '../css/index.mjs';

const FloatingConfig = {
transform() {
  return {
    backgroundColor: "bg.opaque",
    backdropBlur: "frosted",
    backdropFilter: "auto",
    borderRadius: "lg",
    boxShadow: "sm"
  };
}}

export const getFloatingStyle = (styles = {}) => FloatingConfig.transform(styles, { map: mapObject })

export const Floating = (styles) => css(getFloatingStyle(styles))
Floating.raw = getFloatingStyle