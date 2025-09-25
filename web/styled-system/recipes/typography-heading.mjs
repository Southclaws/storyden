import { memo, splitProps } from '../helpers.mjs';
import { createRecipe, mergeRecipes } from './create-recipe.mjs';

const typographyHeadingFn = /* @__PURE__ */ createRecipe('typography-heading', {
  "size": "md"
}, [])

const typographyHeadingVariantMap = {
  "size": [
    "xs",
    "sm",
    "md",
    "lg",
    "xl",
    "2xl"
  ]
}

const typographyHeadingVariantKeys = Object.keys(typographyHeadingVariantMap)

export const typographyHeading = /* @__PURE__ */ Object.assign(memo(typographyHeadingFn.recipeFn), {
  __recipe__: true,
  __name__: 'typographyHeading',
  __getCompoundVariantCss__: typographyHeadingFn.__getCompoundVariantCss__,
  raw: (props) => props,
  variantKeys: typographyHeadingVariantKeys,
  variantMap: typographyHeadingVariantMap,
  merge(recipe) {
    return mergeRecipes(this, recipe)
  },
  splitVariantProps(props) {
    return splitProps(props, typographyHeadingVariantKeys)
  },
  getVariantProps: typographyHeadingFn.getVariantProps,
})