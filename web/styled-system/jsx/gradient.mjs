import { createElement, forwardRef } from 'react'
import { styled } from './factory.mjs';
import { getGradientStyle } from '../patterns/gradient.mjs';

export const Gradient = /* @__PURE__ */ forwardRef(function Gradient(props, ref) {
  const styleProps = getGradientStyle()
return createElement(styled.div, { ref, ...styleProps, ...props })
})