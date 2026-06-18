import { createElement, forwardRef } from 'react';
import { splitProps } from '../helpers';
import { FloatingRaw } from '../patterns/floating';
import { styled } from './factory';

export const Floating = /* @__PURE__ */ forwardRef(function Floating(props, ref) {
  const [patternProps, restProps] = splitProps(props, [])
  const styleProps = FloatingRaw(patternProps)
  const mergedProps = { ref, ...styleProps, ...restProps }
  return createElement(styled["div"], mergedProps)
})