import { createElement, forwardRef } from 'react';
import { cq } from '../patterns/cq';
import { splitProps } from '../helpers';
import { styled } from './factory';

export const Cq = /* @__PURE__ */ forwardRef(function Cq(props, ref) {
  const [patternProps, restProps] = splitProps(props, cq.propKeys)
  const styleProps = cq.raw(patternProps)
  const mergedProps = { ref, ...styleProps, ...restProps }
  return createElement(styled["div"], mergedProps)
})