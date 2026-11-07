import { createSlotRecipe } from './runtime';

const progressConfig = {"name":"progress","slots":["root","label","track","range","valueText","view","circle","circleTrack","circleRange"],"defaultVariants":{"shape":"horizontal","size":"md"},"variantMap":{"shape":["circle","horizontal"],"size":["sm","md","lg"]}}

export const progress = /* @__PURE__ */ createSlotRecipe(progressConfig)