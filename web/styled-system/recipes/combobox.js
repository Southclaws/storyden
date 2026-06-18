import { createSlotRecipe } from './runtime';

const comboboxConfig = {"name":"combobox","slots":["root","clearTrigger","content","control","input","item","itemGroup","itemGroupLabel","itemIndicator","itemText","label","list","positioner","trigger","empty"],"defaultVariants":{"size":"md"},"variantMap":{"size":["lg","md","sm","xs"]}}

export const combobox = /* @__PURE__ */ createSlotRecipe(comboboxConfig)