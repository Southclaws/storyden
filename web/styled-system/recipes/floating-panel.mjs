import { compact, getSlotCompoundVariant, memo, splitProps } from '../helpers.mjs';
import { createRecipe } from './create-recipe.mjs';

const floatingPanelDefaultVariants = {}
const floatingPanelCompoundVariants = []

const floatingPanelSlotNames = [
  [
    "trigger",
    "floating-panel__trigger"
  ],
  [
    "positioner",
    "floating-panel__positioner"
  ],
  [
    "content",
    "floating-panel__content"
  ],
  [
    "header",
    "floating-panel__header"
  ],
  [
    "body",
    "floating-panel__body"
  ],
  [
    "title",
    "floating-panel__title"
  ],
  [
    "resizeTrigger",
    "floating-panel__resizeTrigger"
  ],
  [
    "dragTrigger",
    "floating-panel__dragTrigger"
  ],
  [
    "stageTrigger",
    "floating-panel__stageTrigger"
  ],
  [
    "closeTrigger",
    "floating-panel__closeTrigger"
  ],
  [
    "control",
    "floating-panel__control"
  ]
]
const floatingPanelSlotFns = /* @__PURE__ */ floatingPanelSlotNames.map(([slotName, slotKey]) => [slotName, createRecipe(slotKey, floatingPanelDefaultVariants, getSlotCompoundVariant(floatingPanelCompoundVariants, slotName))])

const floatingPanelFn = memo((props = {}) => {
  return Object.fromEntries(floatingPanelSlotFns.map(([slotName, slotFn]) => [slotName, slotFn.recipeFn(props)]))
})

const floatingPanelVariantKeys = []
const getVariantProps = (variants) => ({ ...floatingPanelDefaultVariants, ...compact(variants) })

export const floatingPanel = /* @__PURE__ */ Object.assign(floatingPanelFn, {
  __recipe__: false,
  __name__: 'floatingPanel',
  raw: (props) => props,
  classNameMap: {},
  variantKeys: floatingPanelVariantKeys,
  variantMap: {},
  splitVariantProps(props) {
    return splitProps(props, floatingPanelVariantKeys)
  },
  getVariantProps
})