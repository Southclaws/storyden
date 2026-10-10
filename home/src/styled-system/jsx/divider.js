import { createElement, forwardRef } from 'react';
import { divider } from '../patterns/divider';
import { splitProps } from '../helpers';
import { styled } from './factory';

export const Divider = /* @__PURE__ */ forwardRef(function Divider(props, ref) {
  const [patternProps, restProps] = splitProps(props, divider.propKeys)
  const styleProps = divider.raw(patternProps)
  const mergedProps = { ref, ...styleProps, ...restProps }
  return createElement(styled["div"], mergedProps)
})