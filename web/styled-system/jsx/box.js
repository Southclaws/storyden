import { createElement, forwardRef } from 'react';
import { splitProps } from '../helpers';
import { boxRaw } from '../patterns/box';
import { styled } from './factory';

export const Box = /* @__PURE__ */ forwardRef(function Box(props, ref) {
  const [patternProps, restProps] = splitProps(props, [])
  const styleProps = boxRaw(patternProps)
  const mergedProps = { ref, ...styleProps, ...restProps }
  return createElement(styled["div"], mergedProps)
})