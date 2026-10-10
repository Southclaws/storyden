import { createElement, forwardRef } from 'react';
import { stack } from '../patterns/stack';
import { splitProps } from '../helpers';
import { styled } from './factory';

export const Stack = /* @__PURE__ */ forwardRef(function Stack(props, ref) {
  const [patternProps, restProps] = splitProps(props, stack.propKeys)
  const styleProps = stack.raw(patternProps)
  const mergedProps = { ref, ...styleProps, ...restProps }
  return createElement(styled["div"], mergedProps)
})