import type { ConditionalValue } from '../types/system';
import type { RecipeRuntimeFn, RecipeVariantMap } from '../types/recipe';

export type BadgeVariant = {
  size?: "lg" | "md" | "sm"
  variant?: "outline" | "solid" | "subtle"
}

export type BadgeVariantProps = {
  [K in keyof BadgeVariant]?: ConditionalValue<BadgeVariant[K]>
}

export type BadgeVariantMap = RecipeVariantMap<BadgeVariant>

export type BadgeRecipe = RecipeRuntimeFn<BadgeVariantProps, BadgeVariantMap>

export declare const badge: BadgeRecipe;