import { createElement, forwardRef } from 'react'

import { splitProps } from '../helpers.mjs';
import { getFrostedGlassStyle } from '../patterns/frosted-glass.mjs';
import { styled } from './factory.mjs';

export const FrostedGlass = /* @__PURE__ */ forwardRef(function FrostedGlass(props, ref) {
  const [patternProps, restProps] = splitProps(props, [])

const styleProps = getFrostedGlassStyle(patternProps)
const mergedProps = { ref, ...styleProps, ...restProps }

return createElement(styled.div, mergedProps)
  })