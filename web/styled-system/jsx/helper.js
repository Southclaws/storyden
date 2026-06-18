import { css } from '../css/index';

export const composeShouldForwardProps = (tag, shouldForwardProp) => {
  if (!tag.__shouldForwardProps__ || !shouldForwardProp) return shouldForwardProp
  return (prop) => tag.__shouldForwardProps__(prop) && shouldForwardProp(prop)
}

export const composeCvaFn = (cvaA, cvaB) => {
  if (cvaA && !cvaB) return cvaA
  if (!cvaA && cvaB) return cvaB
  if ((cvaA.__cva__ && cvaB.__cva__) || (cvaA.__recipe__ && cvaB.__recipe__)) return cvaA.merge(cvaB)
  const error = new TypeError('Cannot merge cva with recipe. Please use either cva or recipe.')
  TypeError.captureStackTrace?.(error)
  throw error
}

export const getDisplayName = (Component) => {
  if (typeof Component === 'string') return Component
  return Component?.displayName || Component?.name || 'Component'
}

const htmlPropsMap = {
  htmlWidth: 'width',
  htmlHeight: 'height',
  htmlTranslate: 'translate',
  htmlContent: 'content',
}
const hasOwn = Object.prototype.hasOwnProperty

export function splitJsxProps(props, shouldForwardProp, variantSet, isCssProperty, skipClass) {
  let htmlProps
  let forwardedProps
  let variantProps
  let propStyles
  let cssStyles
  let elementProps
  const keys = Object.keys(props)
  for (let i = 0; i < keys.length; i++) {
    const key = keys[i]
    const value = props[key]
    if (value === void 0) continue
    if (key === 'className' || (skipClass && key === 'class') || key === 'as' || key === 'unstyled' || key === 'children') continue
    if (hasOwn.call(htmlPropsMap, key)) {
      htmlProps ||= Object.create(null)
      htmlProps[htmlPropsMap[key]] = value
    } else if (shouldForwardProp(key)) {
      forwardedProps ||= Object.create(null)
      forwardedProps[key] = value
    } else if (variantSet.has(key)) {
      variantProps ||= Object.create(null)
      variantProps[key] = value
    } else if (key === 'css') {
      cssStyles = value
    } else if (isCssProperty(key)) {
      (propStyles ||= {})[key] = value
    } else {
      elementProps ||= Object.create(null)
      elementProps[key] = value
    }
  }
  return [htmlProps, forwardedProps, variantProps || {}, propStyles, cssStyles, elementProps]
}

export function serializeSplitStyles(propStyles, cssStyles, baseStyles) {
  if (baseStyles !== void 0) {
    return propStyles ? cssStyles !== void 0 ? css(baseStyles, propStyles, cssStyles) : css(baseStyles, propStyles) : css(baseStyles, cssStyles)
  }
  return propStyles ? cssStyles !== void 0 ? css(propStyles, cssStyles) : css(propStyles) : css(cssStyles)
}

export function splitStyleProps(styleProps) {
  let propStyles
  let cssStyles
  const keys = Object.keys(styleProps)
  for (let i = 0; i < keys.length; i++) {
    const key = keys[i]
    const value = styleProps[key]
    if (value === void 0) continue
    if (key === 'css') cssStyles = value
    else (propStyles ||= {})[key] = value
  }
  return [propStyles, cssStyles]
}