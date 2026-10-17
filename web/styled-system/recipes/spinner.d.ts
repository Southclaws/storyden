import type { ConditionalValue } from '../types/system';
import type { RecipeRuntimeFn, RecipeVariantMap } from '../types/recipe';

export type SpinnerVariant = {
  size?: "lg" | "md" | "sm"
}

export type SpinnerVariantProps = {
  [K in keyof SpinnerVariant]?: ConditionalValue<SpinnerVariant[K]>
}

export type SpinnerVariantMap = RecipeVariantMap<SpinnerVariant>

export type SpinnerRecipe = RecipeRuntimeFn<SpinnerVariantProps, SpinnerVariantMap>

export declare const spinner: SpinnerRecipe;