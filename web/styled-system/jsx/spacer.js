import { createElement, forwardRef } from 'react';
import { spacer } from '../patterns/spacer';
import { splitProps } from '../helpers';
import { styled } from './factory';

export const Spacer = /* @__PURE__ */ forwardRef(function Spacer(props, ref) {
  const [patternProps, restProps] = splitProps(props, spacer.propKeys)
  const styleProps = spacer.raw(patternProps)
  const mergedProps = { ref, ...styleProps, ...restProps }
  return createElement(styled["div"], mergedProps)
})