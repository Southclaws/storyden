import { createCssRuntime, getCompoundVariantClassNames, getSlotCompoundVariant, memo, splitProps, toHash, uniq, withDefaults, withoutSpace } from '../helpers';
import { breakpointKeys, finalizeConditions, sortConditions } from '../css/conditions';
import { cx } from '../css/cx';

function normalize(config) {
  const variantMap = config.variantMap ?? {}
  return {
    name: config.name,
    className: config.className ?? config.name,
    slots: config.slots ?? [],
    variantMap,
    variantKeys: Object.keys(variantMap),
    defaults: config.defaultVariants ?? {},
    compounds: config.compoundVariants ?? [],
  }
}

export function createRecipe(config) {
  const { name, className, variantMap, variantKeys, defaults, compounds } = normalize(config)
  const classPrefix = null

  const { serializeCss: recipeCss } = createCssRuntime({
    hash: false,
    conditions: {
      shift: sortConditions,
      finalize: finalizeConditions,
      breakpoints: { keys: breakpointKeys },
    },
    utility: {
      prefix: classPrefix,
      toHash,
      transform(prop, value) {
        return { className: value === "__ignore__" ? className : `${className}--${prop}_${withoutSpace(value)}` }
      },
    },
  })
  const formatClassName = (name) => {
    const next = false ? toHash(name) : name
    return classPrefix ? `${classPrefix}-${next}` : next
  }

  function resolve(props = {}) {
    return { [className]: "__ignore__", ...withDefaults(defaults, props) }
  }

  function assertNoConditions(props) {
    if (compounds.length === 0) return
    for (const key of variantKeys) {
      const value = props[key]
      if (value != null && typeof value === "object") {
        throw new Error(`[recipe:${className}:${key}] Conditions are not supported when using compound variants.`)
      }
    }
  }

  function compoundClasses(props) {
    return getCompoundVariantClassNames(compounds, resolve(props), formatClassName)
  }

  const recipe = attach(memo(function recipeFn(props = {}, withCompoundVariants = true) {
    assertNoConditions(props)
    const recipeClass = recipeCss(resolve(props))
    if (!withCompoundVariants) return recipeClass
    const compoundsClass = compoundClasses(props)
    return cx(recipeClass, compoundsClass)
  }), name, variantKeys, variantMap, resolve)
  recipe.__recipe__ = true
  recipe.__getCompoundVariantClasses__ = compoundClasses
  recipe.merge = function merge(other) {
    return mergeRecipes(recipe, other)
  }
  return recipe
}

export function createSlotRecipe(config) {
  const { name, className, slots, variantMap, variantKeys, defaults, compounds } = normalize(config)

  const slotFns = slots.map(function toSlotRecipe(slot) {
    return [slot, createRecipe({
      name,
      className: `${className}__${slot}`,
      variantMap,
      defaultVariants: defaults,
      compoundVariants: getSlotCompoundVariant(compounds, slot),
    })]
  })

  const recipe = memo(function slotRecipeFn(props = {}) {
    const result = {}
    for (const [slot, slotFn] of slotFns) result[slot] = slotFn(props)
    return result
  })
  attach(recipe, name, variantKeys, variantMap, function getVariantProps(props = {}) {
    return withDefaults(defaults, props)
  })
  recipe.__recipe__ = false
  recipe.classNameMap = {}
  return recipe
}

function mergeRecipes(recipeA, recipeB) {
  if (recipeA && !recipeB) return recipeA
  if (!recipeA && recipeB) return recipeB
  function merged(...args) {
    const classA = recipeA(...args)
    const classB = recipeB(...args)
    return classA && classB ? `${classA} ${classB}` : classA || classB
  }
  const variantKeys = uniq(recipeA.variantKeys, recipeB.variantKeys)
  const variantMap = {}
  for (const key of variantKeys) variantMap[key] = uniq(recipeA.variantMap[key], recipeB.variantMap[key])
  attach(merged, `${recipeA.__name__} ${recipeB.__name__}`, variantKeys, variantMap, function getVariantProps(props) {
    return props
  })
  merged.__recipe__ = true
  return merged
}

function attach(recipe, name, variantKeys, variantMap, getVariantProps) {
  recipe.__name__ = name
  recipe.raw = function raw(props) {
    return props
  }
  recipe.variantKeys = variantKeys
  recipe.variantMap = variantMap
  recipe.splitVariantProps = function splitVariantProps(props) {
    return splitProps(props, variantKeys)
  }
  recipe.getVariantProps = getVariantProps
  return recipe
}