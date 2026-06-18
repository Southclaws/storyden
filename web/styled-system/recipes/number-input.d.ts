import type { ConditionalValue } from '../types/system';
import type { SlotRecipeRuntimeFn, RecipeVariantMap } from '../types/recipe';

export type NumberInputVariant = {
  size?: "lg" | "md" | "sm" | "xl"
  variant?: "ghost" | "outline"
}

export type NumberInputVariantProps = {
  [K in keyof NumberInputVariant]?: ConditionalValue<NumberInputVariant[K]>
}

export type NumberInputVariantMap = RecipeVariantMap<NumberInputVariant>

export type NumberInputSlot = "root" | "label" | "input" | "control" | "valueText" | "incrementTrigger" | "decrementTrigger" | "scrubber"

export type NumberInputRecipe = SlotRecipeRuntimeFn<NumberInputSlot, NumberInputVariantProps, NumberInputVariantMap>

export declare const numberInput: NumberInputRecipe;