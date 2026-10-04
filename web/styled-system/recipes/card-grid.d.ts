import type { ConditionalValue } from '../types/system';
import type { SlotRecipeRuntimeFn, RecipeVariantMap } from '../types/recipe';

export type CardGridVariant = {}

export type CardGridVariantProps = {
  [K in keyof CardGridVariant]?: ConditionalValue<CardGridVariant[K]>
}

export type CardGridVariantMap = RecipeVariantMap<CardGridVariant>

export type CardGridSlot = "container" | "grid"

export type CardGridRecipe = SlotRecipeRuntimeFn<CardGridSlot, CardGridVariantProps, CardGridVariantMap>

export declare const cardGrid: CardGridRecipe;