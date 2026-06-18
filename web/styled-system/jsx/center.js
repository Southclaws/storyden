import { createElement, forwardRef } from 'react';
import { splitProps } from '../helpers';
import { centerRaw } from '../patterns/center';
import { styled } from './factory';

export const Center = /* @__PURE__ */ forwardRef(function Center(props, ref) {
  const [patternProps, restProps] = splitProps(props, ["inline"])
  const styleProps = centerRaw(patternProps)
  const mergedProps = { ref, ...styleProps, ...restProps }
  return createElement(styled["div"], mergedProps)
})