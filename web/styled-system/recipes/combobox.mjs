import { compact, getSlotCompoundVariant, memo, splitProps } from '../helpers.mjs';
import { createRecipe } from './create-recipe.mjs';

const comboboxDefaultVariants = {
  "size": "md"
}
const comboboxCompoundVariants = []

const comboboxSlotNames = [
  [
    "root",
    "combobox__root"
  ],
  [
    "clearTrigger",
    "combobox__clearTrigger"
  ],
  [
    "content",
    "combobox__content"
  ],
  [
    "control",
    "combobox__control"
  ],
  [
    "input",
    "combobox__input"
  ],
  [
    "item",
    "combobox__item"
  ],
  [
    "itemGroup",
    "combobox__itemGroup"
  ],
  [
    "itemGroupLabel",
    "combobox__itemGroupLabel"
  ],
  [
    "itemIndicator",
    "combobox__itemIndicator"
  ],
  [
    "itemText",
    "combobox__itemText"
  ],
  [
    "label",
    "combobox__label"
  ],
  [
    "list",
    "combobox__list"
  ],
  [
    "positioner",
    "combobox__positioner"
  ],
  [
    "trigger",
    "combobox__trigger"
  ],
  [
    "root",
    "combobox__root"
  ],
  [
    "clearTrigger",
    "combobox__clearTrigger"
  ],
  [
    "content",
    "combobox__content"
  ],
  [
    "control",
    "combobox__control"
  ],
  [
    "input",
    "combobox__input"
  ],
  [
    "item",
    "combobox__item"
  ],
  [
    "itemGroup",
    "combobox__itemGroup"
  ],
  [
    "itemGroupLabel",
    "combobox__itemGroupLabel"
  ],
  [
    "itemIndicator",
    "combobox__itemIndicator"
  ],
  [
    "itemText",
    "combobox__itemText"
  ],
  [
    "label",
    "combobox__label"
  ],
  [
    "list",
    "combobox__list"
  ],
  [
    "positioner",
    "combobox__positioner"
  ],
  [
    "trigger",
    "combobox__trigger"
  ]
]
const comboboxSlotFns = /* @__PURE__ */ comboboxSlotNames.map(([slotName, slotKey]) => [slotName, createRecipe(slotKey, comboboxDefaultVariants, getSlotCompoundVariant(comboboxCompoundVariants, slotName))])

const comboboxFn = memo((props = {}) => {
  return Object.fromEntries(comboboxSlotFns.map(([slotName, slotFn]) => [slotName, slotFn.recipeFn(props)]))
})

const comboboxVariantKeys = [
  "size"
]
const getVariantProps = (variants) => ({ ...comboboxDefaultVariants, ...compact(variants) })

export const combobox = /* @__PURE__ */ Object.assign(comboboxFn, {
  __recipe__: false,
  __name__: 'combobox',
  raw: (props) => props,
  variantKeys: comboboxVariantKeys,
  variantMap: {
  "size": [
    "sm",
    "md",
    "lg"
  ]
},
  splitVariantProps(props) {
    return splitProps(props, comboboxVariantKeys)
  },
  getVariantProps
})