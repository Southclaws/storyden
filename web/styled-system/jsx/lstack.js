import { createElement, forwardRef } from 'react';
import { splitProps } from '../helpers';
import { lstackRaw } from '../patterns/lstack';
import { styled } from './factory';

export const LStack = /* @__PURE__ */ forwardRef(function LStack(props, ref) {
  const [patternProps, restProps] = splitProps(props, [])
  const styleProps = lstackRaw(patternProps)
  const mergedProps = { ref, ...styleProps, ...restProps }
  return createElement(styled["div"], mergedProps)
})