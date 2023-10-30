import { splitProps } from '../helpers.mjs';
import { createRecipe } from './create-recipe.mjs';

const inputFn = /* @__PURE__ */ createRecipe('input', {
  "size": "md"
}, [])

const inputVariantMap = {
  "size": [
    "2xs",
    "xs",
    "sm",
    "md",
    "lg",
    "xl",
    "2xl"
  ]
}

const inputVariantKeys = Object.keys(inputVariantMap)

export const input = /* @__PURE__ */ Object.assign(inputFn, {
  __recipe__: true,
  __name__: 'input',
  raw: (props) => props,
  variantKeys: inputVariantKeys,
  variantMap: inputVariantMap,
  splitVariantProps(props) {
    return splitProps(props, inputVariantKeys)
  },
})