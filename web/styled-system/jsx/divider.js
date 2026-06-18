import { createElement, forwardRef } from 'react';
import { splitProps } from '../helpers';
import { dividerRaw } from '../patterns/divider';
import { styled } from './factory';

export const Divider = /* @__PURE__ */ forwardRef(function Divider(props, ref) {
  const [patternProps, restProps] = splitProps(props, ["color","orientation","thickness"])
  const styleProps = dividerRaw(patternProps)
  const mergedProps = { ref, ...styleProps, ...restProps }
  return createElement(styled["div"], mergedProps)
})