import { createElement, forwardRef } from 'react';
import { splitProps } from '../helpers';
import { cqRaw } from '../patterns/cq';
import { styled } from './factory';

export const Cq = /* @__PURE__ */ forwardRef(function Cq(props, ref) {
  const [patternProps, restProps] = splitProps(props, ["name","type"])
  const styleProps = cqRaw(patternProps)
  const mergedProps = { ref, ...styleProps, ...restProps }
  return createElement(styled["div"], mergedProps)
})