import { createElement, forwardRef } from 'react';
import { float } from '../patterns/float';
import { splitProps } from '../helpers';
import { styled } from './factory';

export const Float = /* @__PURE__ */ forwardRef(function Float(props, ref) {
  const [patternProps, restProps] = splitProps(props, float.propKeys)
  const styleProps = float.raw(patternProps)
  const mergedProps = { ref, ...styleProps, ...restProps }
  return createElement(styled["div"], mergedProps)
})