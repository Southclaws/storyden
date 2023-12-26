import { splitProps } from '../helpers.mjs';
import { createRecipe } from './create-recipe.mjs';

const titleInputFn = /* @__PURE__ */ createRecipe('titleInput', {}, [])

const titleInputVariantMap = {}

const titleInputVariantKeys = Object.keys(titleInputVariantMap)

export const titleInput = /* @__PURE__ */ Object.assign(titleInputFn, {
  __recipe__: true,
  __name__: 'titleInput',
  raw: (props) => props,
  variantKeys: titleInputVariantKeys,
  variantMap: titleInputVariantMap,
  splitVariantProps(props) {
    return splitProps(props, titleInputVariantKeys)
  },
})