import { createSlotRecipe } from './runtime';

const floatingPanelConfig = {"name":"floatingPanel","className":"floating-panel","slots":["trigger","positioner","content","header","body","title","resizeTrigger","dragTrigger","stageTrigger","closeTrigger","control"]}

export const floatingPanel = /* @__PURE__ */ createSlotRecipe(floatingPanelConfig)