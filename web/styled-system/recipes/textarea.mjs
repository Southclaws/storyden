import { memo, splitProps } from '../helpers.mjs';
import { createRecipe, mergeRecipes } from './create-recipe.mjs';

const textareaFn = /* @__PURE__ */ createRecipe('textarea', {
  "size": "md",
  "variant": "outline"
}, [
  {
    "size": "sm",
    "variant": "outline",
    "css": {
      "px": "2",
      "py": "1.5"
    }
  },
  {
    "size": "md",
    "variant": "outline",
    "css": {
      "px": "2.5",
      "py": "2"
    }
  },
  {
    "size": "lg",
    "variant": "outline",
    "css": {
      "px": "3",
      "py": "2.5"
    }
  }
])

const textareaVariantMap = {
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

const textareaVariantKeys = Object.keys(textareaVariantMap)

export const textarea = /* @__PURE__ */ Object.assign(memo(textareaFn.recipeFn), {
  __recipe__: true,
  __name__: 'textarea',
  __getCompoundVariantCss__: textareaFn.__getCompoundVariantCss__,
  raw: (props) => props,
  variantKeys: textareaVariantKeys,
  variantMap: textareaVariantMap,
  merge(recipe) {
    return mergeRecipes(this, recipe)
  },
  splitVariantProps(props) {
    return splitProps(props, textareaVariantKeys)
  },
  getVariantProps: textareaFn.getVariantProps,
})