import { createElement, forwardRef } from 'react'

import { splitProps } from '../helpers.mjs';
import { getCardStyle } from '../patterns/card.mjs';
import { styled } from './factory.mjs';

export const Card = /* @__PURE__ */ forwardRef(function Card(props, ref) {
  const [patternProps, restProps] = splitProps(props, ["kind","display"])

const styleProps = getCardStyle(patternProps)
const mergedProps = { ref, ...styleProps, ...restProps }

return createElement(styled.div, mergedProps)
  })