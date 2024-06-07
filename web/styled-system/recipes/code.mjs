import { memo, splitProps } from '../helpers.mjs';
import { createRecipe, mergeRecipes } from './create-recipe.mjs';

const codeFn = /* @__PURE__ */ createRecipe('code', {
  "size": "md",
  "variant": "outline"
}, [])

const codeVariantMap = {
  "variant": [
    "outline",
    "ghost"
  ],
  "size": [
    "sm",
    "md",
    "lg"
  ]
}

const codeVariantKeys = Object.keys(codeVariantMap)

export const code = /* @__PURE__ */ Object.assign(memo(codeFn.recipeFn), {
  __recipe__: true,
  __name__: 'code',
  __getCompoundVariantCss__: codeFn.__getCompoundVariantCss__,
  raw: (props) => props,
  variantKeys: codeVariantKeys,
  variantMap: codeVariantMap,
  merge(recipe) {
    return mergeRecipes(this, recipe)
  },
  splitVariantProps(props) {
    return splitProps(props, codeVariantKeys)
  },
  getVariantProps: codeFn.getVariantProps,
})