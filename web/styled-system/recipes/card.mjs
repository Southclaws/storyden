import { splitProps, getSlotCompoundVariant } from '../helpers.mjs';
import { createRecipe } from './create-recipe.mjs';

const cardDefaultVariants = {
  "shape": "box"
}
const cardCompoundVariants = []

const cardSlotNames = [
  [
    "root",
    "card__root"
  ],
  [
    "textArea",
    "card__textArea"
  ],
  [
    "title",
    "card__title"
  ],
  [
    "text",
    "card__text"
  ],
  [
    "media",
    "card__media"
  ],
  [
    "mediaContainer",
    "card__mediaContainer"
  ],
  [
    "mediaBackdrop",
    "card__mediaBackdrop"
  ],
  [
    "mediaBackdropContainer",
    "card__mediaBackdropContainer"
  ]
]
const cardSlotFns = /* @__PURE__ */ cardSlotNames.map(([slotName, slotKey]) => [slotName, createRecipe(slotKey, cardDefaultVariants, getSlotCompoundVariant(cardCompoundVariants, slotName))])

const cardFn = (props = {}) => {
  return Object.fromEntries(cardSlotFns.map(([slotName, slotFn]) => [slotName, slotFn(props)]))
}

const cardVariantKeys = [
  "shape"
]

export const card = /* @__PURE__ */ Object.assign(cardFn, {
  __recipe__: false,
  __name__: 'card',
  raw: (props) => props,
  variantKeys: cardVariantKeys,
  variantMap: {
  "shape": [
    "box",
    "row"
  ]
},
  splitVariantProps(props) {
    return splitProps(props, cardVariantKeys)
  },
})