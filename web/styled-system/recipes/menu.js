import { createSlotRecipe } from './runtime';

const menuConfig = {"name":"menu","slots":["arrow","arrowTip","content","contextTrigger","indicator","item","itemGroup","itemGroupLabel","itemIndicator","itemText","positioner","separator","trigger","triggerItem"],"defaultVariants":{"size":"xs"},"variantMap":{"size":["lg","md","sm","xs"]}}

export const menu = /* @__PURE__ */ createSlotRecipe(menuConfig)