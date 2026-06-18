import { createElement, forwardRef } from 'react';
import { splitProps } from '../helpers';
import { squareRaw } from '../patterns/square';
import { styled } from './factory';

export const Square = /* @__PURE__ */ forwardRef(function Square(props, ref) {
  const [patternProps, restProps] = splitProps(props, ["size"])
  const styleProps = squareRaw(patternProps)
  const mergedProps = { ref, ...styleProps, ...restProps }
  return createElement(styled["div"], mergedProps)
})