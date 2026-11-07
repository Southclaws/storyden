import type { ConditionalValue } from '../types/system';
import type { ConfigRecipeRuntimeFn, RecipeVariantMap } from '../types/recipe';

export type BadgeVariant = {
  variant?: "solid" | "subtle" | "outline"
  size?: "sm" | "md" | "lg"
}

export type BadgeVariantProps = {
  [K in keyof BadgeVariant]?: ConditionalValue<BadgeVariant[K]>
}

export type BadgeVariantMap = RecipeVariantMap<BadgeVariant>

export type BadgeRecipe = ConfigRecipeRuntimeFn<BadgeVariantProps, BadgeVariantMap>

export declare const badge: BadgeRecipe;