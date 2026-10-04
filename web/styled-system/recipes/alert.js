import { createSlotRecipe } from './runtime';

const alertConfig = {"name":"alert","slots":["root","content","description","icon","title"],"variantMap":{"tone":["danger","info","success","warning"]}}

export const alert = /* @__PURE__ */ createSlotRecipe(alertConfig)