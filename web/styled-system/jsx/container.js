import { createElement, forwardRef } from 'react';
import { splitProps } from '../helpers';
import { containerRaw } from '../patterns/container';
import { styled } from './factory';

export const Container = /* @__PURE__ */ forwardRef(function Container(props, ref) {
  const [patternProps, restProps] = splitProps(props, [])
  const styleProps = containerRaw(patternProps)
  const mergedProps = { ref, ...styleProps, ...restProps }
  return createElement(styled["div"], mergedProps)
})