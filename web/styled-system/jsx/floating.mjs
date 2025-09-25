import { createElement, forwardRef } from 'react'

import { splitProps } from '../helpers.mjs';
import { getFloatingStyle } from '../patterns/floating.mjs';
import { styled } from './factory.mjs';

export const Floating = /* @__PURE__ */ forwardRef(function Floating(props, ref) {
  const [patternProps, restProps] = splitProps(props, [])

const styleProps = getFloatingStyle(patternProps)
const mergedProps = { ref, ...styleProps, ...restProps }

return createElement(styled.div, mergedProps)
  })