import { createElement, forwardRef } from 'react';
import { splitProps } from '../helpers';
import { aspectRatioRaw } from '../patterns/aspect-ratio';
import { styled } from './factory';

export const AspectRatio = /* @__PURE__ */ forwardRef(function AspectRatio(props, ref) {
  const [patternProps, restProps] = splitProps(props, ["ratio"])
  const styleProps = aspectRatioRaw(patternProps)
  const mergedProps = { ref, ...styleProps, ...restProps }
  return createElement(styled["div"], mergedProps)
})