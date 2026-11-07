import { mapObject, withDefaults } from '../helpers';

const CSS_FUNCTION_REGEX = /^(min|max|clamp|calc)\(.*\)/
const CSS_VAR_REGEX = /^var\(--.+\)$/
const CSS_UNIT_REGEX = /^[+-]?[0-9]*.?[0-9]+(?:[eE][+-]?[0-9]+)?(?:cm|mm|Q|in|pc|pt|px|em|ex|ch|rem|lh|rlh|vw|vh|vmin|vmax|vb|vi|svw|svh|lvw|lvh|dvw|dvh|cqw|cqh|cqi|cqb|cqmin|cqmax|%)$/

export function isCssFunction(v) {
  return typeof v === "string" && CSS_FUNCTION_REGEX.test(v)
}

export function isCssVar(v) {
  return typeof v === "string" && CSS_VAR_REGEX.test(v)
}

export function isCssUnit(v) {
  return typeof v === "string" && CSS_UNIT_REGEX.test(v)
}

export const patternFns = { map: mapObject, isCssFunction, isCssVar, isCssUnit }

export function getPatternStyles(pattern, styles) {
  if (!pattern?.defaultValues) return styles
  const defaults = typeof pattern.defaultValues === "function" ? pattern.defaultValues(styles) : pattern.defaultValues
  return withDefaults(defaults, styles)
}