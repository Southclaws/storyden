import { createElement, forwardRef } from 'react';
import { circle } from '../patterns/circle';
import { splitProps } from '../helpers';
import { styled } from './factory';

export const Circle = /* @__PURE__ */ forwardRef(function Circle(props, ref) {
  const [patternProps, restProps] = splitProps(props, circle.propKeys)
  const styleProps = circle.raw(patternProps)
  const mergedProps = { ref, ...styleProps, ...restProps }
  return createElement(styled["div"], mergedProps)
})