import { createElement, forwardRef } from 'react';
import { splitProps } from '../helpers';
import { flexRaw } from '../patterns/flex';
import { styled } from './factory';

export const Flex = /* @__PURE__ */ forwardRef(function Flex(props, ref) {
  const [patternProps, restProps] = splitProps(props, ["align","basis","direction","grow","justify","shrink","wrap"])
  const styleProps = flexRaw(patternProps)
  const mergedProps = { ref, ...styleProps, ...restProps }
  return createElement(styled["div"], mergedProps)
})