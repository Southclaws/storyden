import type { ConditionalValue } from '../types/system';
import type { ConfigRecipeRuntimeFn, RecipeVariantMap } from '../types/recipe';

export type ButtonVariant = {
  size?: "sm" | "md" | "lg"
  variant?: "solid" | "outline" | "ghost" | "subtle" | "plain"
  intent?: "success" | "warning" | "destructive"
}

export type ButtonVariantProps = {
  [K in keyof ButtonVariant]?: ConditionalValue<ButtonVariant[K]>
}

export type ButtonVariantMap = RecipeVariantMap<ButtonVariant>

export type ButtonRecipe = ConfigRecipeRuntimeFn<ButtonVariantProps, ButtonVariantMap>

export declare const button: ButtonRecipe;