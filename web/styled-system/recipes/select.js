import { createSlotRecipe } from './runtime';

const selectConfig = {"name":"select","slots":["label","positioner","trigger","indicator","clearTrigger","item","itemText","itemIndicator","itemGroup","itemGroupLabel","list","content","root","control","valueText"],"defaultVariants":{"size":"sm","variant":"outline"},"variantMap":{"variant":["outline","ghost"],"size":["sm","md","lg"]}}

export const select = /* @__PURE__ */ createSlotRecipe(selectConfig)