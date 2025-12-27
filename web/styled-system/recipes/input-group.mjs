import { compact, getSlotCompoundVariant, memo, splitProps } from '../helpers.mjs';
import { createRecipe } from './create-recipe.mjs';

const inputGroupDefaultVariants = {
  "size": "md"
}
const inputGroupCompoundVariants = []

const inputGroupSlotNames = [
  [
    "root",
    "input-group__root"
  ],
  [
    "element",
    "input-group__element"
  ]
]
const inputGroupSlotFns = /* @__PURE__ */ inputGroupSlotNames.map(([slotName, slotKey]) => [slotName, createRecipe(slotKey, inputGroupDefaultVariants, getSlotCompoundVariant(inputGroupCompoundVariants, slotName))])

const inputGroupFn = memo((props = {}) => {
  return Object.fromEntries(inputGroupSlotFns.map(([slotName, slotFn]) => [slotName, slotFn.recipeFn(props)]))
})

const inputGroupVariantKeys = [
  "size"
]
const getVariantProps = (variants) => ({ ...inputGroupDefaultVariants, ...compact(variants) })

export const inputGroup = /* @__PURE__ */ Object.assign(inputGroupFn, {
  __recipe__: false,
  __name__: 'inputGroup',
  raw: (props) => props,
  classNameMap: {},
  variantKeys: inputGroupVariantKeys,
  variantMap: {
  "size": [
    "xs",
    "sm",
    "md",
    "lg",
    "xl"
  ]
},
  splitVariantProps(props) {
    return splitProps(props, inputGroupVariantKeys)
  },
  getVariantProps
})