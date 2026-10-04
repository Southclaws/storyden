import { createSlotRecipe } from './runtime';

const sliderConfig = {"name":"slider","slots":["root","label","thumb","valueText","track","range","control","markerGroup","marker","draggingIndicator"],"defaultVariants":{"size":"md"},"variantMap":{"size":["lg","md","sm"]}}

export const slider = /* @__PURE__ */ createSlotRecipe(sliderConfig)