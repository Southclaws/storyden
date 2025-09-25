import { createElement, forwardRef } from 'react'

import { splitProps } from '../helpers.mjs';
import { getLstackStyle } from '../patterns/lstack.mjs';
import { styled } from './factory.mjs';

export const LStack = /* @__PURE__ */ forwardRef(function LStack(props, ref) {
  const [patternProps, restProps] = splitProps(props, [])

const styleProps = getLstackStyle(patternProps)
const mergedProps = { ref, ...styleProps, ...restProps }

return createElement(styled.div, mergedProps)
  })