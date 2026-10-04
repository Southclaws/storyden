import { createRecipe } from './runtime';

const badgeConfig = {"name":"badge","defaultVariants":{"size":"md","variant":"subtle"},"variantMap":{"size":["lg","md","sm"],"variant":["outline","solid","subtle"]}}

export const badge = /* @__PURE__ */ createRecipe(badgeConfig)