import { memo, splitProps } from '../helpers.mjs';
import { createRecipe, mergeRecipes } from './create-recipe.mjs';

const headingFn = /* @__PURE__ */ createRecipe('heading', {
  "size": "md"
}, [])

const headingVariantMap = {
  "size": [
    "xs",
    "sm",
    "md",
    "lg",
    "xl",
    "2xl"
  ]
}

const headingVariantKeys = Object.keys(headingVariantMap)

export const heading = /* @__PURE__ */ Object.assign(memo(headingFn), {
  __recipe__: true,
  __name__: 'heading',
  raw: (props) => props,
  variantKeys: headingVariantKeys,
  variantMap: headingVariantMap,
  merge(recipe) {
    return mergeRecipes(this, recipe)
  },
  splitVariantProps(props) {
    return splitProps(props, headingVariantKeys)
  },
})