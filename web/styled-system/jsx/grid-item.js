import { createElement, forwardRef } from 'react';
import { splitProps } from '../helpers';
import { gridItemRaw } from '../patterns/grid-item';
import { styled } from './factory';

export const GridItem = /* @__PURE__ */ forwardRef(function GridItem(props, ref) {
  const [patternProps, restProps] = splitProps(props, ["colEnd","colSpan","colStart","rowEnd","rowSpan","rowStart"])
  const styleProps = gridItemRaw(patternProps)
  const mergedProps = { ref, ...styleProps, ...restProps }
  return createElement(styled["div"], mergedProps)
})