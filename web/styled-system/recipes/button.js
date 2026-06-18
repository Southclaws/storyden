import { createRecipe } from './runtime';

const buttonConfig = {"name":"button","defaultVariants":{"size":"md","variant":"solid"},"variantMap":{"size":["2xl","lg","md","sm","xl","xs"],"variant":["ghost","link","outline","solid","subtle"]}}

export const button = /* @__PURE__ */ createRecipe(buttonConfig)