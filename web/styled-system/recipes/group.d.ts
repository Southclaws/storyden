import type { ConditionalValue } from '../types/system';
import type { RecipeRuntimeFn, RecipeVariantMap } from '../types/recipe';

export type GroupVariant = {
  attached?: boolean
  grow?: boolean
  orientation?: "horizontal" | "vertical"
}

export type GroupVariantProps = {
  [K in keyof GroupVariant]?: GroupVariant[K]
}

export type GroupVariantMap = RecipeVariantMap<GroupVariant>

export type GroupRecipe = RecipeRuntimeFn<GroupVariantProps, GroupVariantMap>

export declare const group: GroupRecipe;