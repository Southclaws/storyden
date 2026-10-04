import { createSlotRecipe } from './runtime';

const treeViewConfig = {"name":"treeView","slots":["branch","branchContent","branchControl","branchIndentGuide","branchIndicator","branchText","branchTrigger","item","itemIndicator","itemText","label","nodeCheckbox","nodeRenameInput","root","tree"],"defaultVariants":{"variant":"clamped"},"variantMap":{"variant":["clamped","scrollable"]}}

export const treeView = /* @__PURE__ */ createSlotRecipe(treeViewConfig)