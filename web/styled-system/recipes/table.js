import { createSlotRecipe } from './runtime';

const tableConfig = {"name":"table","slots":["root","body","cell","footer","head","header","row","caption"],"defaultVariants":{"size":"md","variant":"plain"},"variantMap":{"variant":["dense","plain"],"size":["sm","md"]}}

export const table = /* @__PURE__ */ createSlotRecipe(tableConfig)