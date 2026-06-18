import { createElement, forwardRef } from 'react';
import { splitProps } from '../helpers';
import { visuallyHiddenRaw } from '../patterns/visually-hidden';
import { styled } from './factory';

export const VisuallyHidden = /* @__PURE__ */ forwardRef(function VisuallyHidden(props, ref) {
  const [patternProps, restProps] = splitProps(props, [])
  const styleProps = visuallyHiddenRaw(patternProps)
  const mergedProps = { ref, ...styleProps, ...restProps }
  return createElement(styled["div"], mergedProps)
})