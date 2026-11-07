import type { ConfigRecipeRuntimeFn, RecipeVariantMap } from '../types/recipe';

export type InputVariant = {
  size?: "sm" | "md" | "lg"
  variant?: "outline" | "ghost" | "inset"
}

export type InputVariantProps = {
  [K in keyof InputVariant]?: InputVariant[K] | undefined
}

export type InputVariantMap = RecipeVariantMap<InputVariant>

export type InputRecipe = ConfigRecipeRuntimeFn<InputVariantProps, InputVariantMap>

export declare const input: InputRecipe;