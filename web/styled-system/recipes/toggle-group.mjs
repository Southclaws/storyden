import { compact, getSlotCompoundVariant, memo, splitProps } from '../helpers.mjs';
import { createRecipe } from './create-recipe.mjs';

const toggleGroupDefaultVariants = {
  "size": "md",
  "variant": "outline"
}
const toggleGroupCompoundVariants = [
  {
    "variant": "outline",
    "size": "xs",
    "css": {
      "item": {
        "h": "{sizes.5.5}"
      }
    }
  },
  {
    "variant": "outline",
    "size": "sm",
    "css": {
      "item": {
        "h": "7.5"
      }
    }
  },
  {
    "variant": "outline",
    "size": "md",
    "css": {
      "item": {
        "h": "9.5"
      }
    }
  },
  {
    "variant": "outline",
    "size": "lg",
    "css": {
      "item": {
        "h": "10.5"
      }
    }
  }
]

const toggleGroupSlotNames = [
  [
    "root",
    "toggleGroup__root"
  ],
  [
    "item",
    "toggleGroup__item"
  ]
]
const toggleGroupSlotFns = /* @__PURE__ */ toggleGroupSlotNames.map(([slotName, slotKey]) => [slotName, createRecipe(slotKey, toggleGroupDefaultVariants, getSlotCompoundVariant(toggleGroupCompoundVariants, slotName))])

const toggleGroupFn = memo((props = {}) => {
  return Object.fromEntries(toggleGroupSlotFns.map(([slotName, slotFn]) => [slotName, slotFn.recipeFn(props)]))
})

const toggleGroupVariantKeys = [
  "variant",
  "size"
]
const getVariantProps = (variants) => ({ ...toggleGroupDefaultVariants, ...compact(variants) })

export const toggleGroup = /* @__PURE__ */ Object.assign(toggleGroupFn, {
  __recipe__: false,
  __name__: 'toggleGroup',
  raw: (props) => props,
  classNameMap: {},
  variantKeys: toggleGroupVariantKeys,
  variantMap: {
  "variant": [
    "outline",
    "ghost"
  ],
  "size": [
    "xs",
    "sm",
    "md",
    "lg"
  ]
},
  splitVariantProps(props) {
    return splitProps(props, toggleGroupVariantKeys)
  },
  getVariantProps
})