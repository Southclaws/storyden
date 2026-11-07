import type { ConditionalValue } from '../types/system';
import type { ConfigSlotRecipeRuntimeFn, RecipeVariantMap } from '../types/recipe';

export type CheckboxVariant = {
  size?: "sm" | "md" | "lg"
}

export type CheckboxVariantProps = {
  [K in keyof CheckboxVariant]?: ConditionalValue<CheckboxVariant[K]>
}

export type CheckboxVariantMap = RecipeVariantMap<CheckboxVariant>

export type CheckboxSlot = "root" | "label" | "control" | "indicator" | "group"

export type CheckboxRecipe = ConfigSlotRecipeRuntimeFn<CheckboxSlot, CheckboxVariantProps, CheckboxVariantMap>

export declare const checkbox: CheckboxRecipe;