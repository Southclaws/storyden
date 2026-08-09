import { memo, splitProps } from '../helpers.mjs';
import { createRecipe, mergeRecipes } from './create-recipe.mjs';

const cardRowsFn = /* @__PURE__ */ createRecipe('card-rows', {}, [])

const cardRowsVariantMap = {}

const cardRowsVariantKeys = Object.keys(cardRowsVariantMap)

export const cardRows = /* @__PURE__ */ Object.assign(memo(cardRowsFn.recipeFn), {
  __recipe__: true,
  __name__: 'cardRows',
  __getCompoundVariantCss__: cardRowsFn.__getCompoundVariantCss__,
  raw: (props) => props,
  variantKeys: cardRowsVariantKeys,
  variantMap: cardRowsVariantMap,
  merge(recipe) {
    return mergeRecipes(this, recipe)
  },
  splitVariantProps(props) {
    return splitProps(props, cardRowsVariantKeys)
  },
  getVariantProps: cardRowsFn.getVariantProps,
})