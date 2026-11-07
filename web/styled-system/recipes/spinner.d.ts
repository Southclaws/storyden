import type { ConditionalValue } from '../types/system';
import type { ConfigRecipeRuntimeFn, RecipeVariantMap } from '../types/recipe';

export type SpinnerVariant = {
  size?: "sm" | "md" | "lg"
}

export type SpinnerVariantProps = {
  [K in keyof SpinnerVariant]?: ConditionalValue<SpinnerVariant[K]>
}

export type SpinnerVariantMap = RecipeVariantMap<SpinnerVariant>

export type SpinnerRecipe = ConfigRecipeRuntimeFn<SpinnerVariantProps, SpinnerVariantMap>

export declare const spinner: SpinnerRecipe;