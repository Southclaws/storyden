import { createElement, forwardRef } from 'react';
import { splitProps } from '../helpers';
import { stackRaw } from '../patterns/stack';
import { styled } from './factory';

export const Stack = /* @__PURE__ */ forwardRef(function Stack(props, ref) {
  const [patternProps, restProps] = splitProps(props, ["align","direction","gap","justify"])
  const styleProps = stackRaw(patternProps)
  const mergedProps = { ref, ...styleProps, ...restProps }
  return createElement(styled["div"], mergedProps)
})