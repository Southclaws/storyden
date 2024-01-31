import { createElement, forwardRef } from 'react'
import { mergeCss } from '../css/css.mjs';
import { splitProps } from '../helpers.mjs';
import { getLinkButtonStyle } from '../patterns/link-button.mjs';
import { styled } from './factory.mjs';

export const LinkButton = /* @__PURE__ */ forwardRef(function LinkButton(props, ref) {
  const [patternProps, restProps] = splitProps(props, [])

const styleProps = getLinkButtonStyle(patternProps)
const mergedProps = { ref, ...styleProps, ...restProps }

return createElement(styled.div, mergedProps)
  })