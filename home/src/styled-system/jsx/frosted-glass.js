import { createElement, forwardRef } from 'react';
import { frostedGlass } from '../patterns/frosted-glass';
import { splitProps } from '../helpers';
import { styled } from './factory';

export const FrostedGlass = /* @__PURE__ */ forwardRef(function FrostedGlass(props, ref) {
  const [patternProps, restProps] = splitProps(props, frostedGlass.propKeys)
  const styleProps = frostedGlass.raw(patternProps)
  const mergedProps = { ref, ...styleProps, ...restProps }
  return createElement(styled["div"], mergedProps)
})