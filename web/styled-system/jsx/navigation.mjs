import { createElement, forwardRef } from 'react'
import { styled } from './factory.mjs';
import { getNavigationStyle } from '../patterns/navigation.mjs';

export const Navigation = /* @__PURE__ */ forwardRef(function Navigation(props, ref) {
  const styleProps = getNavigationStyle()
return createElement(styled.div, { ref, ...styleProps, ...props })
})