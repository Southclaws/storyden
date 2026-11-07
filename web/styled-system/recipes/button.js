import { createRecipe } from './runtime';

const buttonConfig = {"name":"button","defaultVariants":{"size":"sm","variant":"subtle"},"variantMap":{"size":["sm","md","lg"],"variant":["solid","outline","ghost","subtle","plain"],"intent":["success","warning","destructive"]}}

export const button = /* @__PURE__ */ createRecipe(buttonConfig)