import { compact, getSlotCompoundVariant, memo, splitProps } from '../helpers.mjs';
import { createRecipe } from './create-recipe.mjs';

const cardGridDefaultVariants = {}
const cardGridCompoundVariants = []

const cardGridSlotNames = [
  [
    "container",
    "card-grid__container"
  ],
  [
    "grid",
    "card-grid__grid"
  ]
]
const cardGridSlotFns = /* @__PURE__ */ cardGridSlotNames.map(([slotName, slotKey]) => [slotName, createRecipe(slotKey, cardGridDefaultVariants, getSlotCompoundVariant(cardGridCompoundVariants, slotName))])

const cardGridFn = memo((props = {}) => {
  return Object.fromEntries(cardGridSlotFns.map(([slotName, slotFn]) => [slotName, slotFn.recipeFn(props)]))
})

const cardGridVariantKeys = []
const getVariantProps = (variants) => ({ ...cardGridDefaultVariants, ...compact(variants) })

export const cardGrid = /* @__PURE__ */ Object.assign(cardGridFn, {
  __recipe__: false,
  __name__: 'cardGrid',
  raw: (props) => props,
  classNameMap: {},
  variantKeys: cardGridVariantKeys,
  variantMap: {},
  splitVariantProps(props) {
    return splitProps(props, cardGridVariantKeys)
  },
  getVariantProps
})