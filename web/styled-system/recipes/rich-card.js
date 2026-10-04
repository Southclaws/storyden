import { createSlotRecipe } from './runtime';

const richCardConfig = {"name":"richCard","className":"rich-card","slots":["container","root","headerContainer","menuContainer","titleContainer","titleIcon","title","contentContainer","mediaContainer","footerContainer","mediaBackdropContainer","mediaBackdrop","textArea","text","media","mediaMissing"],"defaultVariants":{"backgroundColor":"default","shape":"row"},"variantMap":{"backgroundColor":["accent","default","emphasized"],"shape":["box","fill","responsive","row"]}}

export const richCard = /* @__PURE__ */ createSlotRecipe(richCardConfig)