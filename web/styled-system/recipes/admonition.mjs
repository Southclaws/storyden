import { splitProps } from '../helpers.mjs';
import { createRecipe } from './create-recipe.mjs';

const admonitionFn = /* @__PURE__ */ createRecipe('admonition', {}, [])

const admonitionVariantMap = {
  "kind": [
    "neutral",
    "success",
    "failure"
  ]
}

const admonitionVariantKeys = Object.keys(admonitionVariantMap)

export const admonition = /* @__PURE__ */ Object.assign(admonitionFn, {
  __recipe__: true,
  __name__: 'admonition',
  raw: (props) => props,
  variantKeys: admonitionVariantKeys,
  variantMap: admonitionVariantMap,
  splitVariantProps(props) {
    return splitProps(props, admonitionVariantKeys)
  },
})