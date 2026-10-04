import type { ConditionalValue } from '../types/system';
import type { SlotRecipeRuntimeFn, RecipeVariantMap } from '../types/recipe';

export type PinInputVariant = {
  size?: "lg" | "md" | "sm"
}

export type PinInputVariantProps = {
  [K in keyof PinInputVariant]?: ConditionalValue<PinInputVariant[K]>
}

export type PinInputVariantMap = RecipeVariantMap<PinInputVariant>

export type PinInputSlot = "root" | "label" | "input" | "control"

export type PinInputRecipe = SlotRecipeRuntimeFn<PinInputSlot, PinInputVariantProps, PinInputVariantMap>

export declare const pinInput: PinInputRecipe;