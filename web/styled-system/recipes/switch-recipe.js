import { createSlotRecipe } from './runtime';

const switchRecipeConfig = {"name":"switchRecipe","className":"switch","slots":["root","label","control","thumb"],"defaultVariants":{"size":"md"},"variantMap":{"size":["sm","md","lg"]}}

export const switchRecipe = /* @__PURE__ */ createSlotRecipe(switchRecipeConfig)