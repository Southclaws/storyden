import { createSlotRecipe } from './runtime';

const dragTreeConfig = {"name":"dragTree","slots":["branch","branchContent","branchControl","branchIndentGuide","branchIndicator","branchText","branchTrigger","item","itemIndicator","itemText","label","nodeCheckbox","nodeRenameInput","root","tree","actions","link","row","dropIndicator","dropIndicatorLine"],"defaultVariants":{"variant":"clamped"},"variantMap":{"variant":["clamped","scrollable"]}}

export const dragTree = /* @__PURE__ */ createSlotRecipe(dragTreeConfig)