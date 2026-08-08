import { compact, getSlotCompoundVariant, memo, splitProps } from '../helpers.mjs';
import { createRecipe } from './create-recipe.mjs';

const dragTreeDefaultVariants = {
  "variant": "clamped"
}
const dragTreeCompoundVariants = []

const dragTreeSlotNames = [
  [
    "branch",
    "dragTree__branch"
  ],
  [
    "branchContent",
    "dragTree__branchContent"
  ],
  [
    "branchControl",
    "dragTree__branchControl"
  ],
  [
    "branchIndentGuide",
    "dragTree__branchIndentGuide"
  ],
  [
    "branchIndicator",
    "dragTree__branchIndicator"
  ],
  [
    "branchText",
    "dragTree__branchText"
  ],
  [
    "branchTrigger",
    "dragTree__branchTrigger"
  ],
  [
    "item",
    "dragTree__item"
  ],
  [
    "itemIndicator",
    "dragTree__itemIndicator"
  ],
  [
    "itemText",
    "dragTree__itemText"
  ],
  [
    "label",
    "dragTree__label"
  ],
  [
    "nodeCheckbox",
    "dragTree__nodeCheckbox"
  ],
  [
    "nodeRenameInput",
    "dragTree__nodeRenameInput"
  ],
  [
    "root",
    "dragTree__root"
  ],
  [
    "tree",
    "dragTree__tree"
  ],
  [
    "actions",
    "dragTree__actions"
  ],
  [
    "link",
    "dragTree__link"
  ],
  [
    "row",
    "dragTree__row"
  ],
  [
    "dropIndicator",
    "dragTree__dropIndicator"
  ],
  [
    "dropIndicatorLine",
    "dragTree__dropIndicatorLine"
  ]
]
const dragTreeSlotFns = /* @__PURE__ */ dragTreeSlotNames.map(([slotName, slotKey]) => [slotName, createRecipe(slotKey, dragTreeDefaultVariants, getSlotCompoundVariant(dragTreeCompoundVariants, slotName))])

const dragTreeFn = memo((props = {}) => {
  return Object.fromEntries(dragTreeSlotFns.map(([slotName, slotFn]) => [slotName, slotFn.recipeFn(props)]))
})

const dragTreeVariantKeys = [
  "variant"
]
const getVariantProps = (variants) => ({ ...dragTreeDefaultVariants, ...compact(variants) })

export const dragTree = /* @__PURE__ */ Object.assign(dragTreeFn, {
  __recipe__: false,
  __name__: 'dragTree',
  raw: (props) => props,
  classNameMap: {},
  variantKeys: dragTreeVariantKeys,
  variantMap: {
  "variant": [
    "clamped",
    "scrollable"
  ]
},
  splitVariantProps(props) {
    return splitProps(props, dragTreeVariantKeys)
  },
  getVariantProps
})