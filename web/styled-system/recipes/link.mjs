import { memo, splitProps } from '../helpers.mjs';
import { createRecipe, mergeRecipes } from './create-recipe.mjs';

const linkFn = /* @__PURE__ */ createRecipe('link', {
  "kind": "neutral",
  "size": "md"
}, [])

const linkVariantMap = {
  "kind": [
    "neutral",
    "primary",
    "destructive",
    "ghost",
    "blank"
  ],
  "size": [
    "xs",
    "sm",
    "md",
    "lg",
    "xl",
    "2xl"
  ]
}

const linkVariantKeys = Object.keys(linkVariantMap)

export const link = /* @__PURE__ */ Object.assign(memo(linkFn.recipeFn), {
  __recipe__: true,
  __name__: 'link',
  __getCompoundVariantCss__: linkFn.__getCompoundVariantCss__,
  raw: (props) => props,
  variantKeys: linkVariantKeys,
  variantMap: linkVariantMap,
  merge(recipe) {
    return mergeRecipes(this, recipe)
  },
  splitVariantProps(props) {
    return splitProps(props, linkVariantKeys)
  },
  getVariantProps: linkFn.getVariantProps,
})