import { createElement, forwardRef } from 'react';
import { splitProps } from '../helpers';
import { CardBoxRaw } from '../patterns/card-box';
import { styled } from './factory';

export const CardBox = /* @__PURE__ */ forwardRef(function CardBox(props, ref) {
  const [patternProps, restProps] = splitProps(props, ["display","kind"])
  const styleProps = CardBoxRaw(patternProps)
  const mergedProps = { ref, ...styleProps, ...restProps }
  return createElement(styled["div"], mergedProps)
})