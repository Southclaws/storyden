import { createElement, forwardRef } from 'react';
import { visuallyHidden } from '../patterns/visually-hidden';
import { splitProps } from '../helpers';
import { styled } from './factory';

export const VisuallyHidden = /* @__PURE__ */ forwardRef(function VisuallyHidden(props, ref) {
  const [patternProps, restProps] = splitProps(props, visuallyHidden.propKeys)
  const styleProps = visuallyHidden.raw(patternProps)
  const mergedProps = { ref, ...styleProps, ...restProps }
  return createElement(styled["div"], mergedProps)
})