import { createElement, forwardRef } from 'react';
import { grid } from '../patterns/grid';
import { splitProps } from '../helpers';
import { styled } from './factory';

export const Grid = /* @__PURE__ */ forwardRef(function Grid(props, ref) {
  const [patternProps, restProps] = splitProps(props, grid.propKeys)
  const styleProps = grid.raw(patternProps)
  const mergedProps = { ref, ...styleProps, ...restProps }
  return createElement(styled["div"], mergedProps)
})