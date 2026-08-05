import { compact, getSlotCompoundVariant, memo, splitProps } from '../helpers.mjs';
import { createRecipe } from './create-recipe.mjs';

const reactListDefaultVariants = {}
const reactListCompoundVariants = []

const reactListSlotNames = [
  [
    "root",
    "reactList__root"
  ],
  [
    "reaction",
    "reactList__reaction"
  ],
  [
    "count",
    "reactList__count"
  ],
  [
    "picker",
    "reactList__picker"
  ]
]
const reactListSlotFns = /* @__PURE__ */ reactListSlotNames.map(([slotName, slotKey]) => [slotName, createRecipe(slotKey, reactListDefaultVariants, getSlotCompoundVariant(reactListCompoundVariants, slotName))])

const reactListFn = memo((props = {}) => {
  return Object.fromEntries(reactListSlotFns.map(([slotName, slotFn]) => [slotName, slotFn.recipeFn(props)]))
})

const reactListVariantKeys = []
const getVariantProps = (variants) => ({ ...reactListDefaultVariants, ...compact(variants) })

export const reactList = /* @__PURE__ */ Object.assign(reactListFn, {
  __recipe__: false,
  __name__: 'reactList',
  raw: (props) => props,
  classNameMap: {},
  variantKeys: reactListVariantKeys,
  variantMap: {},
  splitVariantProps(props) {
    return splitProps(props, reactListVariantKeys)
  },
  getVariantProps
})