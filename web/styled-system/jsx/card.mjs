import { createElement, forwardRef } from 'react'
import { styled } from './factory.mjs';
import { getCardStyle } from '../patterns/card.mjs';

export const Card = /* @__PURE__ */ forwardRef(function Card(props, ref) {
  const { kind, display, ...restProps } = props
const styleProps = getCardStyle({kind, display})
return createElement(styled.div, { ref, ...styleProps, ...restProps })
})