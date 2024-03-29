import { memo, splitProps } from '../helpers.mjs';
import { createRecipe, mergeRecipes } from './create-recipe.mjs';

const buttonFn = /* @__PURE__ */ createRecipe('button', {
  "kind": "neutral",
  "size": "md"
}, [
  {
    "kind": "blank",
    "css": {
      "px": "0",
      "color": "accent.text.50"
    }
  }
])

const buttonVariantMap = {
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

const buttonVariantKeys = Object.keys(buttonVariantMap)

export const button = /* @__PURE__ */ Object.assign(memo(buttonFn), {
  __recipe__: true,
  __name__: 'button',
  raw: (props) => props,
  variantKeys: buttonVariantKeys,
  variantMap: buttonVariantMap,
  merge(recipe) {
    return mergeRecipes(this, recipe)
  },
  splitVariantProps(props) {
    return splitProps(props, buttonVariantKeys)
  },
})