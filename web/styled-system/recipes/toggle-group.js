import { createSlotRecipe } from './runtime';

const toggleGroupConfig = {"name":"toggleGroup","slots":["root","item"],"defaultVariants":{"size":"md","variant":"outline"},"compoundVariants":[{"size":"xs","variant":"outline","classNames":{"item":"toggleGroup__item--compound__size_xs__variant_outline"}},{"size":"sm","variant":"outline","classNames":{"item":"toggleGroup__item--compound__size_sm__variant_outline"}},{"size":"md","variant":"outline","classNames":{"item":"toggleGroup__item--compound__size_md__variant_outline"}},{"size":"lg","variant":"outline","classNames":{"item":"toggleGroup__item--compound__size_lg__variant_outline"}}],"variantMap":{"size":["lg","md","sm","xs"],"variant":["ghost","outline"]}}

export const toggleGroup = /* @__PURE__ */ createSlotRecipe(toggleGroupConfig)