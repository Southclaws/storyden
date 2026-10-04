import type { ConditionalValue } from '../types/system';
import type { SlotRecipeRuntimeFn, RecipeVariantMap } from '../types/recipe';

export type CheckboxVariant = {
  size?: "lg" | "md" | "sm"
}

export type CheckboxVariantProps = {
  [K in keyof CheckboxVariant]?: ConditionalValue<CheckboxVariant[K]>
}

export type CheckboxVariantMap = RecipeVariantMap<CheckboxVariant>

export type CheckboxSlot = "root" | "label" | "control" | "indicator" | "group"

export type CheckboxRecipe = SlotRecipeRuntimeFn<CheckboxSlot, CheckboxVariantProps, CheckboxVariantMap>

export declare const checkbox: CheckboxRecipe;