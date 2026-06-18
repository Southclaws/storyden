import { createElement, forwardRef } from 'react';
import { splitProps } from '../helpers';
import { menuItemColorPaletteRaw } from '../patterns/menu-item-color-palette';
import { styled } from './factory';

export const MenuItemColorPalette = /* @__PURE__ */ forwardRef(function MenuItemColorPalette(props, ref) {
  const [patternProps, restProps] = splitProps(props, [])
  const styleProps = menuItemColorPaletteRaw(patternProps)
  const mergedProps = { ref, ...styleProps, ...restProps }
  return createElement(styled["div"], mergedProps)
})