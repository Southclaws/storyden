import { createElement, forwardRef } from 'react';
import { wrap } from '../patterns/wrap';
import { splitProps } from '../helpers';
import { styled } from './factory';

export const Wrap = /* @__PURE__ */ forwardRef(function Wrap(props, ref) {
  const [patternProps, restProps] = splitProps(props, wrap.propKeys)
  const styleProps = wrap.raw(patternProps)
  const mergedProps = { ref, ...styleProps, ...restProps }
  return createElement(styled["div"], mergedProps)
})