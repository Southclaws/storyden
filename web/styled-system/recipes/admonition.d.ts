import type { ConditionalValue } from '../types/system';
import type { RecipeRuntimeFn, RecipeVariantMap } from '../types/recipe';

export type AdmonitionVariant = {
  kind?: "failure" | "neutral" | "success"
}

export type AdmonitionVariantProps = {
  [K in keyof AdmonitionVariant]?: ConditionalValue<AdmonitionVariant[K]>
}

export type AdmonitionVariantMap = RecipeVariantMap<AdmonitionVariant>

export type AdmonitionRecipe = RecipeRuntimeFn<AdmonitionVariantProps, AdmonitionVariantMap>

export declare const admonition: AdmonitionRecipe;