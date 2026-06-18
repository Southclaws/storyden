import { createRecipe } from './runtime';

const richCardConfig = {"name":"richCard","className":"rich-card","defaultVariants":{"backgroundColor":"default","shape":"row"},"variantMap":{"backgroundColor":["accent","default","emphasized"],"shape":["box","fill","responsive","row"]}}

export const richCard = /* @__PURE__ */ createRecipe(richCardConfig)