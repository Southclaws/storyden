import { createSlotRecipe } from './runtime';

const numberInputConfig = {"name":"numberInput","className":"number-input","slots":["root","label","input","control","valueText","incrementTrigger","decrementTrigger","scrubber"],"defaultVariants":{"size":"md","variant":"outline"},"variantMap":{"size":["lg","md","sm","xl"],"variant":["ghost","outline"]}}

export const numberInput = /* @__PURE__ */ createSlotRecipe(numberInputConfig)