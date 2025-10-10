import { compact, getSlotCompoundVariant, memo, splitProps } from '../helpers.mjs';
import { createRecipe } from './create-recipe.mjs';

const datePickerDefaultVariants = {}
const datePickerCompoundVariants = []

const datePickerSlotNames = [
  [
    "clearTrigger",
    "datePicker__clearTrigger"
  ],
  [
    "content",
    "datePicker__content"
  ],
  [
    "control",
    "datePicker__control"
  ],
  [
    "input",
    "datePicker__input"
  ],
  [
    "label",
    "datePicker__label"
  ],
  [
    "monthSelect",
    "datePicker__monthSelect"
  ],
  [
    "nextTrigger",
    "datePicker__nextTrigger"
  ],
  [
    "positioner",
    "datePicker__positioner"
  ],
  [
    "presetTrigger",
    "datePicker__presetTrigger"
  ],
  [
    "prevTrigger",
    "datePicker__prevTrigger"
  ],
  [
    "rangeText",
    "datePicker__rangeText"
  ],
  [
    "root",
    "datePicker__root"
  ],
  [
    "table",
    "datePicker__table"
  ],
  [
    "tableBody",
    "datePicker__tableBody"
  ],
  [
    "tableCell",
    "datePicker__tableCell"
  ],
  [
    "tableCellTrigger",
    "datePicker__tableCellTrigger"
  ],
  [
    "tableHead",
    "datePicker__tableHead"
  ],
  [
    "tableHeader",
    "datePicker__tableHeader"
  ],
  [
    "tableRow",
    "datePicker__tableRow"
  ],
  [
    "trigger",
    "datePicker__trigger"
  ],
  [
    "view",
    "datePicker__view"
  ],
  [
    "viewControl",
    "datePicker__viewControl"
  ],
  [
    "viewTrigger",
    "datePicker__viewTrigger"
  ],
  [
    "yearSelect",
    "datePicker__yearSelect"
  ],
  [
    "view",
    "datePicker__view"
  ]
]
const datePickerSlotFns = /* @__PURE__ */ datePickerSlotNames.map(([slotName, slotKey]) => [slotName, createRecipe(slotKey, datePickerDefaultVariants, getSlotCompoundVariant(datePickerCompoundVariants, slotName))])

const datePickerFn = memo((props = {}) => {
  return Object.fromEntries(datePickerSlotFns.map(([slotName, slotFn]) => [slotName, slotFn.recipeFn(props)]))
})

const datePickerVariantKeys = []
const getVariantProps = (variants) => ({ ...datePickerDefaultVariants, ...compact(variants) })

export const datePicker = /* @__PURE__ */ Object.assign(datePickerFn, {
  __recipe__: false,
  __name__: 'datePicker',
  raw: (props) => props,
  classNameMap: {},
  variantKeys: datePickerVariantKeys,
  variantMap: {},
  splitVariantProps(props) {
    return splitProps(props, datePickerVariantKeys)
  },
  getVariantProps
})