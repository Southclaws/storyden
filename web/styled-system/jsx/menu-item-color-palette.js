import { createElement, forwardRef } from 'react';
import { menuItemColorPalette } from '../patterns/menu-item-color-palette';
import { splitProps } from '../helpers';
import { styled } from './factory';

export const MenuItemColorPalette = /* @__PURE__ */ forwardRef(function MenuItemColorPalette(props, ref) {
  const [patternProps, restProps] = splitProps(props, menuItemColorPalette.propKeys)
  const styleProps = menuItemColorPalette.raw(patternProps)
  const mergedProps = { ref, ...styleProps, ...restProps }
  return createElement(styled["div"], mergedProps)
})