import { createRecipe } from './runtime';

const buttonConfig = {"name":"button","defaultVariants":{"size":"sm","variant":"subtle"},"variantMap":{"intent":["destructive","success","warning"],"size":["lg","md","sm"],"variant":["ghost","outline","plain","solid","subtle"]}}

export const button = /* @__PURE__ */ createRecipe(buttonConfig)