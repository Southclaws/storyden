import { createElement, forwardRef } from 'react';
import { splitProps } from '../helpers';
import { wstackRaw } from '../patterns/wstack';
import { styled } from './factory';

export const WStack = /* @__PURE__ */ forwardRef(function WStack(props, ref) {
  const [patternProps, restProps] = splitProps(props, [])
  const styleProps = wstackRaw(patternProps)
  const mergedProps = { ref, ...styleProps, ...restProps }
  return createElement(styled["div"], mergedProps)
})