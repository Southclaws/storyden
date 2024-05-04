import { createElement, forwardRef } from 'react'

import { splitProps } from '../helpers.mjs';
import { getCardBoxStyle } from '../patterns/card-box.mjs';
import { styled } from './factory.mjs';

export const CardBox = /* @__PURE__ */ forwardRef(function CardBox(props, ref) {
  const [patternProps, restProps] = splitProps(props, ["kind","display"])

const styleProps = getCardBoxStyle(patternProps)
const mergedProps = { ref, ...styleProps, ...restProps }

return createElement(styled.div, mergedProps)
  })