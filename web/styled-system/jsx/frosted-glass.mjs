import { createElement, forwardRef } from 'react'
import { styled } from './factory.mjs';
import { getFrostedGlassStyle } from '../patterns/frosted-glass.mjs';

export const FrostedGlass = /* @__PURE__ */ forwardRef(function FrostedGlass(props, ref) {
  const styleProps = getFrostedGlassStyle()
return createElement(styled.div, { ref, ...styleProps, ...props })
})