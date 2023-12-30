import { createElement, forwardRef } from 'react'
import { styled } from './factory.mjs';
import { getCardStyle } from '../patterns/card.mjs';

export const Card = /* @__PURE__ */ forwardRef(function Card(props, ref) {
  const { kind, ...restProps } = props
const styleProps = getCardStyle({kind})
return createElement(styled.div, { ref, ...styleProps, ...restProps })
})