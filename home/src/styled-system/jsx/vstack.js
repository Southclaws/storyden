import { createElement, forwardRef } from 'react';
import { vstack } from '../patterns/vstack';
import { splitProps } from '../helpers';
import { styled } from './factory';

export const VStack = /* @__PURE__ */ forwardRef(function VStack(props, ref) {
  const [patternProps, restProps] = splitProps(props, vstack.propKeys)
  const styleProps = vstack.raw(patternProps)
  const mergedProps = { ref, ...styleProps, ...restProps }
  return createElement(styled["div"], mergedProps)
})