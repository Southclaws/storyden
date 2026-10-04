import type { ConditionalValue } from '../types/system';
import type { RecipeRuntimeFn, RecipeVariantMap } from '../types/recipe';

export type ButtonVariant = {
  intent?: "destructive" | "success" | "warning"
  size?: "lg" | "md" | "sm"
  variant?: "ghost" | "outline" | "plain" | "solid" | "subtle"
}

export type ButtonVariantProps = {
  [K in keyof ButtonVariant]?: ConditionalValue<ButtonVariant[K]>
}

export type ButtonVariantMap = RecipeVariantMap<ButtonVariant>

export type ButtonRecipe = RecipeRuntimeFn<ButtonVariantProps, ButtonVariantMap>

export declare const button: ButtonRecipe;