import { createRecipe } from './runtime';

const switchRecipeConfig = {"name":"switchRecipe","defaultVariants":{"size":"md"},"variantMap":{"size":["lg","md","sm"]}}

export const switchRecipe = /* @__PURE__ */ createRecipe(switchRecipeConfig)