import { createElement, forwardRef } from 'react'
import { styled } from './factory.mjs';
import { getFloatingStyle } from '../patterns/floating.mjs';

export const Floating = /* @__PURE__ */ forwardRef(function Floating(props, ref) {
  const styleProps = getFloatingStyle()
return createElement(styled.div, { ref, ...styleProps, ...props })
})