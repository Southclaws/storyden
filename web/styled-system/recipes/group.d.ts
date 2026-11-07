import type { ConfigRecipeRuntimeFn, RecipeVariantMap } from '../types/recipe';

export type GroupVariant = {
  orientation?: "horizontal" | "vertical"
  attached?: boolean
}

export type GroupVariantProps = {
  [K in keyof GroupVariant]?: GroupVariant[K] | undefined
}

export type GroupVariantMap = RecipeVariantMap<GroupVariant>

export type GroupRecipe = ConfigRecipeRuntimeFn<GroupVariantProps, GroupVariantMap>

export declare const group: GroupRecipe;