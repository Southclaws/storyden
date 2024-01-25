import { mapObject } from '../helpers.mjs';
import { css } from '../css/index.mjs';

const NavigationConfig = {
transform() {
  return {
    backgroundColor: "bg.opaque",
    backdropBlur: "frosted",
    backdropFilter: "auto",
    borderRadius: "lg",
    boxShadow: "sm"
  };
}}

export const getNavigationStyle = (styles = {}) => NavigationConfig.transform(styles, { map: mapObject })

export const Navigation = (styles) => css(getNavigationStyle(styles))
Navigation.raw = getNavigationStyle