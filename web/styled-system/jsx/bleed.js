import { createElement, forwardRef } from 'react';
import { splitProps } from '../helpers';
import { bleedRaw } from '../patterns/bleed';
import { styled } from './factory';

export const Bleed = /* @__PURE__ */ forwardRef(function Bleed(props, ref) {
  const [patternProps, restProps] = splitProps(props, ["block","inline"])
  const styleProps = bleedRaw(patternProps)
  const mergedProps = { ref, ...styleProps, ...restProps }
  return createElement(styled["div"], mergedProps)
})