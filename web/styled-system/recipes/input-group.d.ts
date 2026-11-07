import type { ConditionalValue } from '../types/system';
import type { ConfigSlotRecipeRuntimeFn, RecipeVariantMap } from '../types/recipe';

export type InputGroupVariant = {
  size?: "sm" | "md" | "lg"
}

export type InputGroupVariantProps = {
  [K in keyof InputGroupVariant]?: ConditionalValue<InputGroupVariant[K]>
}

export type InputGroupVariantMap = RecipeVariantMap<InputGroupVariant>

export type InputGroupSlot = "root" | "element"

export type InputGroupRecipe = ConfigSlotRecipeRuntimeFn<InputGroupSlot, InputGroupVariantProps, InputGroupVariantMap>

export declare const inputGroup: InputGroupRecipe;