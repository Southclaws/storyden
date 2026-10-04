import { createElement, forwardRef } from 'react';
import { center } from '../patterns/center';
import { splitProps } from '../helpers';
import { styled } from './factory';

export const Center = /* @__PURE__ */ forwardRef(function Center(props, ref) {
  const [patternProps, restProps] = splitProps(props, center.propKeys)
  const styleProps = center.raw(patternProps)
  const mergedProps = { ref, ...styleProps, ...restProps }
  return createElement(styled["div"], mergedProps)
})