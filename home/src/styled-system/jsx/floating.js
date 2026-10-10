import { createElement, forwardRef } from 'react';
import { floating } from '../patterns/floating';
import { splitProps } from '../helpers';
import { styled } from './factory';

export const Floating = /* @__PURE__ */ forwardRef(function Floating(props, ref) {
  const [patternProps, restProps] = splitProps(props, floating.propKeys)
  const styleProps = floating.raw(patternProps)
  const mergedProps = { ref, ...styleProps, ...restProps }
  return createElement(styled["div"], mergedProps)
})