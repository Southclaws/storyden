import { createElement, forwardRef } from 'react';
import { splitProps } from '../helpers';
import { floatRaw } from '../patterns/float';
import { styled } from './factory';

export const Float = /* @__PURE__ */ forwardRef(function Float(props, ref) {
  const [patternProps, restProps] = splitProps(props, ["offset","offsetX","offsetY","placement"])
  const styleProps = floatRaw(patternProps)
  const mergedProps = { ref, ...styleProps, ...restProps }
  return createElement(styled["div"], mergedProps)
})