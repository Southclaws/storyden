import { memo, splitProps } from '../helpers.mjs';
import { createRecipe, mergeRecipes } from './create-recipe.mjs';

const admonitionFn = /* @__PURE__ */ createRecipe('admonition', {}, [])

const admonitionVariantMap = {
  "kind": [
    "neutral",
    "success",
    "failure"
  ]
}

const admonitionVariantKeys = Object.keys(admonitionVariantMap)

export const admonition = /* @__PURE__ */ Object.assign(memo(admonitionFn), {
  __recipe__: true,
  __name__: 'admonition',
  raw: (props) => props,
  variantKeys: admonitionVariantKeys,
  variantMap: admonitionVariantMap,
  merge(recipe) {
    return mergeRecipes(this, recipe)
  },
  splitVariantProps(props) {
    return splitProps(props, admonitionVariantKeys)
  },
})