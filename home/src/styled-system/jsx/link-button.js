import { createElement, forwardRef } from 'react';
import { linkButton } from '../patterns/link-button';
import { splitProps } from '../helpers';
import { styled } from './factory';

export const LinkButton = /* @__PURE__ */ forwardRef(function LinkButton(props, ref) {
  const [patternProps, restProps] = splitProps(props, linkButton.propKeys)
  const styleProps = linkButton.raw(patternProps)
  const mergedProps = { ref, ...styleProps, ...restProps }
  return createElement(styled["div"], mergedProps)
})