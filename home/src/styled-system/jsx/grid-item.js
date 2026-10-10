import { createElement, forwardRef } from 'react';
import { gridItem } from '../patterns/grid-item';
import { splitProps } from '../helpers';
import { styled } from './factory';

export const GridItem = /* @__PURE__ */ forwardRef(function GridItem(props, ref) {
  const [patternProps, restProps] = splitProps(props, gridItem.propKeys)
  const styleProps = gridItem.raw(patternProps)
  const mergedProps = { ref, ...styleProps, ...restProps }
  return createElement(styled["div"], mergedProps)
})