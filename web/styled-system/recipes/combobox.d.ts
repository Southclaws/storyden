import type { ConditionalValue } from '../types/system';
import type { SlotRecipeRuntimeFn, RecipeVariantMap } from '../types/recipe';

export type ComboboxVariant = {
  size?: "lg" | "md" | "sm" | "xs"
}

export type ComboboxVariantProps = {
  [K in keyof ComboboxVariant]?: ConditionalValue<ComboboxVariant[K]>
}

export type ComboboxVariantMap = RecipeVariantMap<ComboboxVariant>

export type ComboboxSlot = "root" | "clearTrigger" | "content" | "control" | "input" | "item" | "itemGroup" | "itemGroupLabel" | "itemIndicator" | "itemText" | "label" | "list" | "positioner" | "trigger" | "empty"

export type ComboboxRecipe = SlotRecipeRuntimeFn<ComboboxSlot, ComboboxVariantProps, ComboboxVariantMap>

export declare const combobox: ComboboxRecipe;