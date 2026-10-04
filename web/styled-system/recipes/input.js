import { createRecipe } from './runtime';

const inputConfig = {"name":"input","defaultVariants":{"size":"sm","variant":"outline"},"compoundVariants":[{"size":"sm","variant":"outline","className":"input--compound__size_sm__variant_outline"},{"size":"md","variant":"outline","className":"input--compound__size_md__variant_outline"},{"size":"lg","variant":"outline","className":"input--compound__size_lg__variant_outline"}],"variantMap":{"size":["lg","md","sm"],"variant":["ghost","inset","outline"]}}

export const input = /* @__PURE__ */ createRecipe(inputConfig)