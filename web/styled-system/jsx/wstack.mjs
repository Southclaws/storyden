import { createElement, forwardRef } from 'react'

import { splitProps } from '../helpers.mjs';
import { getWstackStyle } from '../patterns/wstack.mjs';
import { styled } from './factory.mjs';

export const WStack = /* @__PURE__ */ forwardRef(function WStack(props, ref) {
  const [patternProps, restProps] = splitProps(props, [])

const styleProps = getWstackStyle(patternProps)
const mergedProps = { ref, ...styleProps, ...restProps }

return createElement(styled.div, mergedProps)
  })