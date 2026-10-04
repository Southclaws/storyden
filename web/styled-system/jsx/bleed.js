import { createElement, forwardRef } from 'react';
import { bleed } from '../patterns/bleed';
import { splitProps } from '../helpers';
import { styled } from './factory';

export const Bleed = /* @__PURE__ */ forwardRef(function Bleed(props, ref) {
  const [patternProps, restProps] = splitProps(props, bleed.propKeys)
  const styleProps = bleed.raw(patternProps)
  const mergedProps = { ref, ...styleProps, ...restProps }
  return createElement(styled["div"], mergedProps)
})