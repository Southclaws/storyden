import { compact, getSlotCompoundVariant, memo, splitProps } from '../helpers.mjs';
import { createRecipe } from './create-recipe.mjs';

const treeViewDefaultVariants = {
  "variant": "clamped"
}
const treeViewCompoundVariants = []

const treeViewSlotNames = [
  [
    "branch",
    "treeView__branch"
  ],
  [
    "branchContent",
    "treeView__branchContent"
  ],
  [
    "branchControl",
    "treeView__branchControl"
  ],
  [
    "branchIndentGuide",
    "treeView__branchIndentGuide"
  ],
  [
    "branchIndicator",
    "treeView__branchIndicator"
  ],
  [
    "branchText",
    "treeView__branchText"
  ],
  [
    "branchTrigger",
    "treeView__branchTrigger"
  ],
  [
    "item",
    "treeView__item"
  ],
  [
    "itemIndicator",
    "treeView__itemIndicator"
  ],
  [
    "itemText",
    "treeView__itemText"
  ],
  [
    "label",
    "treeView__label"
  ],
  [
    "nodeCheckbox",
    "treeView__nodeCheckbox"
  ],
  [
    "nodeRenameInput",
    "treeView__nodeRenameInput"
  ],
  [
    "root",
    "treeView__root"
  ],
  [
    "tree",
    "treeView__tree"
  ]
]
const treeViewSlotFns = /* @__PURE__ */ treeViewSlotNames.map(([slotName, slotKey]) => [slotName, createRecipe(slotKey, treeViewDefaultVariants, getSlotCompoundVariant(treeViewCompoundVariants, slotName))])

const treeViewFn = memo((props = {}) => {
  return Object.fromEntries(treeViewSlotFns.map(([slotName, slotFn]) => [slotName, slotFn.recipeFn(props)]))
})

const treeViewVariantKeys = [
  "variant"
]
const getVariantProps = (variants) => ({ ...treeViewDefaultVariants, ...compact(variants) })

export const treeView = /* @__PURE__ */ Object.assign(treeViewFn, {
  __recipe__: false,
  __name__: 'treeView',
  raw: (props) => props,
  classNameMap: {},
  variantKeys: treeViewVariantKeys,
  variantMap: {
  "variant": [
    "clamped",
    "scrollable"
  ]
},
  splitVariantProps(props) {
    return splitProps(props, treeViewVariantKeys)
  },
  getVariantProps
})