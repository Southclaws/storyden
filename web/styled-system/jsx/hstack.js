import { createElement, forwardRef } from 'react';
import { splitProps } from '../helpers';
import { hstackRaw } from '../patterns/hstack';
import { styled } from './factory';

export const HStack = /* @__PURE__ */ forwardRef(function HStack(props, ref) {
  const [patternProps, restProps] = splitProps(props, ["gap","justify"])
  const styleProps = hstackRaw(patternProps)
  const mergedProps = { ref, ...styleProps, ...restProps }
  return createElement(styled["div"], mergedProps)
})