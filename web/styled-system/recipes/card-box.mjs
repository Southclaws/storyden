import { memo, splitProps } from '../helpers.mjs';
import { createRecipe, mergeRecipes } from './create-recipe.mjs';

const cardBoxFn = /* @__PURE__ */ createRecipe('card-box', {
  "kind": "default"
}, [])

const cardBoxVariantMap = {
  "kind": [
    "default",
    "edge"
  ]
}

const cardBoxVariantKeys = Object.keys(cardBoxVariantMap)

export const cardBox = /* @__PURE__ */ Object.assign(memo(cardBoxFn.recipeFn), {
  __recipe__: true,
  __name__: 'cardBox',
  __getCompoundVariantCss__: cardBoxFn.__getCompoundVariantCss__,
  raw: (props) => props,
  variantKeys: cardBoxVariantKeys,
  variantMap: cardBoxVariantMap,
  merge(recipe) {
    return mergeRecipes(this, recipe)
  },
  splitVariantProps(props) {
    return splitProps(props, cardBoxVariantKeys)
  },
  getVariantProps: cardBoxFn.getVariantProps,
})