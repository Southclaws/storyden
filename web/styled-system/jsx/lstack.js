import { createElement, forwardRef } from 'react';
import { lstack } from '../patterns/lstack';
import { splitProps } from '../helpers';
import { styled } from './factory';

export const LStack = /* @__PURE__ */ forwardRef(function LStack(props, ref) {
  const [patternProps, restProps] = splitProps(props, lstack.propKeys)
  const styleProps = lstack.raw(patternProps)
  const mergedProps = { ref, ...styleProps, ...restProps }
  return createElement(styled["div"], mergedProps)
})