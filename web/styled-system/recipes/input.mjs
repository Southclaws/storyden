import { memo, splitProps } from '../helpers.mjs';
import { createRecipe, mergeRecipes } from './create-recipe.mjs';

const inputFn = /* @__PURE__ */ createRecipe('input', {
  "size": "md"
}, [])

const inputVariantMap = {
  "size": [
    "2xs",
    "xs",
    "sm",
    "md",
    "lg",
    "xl",
    "2xl"
  ]
}

const inputVariantKeys = Object.keys(inputVariantMap)

export const input = /* @__PURE__ */ Object.assign(memo(inputFn), {
  __recipe__: true,
  __name__: 'input',
  raw: (props) => props,
  variantKeys: inputVariantKeys,
  variantMap: inputVariantMap,
  merge(recipe) {
    return mergeRecipes(this, recipe)
  },
  splitVariantProps(props) {
    return splitProps(props, inputVariantKeys)
  },
})