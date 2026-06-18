import { createSlotRecipe } from './runtime';

const radioGroupConfig = {"name":"radioGroup","slots":["root","label","item","itemText","itemControl","indicator"],"defaultVariants":{"size":"md"},"variantMap":{"size":["lg","md","sm"]}}

export const radioGroup = /* @__PURE__ */ createSlotRecipe(radioGroupConfig)