import type { ConditionalValue } from '../types/system';
import type { ConfigSlotRecipeRuntimeFn, RecipeVariantMap } from '../types/recipe';

export type PinInputVariant = {
  size?: "sm" | "md" | "lg"
}

export type PinInputVariantProps = {
  [K in keyof PinInputVariant]?: ConditionalValue<PinInputVariant[K]>
}

export type PinInputVariantMap = RecipeVariantMap<PinInputVariant>

export type PinInputSlot = "root" | "label" | "input" | "control"

export type PinInputRecipe = ConfigSlotRecipeRuntimeFn<PinInputSlot, PinInputVariantProps, PinInputVariantMap>

export declare const pinInput: PinInputRecipe;