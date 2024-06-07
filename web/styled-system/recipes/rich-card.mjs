import { compact, getSlotCompoundVariant, memo, splitProps } from '../helpers.mjs';
import { createRecipe } from './create-recipe.mjs';

const richCardDefaultVariants = {
  "mediaDisplay": "with",
  "shape": "box",
  "size": "default"
}
const richCardCompoundVariants = [
  {
    "size": "small",
    "shape": "row",
    "css": {
      "root": {
        "gridTemplateColumns": "1fr 2fr minmax(0, min-content)"
      },
      "text": {
        "display": "none"
      },
      "title": {
        "fontSize": "sm"
      },
      "controlsOverlayContainer": {
        "display": "flex",
        "justifyContent": "end",
        "alignItems": "start",
        "padding": "2"
      }
    }
  }
]

const richCardSlotNames = [
  [
    "root",
    "rich-card__root"
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
    "contentContainer",
    "rich-card__contentContainer"
  ],
  [
    "mediaContainer",
    "rich-card__mediaContainer"
  ],
  [
    "textArea",
    "rich-card__textArea"
  ],
  [
    "footer",
    "rich-card__footer"
  ],
  [
    "title",
    "rich-card__title"
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
  ],
  [
    "controlsOverlayContainer",
    "rich-card__controlsOverlayContainer"
  ],
  [
    "controls",
    "rich-card__controls"
  ]
]
const richCardSlotFns = /* @__PURE__ */ richCardSlotNames.map(([slotName, slotKey]) => [slotName, createRecipe(slotKey, richCardDefaultVariants, getSlotCompoundVariant(richCardCompoundVariants, slotName))])

const richCardFn = memo((props = {}) => {
  return Object.fromEntries(richCardSlotFns.map(([slotName, slotFn]) => [slotName, slotFn.recipeFn(props)]))
})

const richCardVariantKeys = [
  "mediaDisplay",
  "shape",
  "size"
]
const getVariantProps = (variants) => ({ ...richCardDefaultVariants, ...compact(variants) })

export const richCard = /* @__PURE__ */ Object.assign(richCardFn, {
  __recipe__: false,
  __name__: 'richCard',
  raw: (props) => props,
  variantKeys: richCardVariantKeys,
  variantMap: {
  "mediaDisplay": [
    "with",
    "without"
  ],
  "shape": [
    "box",
    "row"
  ],
  "size": [
    "default",
    "small"
  ]
},
  splitVariantProps(props) {
    return splitProps(props, richCardVariantKeys)
  },
  getVariantProps
})