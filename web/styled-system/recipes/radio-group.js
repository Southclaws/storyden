import { createSlotRecipe } from './runtime';

const radioGroupConfig = {"name":"radioGroup","slots":["root","label","item","itemText","itemControl","indicator"],"defaultVariants":{"size":"sm"},"variantMap":{"size":["sm","md","lg"]}}

export const radioGroup = /* @__PURE__ */ createSlotRecipe(radioGroupConfig)