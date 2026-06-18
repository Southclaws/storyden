import { createElement, forwardRef } from 'react';
import { splitProps } from '../helpers';
import { wrapRaw } from '../patterns/wrap';
import { styled } from './factory';

export const Wrap = /* @__PURE__ */ forwardRef(function Wrap(props, ref) {
  const [patternProps, restProps] = splitProps(props, ["align","columnGap","gap","justify","rowGap"])
  const styleProps = wrapRaw(patternProps)
  const mergedProps = { ref, ...styleProps, ...restProps }
  return createElement(styled["div"], mergedProps)
})