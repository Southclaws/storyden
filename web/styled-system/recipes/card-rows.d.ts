import type { ConditionalValue } from '../types/system';
import type { ConfigRecipeRuntimeFn, RecipeVariantMap } from '../types/recipe';

export type CardRowsVariant = {}

export type CardRowsVariantProps = {
  [K in keyof CardRowsVariant]?: ConditionalValue<CardRowsVariant[K]>
}

export type CardRowsVariantMap = RecipeVariantMap<CardRowsVariant>

export type CardRowsRecipe = ConfigRecipeRuntimeFn<CardRowsVariantProps, CardRowsVariantMap>

export declare const cardRows: CardRowsRecipe;