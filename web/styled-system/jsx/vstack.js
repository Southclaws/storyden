import { createElement, forwardRef } from 'react';
import { splitProps } from '../helpers';
import { vstackRaw } from '../patterns/vstack';
import { styled } from './factory';

export const VStack = /* @__PURE__ */ forwardRef(function VStack(props, ref) {
  const [patternProps, restProps] = splitProps(props, ["gap","justify"])
  const styleProps = vstackRaw(patternProps)
  const mergedProps = { ref, ...styleProps, ...restProps }
  return createElement(styled["div"], mergedProps)
})