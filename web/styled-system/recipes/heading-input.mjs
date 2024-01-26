import { memo, splitProps } from '../helpers.mjs';
import { createRecipe, mergeRecipes } from './create-recipe.mjs';

const headingInputFn = /* @__PURE__ */ createRecipe('headingInput', {}, [])

const headingInputVariantMap = {}

const headingInputVariantKeys = Object.keys(headingInputVariantMap)

export const headingInput = /* @__PURE__ */ Object.assign(memo(headingInputFn), {
  __recipe__: true,
  __name__: 'headingInput',
  raw: (props) => props,
  variantKeys: headingInputVariantKeys,
  variantMap: headingInputVariantMap,
  merge(recipe) {
    return mergeRecipes(this, recipe)
  },
  splitVariantProps(props) {
    return splitProps(props, headingInputVariantKeys)
  },
})