import { compact, getSlotCompoundVariant, memo, splitProps } from '../helpers.mjs';
import { createRecipe } from './create-recipe.mjs';

const richCardDefaultVariants = {
  "shape": "row",
  "backgroundColor": "default"
}
const richCardCompoundVariants = []

const richCardSlotNames = [
  [
    "container",
    "rich-card__container"
  ],
  [
    "root",
    "rich-card__root"
  ],
  [
    "headerContainer",
    "rich-card__headerContainer"
  ],
  [
    "menuContainer",
    "rich-card__menuContainer"
  ],
  [
    "titleContainer",
    "rich-card__titleContainer"
  ],
  [
    "title",
    "rich-card__title"
  ],
  [
    "contentContainer",
    "rich-card__contentContainer"
  ],
  [
    "mediaContainer",
    "rich-card__mediaContainer"
  ],
  [
    "footerContainer",
    "rich-card__footerContainer"
  ],
  [
    "mediaBackdropContainer",
    "rich-card__mediaBackdropContainer"
  ],
  [
    "mediaBackdrop",
    "rich-card__mediaBackdrop"
  ],
  [
    "textArea",
    "rich-card__textArea"
  ],
  [
    "text",
    "rich-card__text"
  ],
  [
    "media",
    "rich-card__media"
  ],
  [
    "mediaMissing",
    "rich-card__mediaMissing"
  ]
]
const richCardSlotFns = /* @__PURE__ */ richCardSlotNames.map(([slotName, slotKey]) => [slotName, createRecipe(slotKey, richCardDefaultVariants, getSlotCompoundVariant(richCardCompoundVariants, slotName))])

const richCardFn = memo((props = {}) => {
  return Object.fromEntries(richCardSlotFns.map(([slotName, slotFn]) => [slotName, slotFn.recipeFn(props)]))
})

const richCardVariantKeys = [
  "backgroundColor",
  "shape"
]
const getVariantProps = (variants) => ({ ...richCardDefaultVariants, ...compact(variants) })

export const richCard = /* @__PURE__ */ Object.assign(richCardFn, {
  __recipe__: false,
  __name__: 'richCard',
  raw: (props) => props,
  classNameMap: {},
  variantKeys: richCardVariantKeys,
  variantMap: {
  "backgroundColor": [
    "default",
    "emphasized",
    "accent"
  ],
  "shape": [
    "row",
    "responsive",
    "box",
    "fill"
  ]
},
  splitVariantProps(props) {
    return splitProps(props, richCardVariantKeys)
  },
  getVariantProps
})