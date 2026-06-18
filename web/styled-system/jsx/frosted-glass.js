import { createElement, forwardRef } from 'react';
import { splitProps } from '../helpers';
import { FrostedGlassRaw } from '../patterns/frosted-glass';
import { styled } from './factory';

export const FrostedGlass = /* @__PURE__ */ forwardRef(function FrostedGlass(props, ref) {
  const [patternProps, restProps] = splitProps(props, [])
  const styleProps = FrostedGlassRaw(patternProps)
  const mergedProps = { ref, ...styleProps, ...restProps }
  return createElement(styled["div"], mergedProps)
})