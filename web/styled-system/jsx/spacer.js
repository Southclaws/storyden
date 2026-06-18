import { createElement, forwardRef } from 'react';
import { splitProps } from '../helpers';
import { spacerRaw } from '../patterns/spacer';
import { styled } from './factory';

export const Spacer = /* @__PURE__ */ forwardRef(function Spacer(props, ref) {
  const [patternProps, restProps] = splitProps(props, ["size"])
  const styleProps = spacerRaw(patternProps)
  const mergedProps = { ref, ...styleProps, ...restProps }
  return createElement(styled["div"], mergedProps)
})