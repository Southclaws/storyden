import { createElement, forwardRef } from 'react';
import { wstack } from '../patterns/wstack';
import { splitProps } from '../helpers';
import { styled } from './factory';

export const WStack = /* @__PURE__ */ forwardRef(function WStack(props, ref) {
  const [patternProps, restProps] = splitProps(props, wstack.propKeys)
  const styleProps = wstack.raw(patternProps)
  const mergedProps = { ref, ...styleProps, ...restProps }
  return createElement(styled["div"], mergedProps)
})