import type { ConditionalValue } from '../types/system';
import type { RecipeRuntimeFn, RecipeVariantMap } from '../types/recipe';

export type CardRowsVariant = {}

export type CardRowsVariantProps = {
  [K in keyof CardRowsVariant]?: ConditionalValue<CardRowsVariant[K]>
}

export type CardRowsVariantMap = RecipeVariantMap<CardRowsVariant>

export type CardRowsRecipe = RecipeRuntimeFn<CardRowsVariantProps, CardRowsVariantMap>

export declare const cardRows: CardRowsRecipe;