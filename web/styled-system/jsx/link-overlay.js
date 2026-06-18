import { createElement, forwardRef } from 'react';
import { splitProps } from '../helpers';
import { linkOverlayRaw } from '../patterns/link-overlay';
import { styled } from './factory';

export const LinkOverlay = /* @__PURE__ */ forwardRef(function LinkOverlay(props, ref) {
  const [patternProps, restProps] = splitProps(props, [])
  const styleProps = linkOverlayRaw(patternProps)
  const mergedProps = { ref, ...styleProps, ...restProps }
  return createElement(styled["a"], mergedProps)
})