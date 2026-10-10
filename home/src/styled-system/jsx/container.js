import { createElement, forwardRef } from 'react';
import { container } from '../patterns/container';
import { splitProps } from '../helpers';
import { styled } from './factory';

export const Container = /* @__PURE__ */ forwardRef(function Container(props, ref) {
  const [patternProps, restProps] = splitProps(props, container.propKeys)
  const styleProps = container.raw(patternProps)
  const mergedProps = { ref, ...styleProps, ...restProps }
  return createElement(styled["div"], mergedProps)
})