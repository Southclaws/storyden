import { createElement, forwardRef } from 'react';
import { card } from '../patterns/card';
import { splitProps } from '../helpers';
import { styled } from './factory';

export const Card = /* @__PURE__ */ forwardRef(function Card(props, ref) {
  const [patternProps, restProps] = splitProps(props, card.propKeys)
  const styleProps = card.raw(patternProps)
  const mergedProps = { ref, ...styleProps, ...restProps }
  return createElement(styled["div"], mergedProps)
})