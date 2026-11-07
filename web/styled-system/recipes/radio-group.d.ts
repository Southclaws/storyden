import type { ConditionalValue } from '../types/system';
import type { ConfigSlotRecipeRuntimeFn, RecipeVariantMap } from '../types/recipe';

export type RadioGroupVariant = {
  size?: "sm" | "md" | "lg"
}

export type RadioGroupVariantProps = {
  [K in keyof RadioGroupVariant]?: ConditionalValue<RadioGroupVariant[K]>
}

export type RadioGroupVariantMap = RecipeVariantMap<RadioGroupVariant>

export type RadioGroupSlot = "root" | "label" | "item" | "itemText" | "itemControl" | "indicator"

export type RadioGroupRecipe = ConfigSlotRecipeRuntimeFn<RadioGroupSlot, RadioGroupVariantProps, RadioGroupVariantMap>

export declare const radioGroup: RadioGroupRecipe;