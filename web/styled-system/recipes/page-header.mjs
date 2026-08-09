import { compact, getSlotCompoundVariant, memo, splitProps } from '../helpers.mjs';
import { createRecipe } from './create-recipe.mjs';

const pageHeaderDefaultVariants = {}
const pageHeaderCompoundVariants = []

const pageHeaderSlotNames = [
  [
    "root",
    "page-header__root"
  ],
  [
    "navigation",
    "page-header__navigation"
  ],
  [
    "row",
    "page-header__row"
  ],
  [
    "heading",
    "page-header__heading"
  ],
  [
    "titleGroup",
    "page-header__titleGroup"
  ],
  [
    "actions",
    "page-header__actions"
  ]
]
const pageHeaderSlotFns = /* @__PURE__ */ pageHeaderSlotNames.map(([slotName, slotKey]) => [slotName, createRecipe(slotKey, pageHeaderDefaultVariants, getSlotCompoundVariant(pageHeaderCompoundVariants, slotName))])

const pageHeaderFn = memo((props = {}) => {
  return Object.fromEntries(pageHeaderSlotFns.map(([slotName, slotFn]) => [slotName, slotFn.recipeFn(props)]))
})

const pageHeaderVariantKeys = []
const getVariantProps = (variants) => ({ ...pageHeaderDefaultVariants, ...compact(variants) })

export const pageHeader = /* @__PURE__ */ Object.assign(pageHeaderFn, {
  __recipe__: false,
  __name__: 'pageHeader',
  raw: (props) => props,
  classNameMap: {},
  variantKeys: pageHeaderVariantKeys,
  variantMap: {},
  splitVariantProps(props) {
    return splitProps(props, pageHeaderVariantKeys)
  },
  getVariantProps
})