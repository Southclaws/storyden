import type { ConditionalValue } from '../types/system';
import type { RecipeRuntimeFn, RecipeVariantMap } from '../types/recipe';

export type InputVariant = {
  size?: "2xl" | "2xs" | "lg" | "md" | "sm" | "xl" | "xs"
  variant?: "ghost" | "outline"
}

export type InputVariantProps = {
  [K in keyof InputVariant]?: ConditionalValue<InputVariant[K]>
}

export type InputVariantMap = RecipeVariantMap<InputVariant>

export type InputRecipe = RecipeRuntimeFn<InputVariantProps, InputVariantMap>

export declare const input: InputRecipe;