import { memo, splitProps } from '../helpers.mjs';
import { createRecipe, mergeRecipes } from './create-recipe.mjs';

const multiSelectPickerFn = /* @__PURE__ */ createRecipe('multi-select-picker', {
  "size": "sm"
}, [])

const multiSelectPickerVariantMap = {
  "size": [
    "sm",
    "md",
    "lg"
  ]
}

const multiSelectPickerVariantKeys = Object.keys(multiSelectPickerVariantMap)

export const multiSelectPicker = /* @__PURE__ */ Object.assign(memo(multiSelectPickerFn.recipeFn), {
  __recipe__: true,
  __name__: 'multiSelectPicker',
  __getCompoundVariantCss__: multiSelectPickerFn.__getCompoundVariantCss__,
  raw: (props) => props,
  variantKeys: multiSelectPickerVariantKeys,
  variantMap: multiSelectPickerVariantMap,
  merge(recipe) {
    return mergeRecipes(this, recipe)
  },
  splitVariantProps(props) {
    return splitProps(props, multiSelectPickerVariantKeys)
  },
  getVariantProps: multiSelectPickerFn.getVariantProps,
})