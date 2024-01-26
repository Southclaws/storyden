import { memo, splitProps } from '../helpers.mjs';
import { createRecipe, mergeRecipes } from './create-recipe.mjs';

const skeletonFn = /* @__PURE__ */ createRecipe('skeleton', {}, [])

const skeletonVariantMap = {}

const skeletonVariantKeys = Object.keys(skeletonVariantMap)

export const skeleton = /* @__PURE__ */ Object.assign(memo(skeletonFn), {
  __recipe__: true,
  __name__: 'skeleton',
  raw: (props) => props,
  variantKeys: skeletonVariantKeys,
  variantMap: skeletonVariantMap,
  merge(recipe) {
    return mergeRecipes(this, recipe)
  },
  splitVariantProps(props) {
    return splitProps(props, skeletonVariantKeys)
  },
})