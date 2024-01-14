import { splitProps } from '../helpers.mjs';
import { createRecipe } from './create-recipe.mjs';

const headingInputFn = /* @__PURE__ */ createRecipe('headingInput', {}, [])

const headingInputVariantMap = {}

const headingInputVariantKeys = Object.keys(headingInputVariantMap)

export const headingInput = /* @__PURE__ */ Object.assign(headingInputFn, {
  __recipe__: true,
  __name__: 'headingInput',
  raw: (props) => props,
  variantKeys: headingInputVariantKeys,
  variantMap: headingInputVariantMap,
  splitVariantProps(props) {
    return splitProps(props, headingInputVariantKeys)
  },
})