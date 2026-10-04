import { createRecipe } from './runtime';

const textareaConfig = {"name":"textarea","defaultVariants":{"size":"md","variant":"outline"},"compoundVariants":[{"size":"sm","variant":"outline","className":"textarea--compound__size_sm__variant_outline"},{"size":"md","variant":"outline","className":"textarea--compound__size_md__variant_outline"},{"size":"lg","variant":"outline","className":"textarea--compound__size_lg__variant_outline"}],"variantMap":{"size":["lg","md","sm"],"variant":["ghost","inset","outline"]}}

export const textarea = /* @__PURE__ */ createRecipe(textareaConfig)