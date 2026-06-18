import type { ConditionalValue } from '../types/system';
import type { SlotRecipeRuntimeFn, RecipeVariantMap } from '../types/recipe';

export type SelectVariant = {
  size?: "lg" | "md" | "sm" | "xs"
  variant?: "ghost" | "outline"
}

export type SelectVariantProps = {
  [K in keyof SelectVariant]?: ConditionalValue<SelectVariant[K]>
}

export type SelectVariantMap = RecipeVariantMap<SelectVariant>

export type SelectSlot = "label" | "positioner" | "trigger" | "indicator" | "clearTrigger" | "item" | "itemText" | "itemIndicator" | "itemGroup" | "itemGroupLabel" | "list" | "content" | "root" | "control" | "valueText"

export type SelectRecipe = SlotRecipeRuntimeFn<SelectSlot, SelectVariantProps, SelectVariantMap>

export declare const select: SelectRecipe;