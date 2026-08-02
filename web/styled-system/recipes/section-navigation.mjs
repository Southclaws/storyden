import { compact, getSlotCompoundVariant, memo, splitProps } from '../helpers.mjs';
import { createRecipe } from './create-recipe.mjs';

const sectionNavigationDefaultVariants = {}
const sectionNavigationCompoundVariants = []

const sectionNavigationSlotNames = [
  [
    "root",
    "section-navigation__root"
  ],
  [
    "trigger",
    "section-navigation__trigger"
  ],
  [
    "triggerLabel",
    "section-navigation__triggerLabel"
  ],
  [
    "triggerIndicator",
    "section-navigation__triggerIndicator"
  ],
  [
    "positioner",
    "section-navigation__positioner"
  ],
  [
    "content",
    "section-navigation__content"
  ],
  [
    "navigation",
    "section-navigation__navigation"
  ],
  [
    "section",
    "section-navigation__section"
  ],
  [
    "sectionLabel",
    "section-navigation__sectionLabel"
  ],
  [
    "items",
    "section-navigation__items"
  ],
  [
    "item",
    "section-navigation__item"
  ],
  [
    "link",
    "section-navigation__link"
  ],
  [
    "linkIndicator",
    "section-navigation__linkIndicator"
  ]
]
const sectionNavigationSlotFns = /* @__PURE__ */ sectionNavigationSlotNames.map(([slotName, slotKey]) => [slotName, createRecipe(slotKey, sectionNavigationDefaultVariants, getSlotCompoundVariant(sectionNavigationCompoundVariants, slotName))])

const sectionNavigationFn = memo((props = {}) => {
  return Object.fromEntries(sectionNavigationSlotFns.map(([slotName, slotFn]) => [slotName, slotFn.recipeFn(props)]))
})

const sectionNavigationVariantKeys = []
const getVariantProps = (variants) => ({ ...sectionNavigationDefaultVariants, ...compact(variants) })

export const sectionNavigation = /* @__PURE__ */ Object.assign(sectionNavigationFn, {
  __recipe__: false,
  __name__: 'sectionNavigation',
  raw: (props) => props,
  classNameMap: {},
  variantKeys: sectionNavigationVariantKeys,
  variantMap: {},
  splitVariantProps(props) {
    return splitProps(props, sectionNavigationVariantKeys)
  },
  getVariantProps
})