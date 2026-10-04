import { createSlotRecipe } from './runtime';

const blockEditorConfig = {"name":"blockEditor","className":"block-editor","slots":["root","gutter","handle","menuTrigger","menuAnchor","menuContent","content"]}

export const blockEditor = /* @__PURE__ */ createSlotRecipe(blockEditorConfig)