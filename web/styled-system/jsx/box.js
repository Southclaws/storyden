import { createElement, forwardRef } from 'react';
import { box } from '../patterns/box';
import { splitProps } from '../helpers';
import { styled } from './factory';

export const Box = /* @__PURE__ */ forwardRef(function Box(props, ref) {
  const [patternProps, restProps] = splitProps(props, box.propKeys)
  const styleProps = box.raw(patternProps)
  const mergedProps = { ref, ...styleProps, ...restProps }
  return createElement(styled["div"], mergedProps)
})