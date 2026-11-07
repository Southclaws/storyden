import type { ConditionalValue } from '../types/system';
import type { ConfigSlotRecipeRuntimeFn, RecipeVariantMap } from '../types/recipe';

export type SelectVariant = {
  variant?: "outline" | "ghost"
  size?: "sm" | "md" | "lg"
}

export type SelectVariantProps = {
  [K in keyof SelectVariant]?: ConditionalValue<SelectVariant[K]>
}

export type SelectVariantMap = RecipeVariantMap<SelectVariant>

export type SelectSlot = "label" | "positioner" | "trigger" | "indicator" | "clearTrigger" | "item" | "itemText" | "itemIndicator" | "itemGroup" | "itemGroupLabel" | "list" | "content" | "root" | "control" | "valueText"

export type SelectRecipe = ConfigSlotRecipeRuntimeFn<SelectSlot, SelectVariantProps, SelectVariantMap>

export declare const select: SelectRecipe;