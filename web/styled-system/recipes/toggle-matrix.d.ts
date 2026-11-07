import type { ConditionalValue } from '../types/system';
import type { ConfigSlotRecipeRuntimeFn, RecipeVariantMap } from '../types/recipe';

export type ToggleMatrixVariant = {}

export type ToggleMatrixVariantProps = {
  [K in keyof ToggleMatrixVariant]?: ConditionalValue<ToggleMatrixVariant[K]>
}

export type ToggleMatrixVariantMap = RecipeVariantMap<ToggleMatrixVariant>

export type ToggleMatrixSlot = "root" | "columnHeader" | "rowHeader" | "description" | "cell" | "item"

export type ToggleMatrixRecipe = ConfigSlotRecipeRuntimeFn<ToggleMatrixSlot, ToggleMatrixVariantProps, ToggleMatrixVariantMap>

export declare const toggleMatrix: ToggleMatrixRecipe;