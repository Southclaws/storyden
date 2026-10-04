import type { ConditionalValue } from '../types/system';
import type { RecipeRuntimeFn, RecipeVariantMap } from '../types/recipe';

export type InputVariant = {
  size?: "lg" | "md" | "sm"
  variant?: "ghost" | "inset" | "outline"
}

export type InputVariantProps = {
  [K in keyof InputVariant]?: InputVariant[K] | undefined
}

export type InputVariantMap = RecipeVariantMap<InputVariant>

export type InputRecipe = RecipeRuntimeFn<InputVariantProps, InputVariantMap>

export declare const input: InputRecipe;