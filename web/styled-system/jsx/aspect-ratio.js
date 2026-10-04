import { createElement, forwardRef } from 'react';
import { aspectRatio } from '../patterns/aspect-ratio';
import { splitProps } from '../helpers';
import { styled } from './factory';

export const AspectRatio = /* @__PURE__ */ forwardRef(function AspectRatio(props, ref) {
  const [patternProps, restProps] = splitProps(props, aspectRatio.propKeys)
  const styleProps = aspectRatio.raw(patternProps)
  const mergedProps = { ref, ...styleProps, ...restProps }
  return createElement(styled["div"], mergedProps)
})