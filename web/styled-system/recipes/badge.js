import { createRecipe } from './runtime';

const badgeConfig = {"name":"badge","defaultVariants":{"size":"md","variant":"subtle"},"variantMap":{"variant":["solid","subtle","outline"],"size":["sm","md","lg"]}}

export const badge = /* @__PURE__ */ createRecipe(badgeConfig)