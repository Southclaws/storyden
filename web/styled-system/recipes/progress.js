import { createSlotRecipe } from './runtime';

const progressConfig = {"name":"progress","slots":["root","label","track","range","valueText","view","circle","circleTrack","circleRange"],"defaultVariants":{"size":"md"},"variantMap":{"size":["lg","md","sm"]}}

export const progress = /* @__PURE__ */ createSlotRecipe(progressConfig)