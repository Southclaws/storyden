import { createElement, forwardRef } from 'react';
import { square } from '../patterns/square';
import { splitProps } from '../helpers';
import { styled } from './factory';

export const Square = /* @__PURE__ */ forwardRef(function Square(props, ref) {
  const [patternProps, restProps] = splitProps(props, square.propKeys)
  const styleProps = square.raw(patternProps)
  const mergedProps = { ref, ...styleProps, ...restProps }
  return createElement(styled["div"], mergedProps)
})