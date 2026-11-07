import type { ConditionalValue } from '../types/system';
import type { ConfigRecipeRuntimeFn, RecipeVariantMap } from '../types/recipe';

export type AdmonitionVariant = {
  kind?: "neutral" | "success" | "failure"
}

export type AdmonitionVariantProps = {
  [K in keyof AdmonitionVariant]?: ConditionalValue<AdmonitionVariant[K]>
}

export type AdmonitionVariantMap = RecipeVariantMap<AdmonitionVariant>

export type AdmonitionRecipe = ConfigRecipeRuntimeFn<AdmonitionVariantProps, AdmonitionVariantMap>

export declare const admonition: AdmonitionRecipe;