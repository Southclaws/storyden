import { createSlotRecipe } from './runtime';

const checkboxConfig = {"name":"checkbox","slots":["root","label","control","indicator","group"],"defaultVariants":{"size":"sm"},"variantMap":{"size":["lg","md","sm"]}}

export const checkbox = /* @__PURE__ */ createSlotRecipe(checkboxConfig)