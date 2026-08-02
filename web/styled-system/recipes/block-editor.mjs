import { compact, getSlotCompoundVariant, memo, splitProps } from '../helpers.mjs';
import { createRecipe } from './create-recipe.mjs';

const blockEditorDefaultVariants = {}
const blockEditorCompoundVariants = []

const blockEditorSlotNames = [
  [
    "root",
    "block-editor__root"
  ],
  [
    "gutter",
    "block-editor__gutter"
  ],
  [
    "handle",
    "block-editor__handle"
  ],
  [
    "content",
    "block-editor__content"
  ]
]
const blockEditorSlotFns = /* @__PURE__ */ blockEditorSlotNames.map(([slotName, slotKey]) => [slotName, createRecipe(slotKey, blockEditorDefaultVariants, getSlotCompoundVariant(blockEditorCompoundVariants, slotName))])

const blockEditorFn = memo((props = {}) => {
  return Object.fromEntries(blockEditorSlotFns.map(([slotName, slotFn]) => [slotName, slotFn.recipeFn(props)]))
})

const blockEditorVariantKeys = []
const getVariantProps = (variants) => ({ ...blockEditorDefaultVariants, ...compact(variants) })

export const blockEditor = /* @__PURE__ */ Object.assign(blockEditorFn, {
  __recipe__: false,
  __name__: 'blockEditor',
  raw: (props) => props,
  classNameMap: {},
  variantKeys: blockEditorVariantKeys,
  variantMap: {},
  splitVariantProps(props) {
    return splitProps(props, blockEditorVariantKeys)
  },
  getVariantProps
})