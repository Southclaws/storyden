import type { ConditionalValue } from '../types/system';
import type { ConfigSlotRecipeRuntimeFn, RecipeVariantMap } from '../types/recipe';

export type ComboboxVariant = {
  size?: "sm" | "md" | "lg"
}

export type ComboboxVariantProps = {
  [K in keyof ComboboxVariant]?: ConditionalValue<ComboboxVariant[K]>
}

export type ComboboxVariantMap = RecipeVariantMap<ComboboxVariant>

export type ComboboxSlot = "root" | "clearTrigger" | "content" | "control" | "input" | "item" | "itemGroup" | "itemGroupLabel" | "itemIndicator" | "itemText" | "label" | "list" | "positioner" | "trigger" | "empty"

export type ComboboxRecipe = ConfigSlotRecipeRuntimeFn<ComboboxSlot, ComboboxVariantProps, ComboboxVariantMap>

export declare const combobox: ComboboxRecipe;