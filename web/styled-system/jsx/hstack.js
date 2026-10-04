import { createElement, forwardRef } from 'react';
import { hstack } from '../patterns/hstack';
import { splitProps } from '../helpers';
import { styled } from './factory';

export const HStack = /* @__PURE__ */ forwardRef(function HStack(props, ref) {
  const [patternProps, restProps] = splitProps(props, hstack.propKeys)
  const styleProps = hstack.raw(patternProps)
  const mergedProps = { ref, ...styleProps, ...restProps }
  return createElement(styled["div"], mergedProps)
})