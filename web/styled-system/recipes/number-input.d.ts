import type { ConditionalValue } from '../types/system';
import type { ConfigSlotRecipeRuntimeFn, RecipeVariantMap } from '../types/recipe';

export type NumberInputVariant = {
  variant?: "outline" | "ghost"
  size?: "sm" | "md" | "lg"
}

export type NumberInputVariantProps = {
  [K in keyof NumberInputVariant]?: ConditionalValue<NumberInputVariant[K]>
}

export type NumberInputVariantMap = RecipeVariantMap<NumberInputVariant>

export type NumberInputSlot = "root" | "label" | "input" | "control" | "valueText" | "incrementTrigger" | "decrementTrigger" | "scrubber"

export type NumberInputRecipe = ConfigSlotRecipeRuntimeFn<NumberInputSlot, NumberInputVariantProps, NumberInputVariantMap>

export declare const numberInput: NumberInputRecipe;