import { createElement, forwardRef } from 'react'

import { splitProps } from '../helpers.mjs';
import { getMenuItemColorPaletteStyle } from '../patterns/menu-item-color-palette.mjs';
import { styled } from './factory.mjs';

export const MenuItemColorPalette = /* @__PURE__ */ forwardRef(function MenuItemColorPalette(props, ref) {
  const [patternProps, restProps] = splitProps(props, [])

const styleProps = getMenuItemColorPaletteStyle(patternProps)
const mergedProps = { ref, ...styleProps, ...restProps }

return createElement(styled.div, mergedProps)
  })