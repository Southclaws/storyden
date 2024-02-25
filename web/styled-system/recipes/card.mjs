import { getSlotCompoundVariant, memo, splitProps } from '../helpers.mjs';
import { createRecipe } from './create-recipe.mjs';

const cardDefaultVariants = {
  "mediaDisplay": "with",
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
  ]
]
const cardSlotFns = /* @__PURE__ */ cardSlotNames.map(([slotName, slotKey]) => [slotName, createRecipe(slotKey, cardDefaultVariants, getSlotCompoundVariant(cardCompoundVariants, slotName))])

const cardFn = memo((props = {}) => {
  return Object.fromEntries(cardSlotFns.map(([slotName, slotFn]) => [slotName, slotFn(props)]))
})

const cardVariantKeys = [
  "mediaDisplay",
  "shape"
]

export const card = /* @__PURE__ */ Object.assign(cardFn, {
  __recipe__: false,
  __name__: 'card',
  raw: (props) => props,
  variantKeys: cardVariantKeys,
  variantMap: {
  "mediaDisplay": [
    "with",
    "without"
  ],
  "shape": [
    "box",
    "row"
  ]
},
  splitVariantProps(props) {
    return splitProps(props, cardVariantKeys)
  },
})