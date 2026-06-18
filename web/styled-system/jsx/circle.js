import { createElement, forwardRef } from 'react';
import { splitProps } from '../helpers';
import { circleRaw } from '../patterns/circle';
import { styled } from './factory';

export const Circle = /* @__PURE__ */ forwardRef(function Circle(props, ref) {
  const [patternProps, restProps] = splitProps(props, ["size"])
  const styleProps = circleRaw(patternProps)
  const mergedProps = { ref, ...styleProps, ...restProps }
  return createElement(styled["div"], mergedProps)
})