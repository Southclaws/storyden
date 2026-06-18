import { mapObject, withDefaults } from '../helpers';

export function isCssFunction(v) {
  return typeof v === "string" && /^(min|max|clamp|calc)\(.*\)/.test(v)
}

export function isCssVar(v) {
  return typeof v === "string" && /^var\(--.+\)$/.test(v)
}

export function isCssUnit(v) {
  return typeof v === "string" && /^[+-]?[0-9]*.?[0-9]+(?:[eE][+-]?[0-9]+)?(?:cm|mm|Q|in|pc|pt|px|em|ex|ch|rem|lh|rlh|vw|vh|vmin|vmax|vb|vi|svw|svh|lvw|lvh|dvw|dvh|cqw|cqh|cqi|cqb|cqmin|cqmax|%)$/.test(v)
}

export const patternFns = { map: mapObject, isCssFunction, isCssVar, isCssUnit }

export function getPatternStyles(pattern, styles) {
  if (!pattern?.defaultValues) return styles
  const defaults = typeof pattern.defaultValues === "function" ? pattern.defaultValues(styles) : pattern.defaultValues
  return withDefaults(defaults, styles)
}