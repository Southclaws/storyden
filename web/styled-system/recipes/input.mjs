import { memo, splitProps } from '../helpers.mjs';
import { createRecipe, mergeRecipes } from './create-recipe.mjs';

const inputFn = /* @__PURE__ */ createRecipe('input', {
  "size": "sm",
  "variant": "outline"
}, [
  {
    "size": "sm",
    "variant": "outline",
    "css": {
      "px": "2"
    }
  },
  {
    "size": "md",
    "variant": "outline",
    "css": {
      "px": "2.5"
    }
  },
  {
    "size": "lg",
    "variant": "outline",
    "css": {
      "px": "3"
    }
  }
])

const inputVariantMap = {
  "size": [
    "sm",
    "md",
    "lg"
  ],
  "variant": [
    "outline",
    "ghost",
    "inset"
  ]
}

const inputVariantKeys = Object.keys(inputVariantMap)

export const input = /* @__PURE__ */ Object.assign(memo(inputFn.recipeFn), {
  __recipe__: true,
  __name__: 'input',
  __getCompoundVariantCss__: inputFn.__getCompoundVariantCss__,
  raw: (props) => props,
  variantKeys: inputVariantKeys,
  variantMap: inputVariantMap,
  merge(recipe) {
    return mergeRecipes(this, recipe)
  },
  splitVariantProps(props) {
    return splitProps(props, inputVariantKeys)
  },
  getVariantProps: inputFn.getVariantProps,
})