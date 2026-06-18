import type { ConditionalValue } from '../types/system';
import type { RecipeRuntimeFn, RecipeVariantMap } from '../types/recipe';

export type CheckboxVariant = {
  size?: "lg" | "md" | "sm"
}

export type CheckboxVariantProps = {
  [K in keyof CheckboxVariant]?: ConditionalValue<CheckboxVariant[K]>
}

export type CheckboxVariantMap = RecipeVariantMap<CheckboxVariant>

export type CheckboxRecipe = RecipeRuntimeFn<CheckboxVariantProps, CheckboxVariantMap>

export declare const checkbox: CheckboxRecipe;