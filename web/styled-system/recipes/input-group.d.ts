import type { ConditionalValue } from '../types/system';
import type { SlotRecipeRuntimeFn, RecipeVariantMap } from '../types/recipe';

export type InputGroupVariant = {
  size?: "lg" | "md" | "sm" | "xl" | "xs"
}

export type InputGroupVariantProps = {
  [K in keyof InputGroupVariant]?: ConditionalValue<InputGroupVariant[K]>
}

export type InputGroupVariantMap = RecipeVariantMap<InputGroupVariant>

export type InputGroupSlot = "root" | "element"

export type InputGroupRecipe = SlotRecipeRuntimeFn<InputGroupSlot, InputGroupVariantProps, InputGroupVariantMap>

export declare const inputGroup: InputGroupRecipe;