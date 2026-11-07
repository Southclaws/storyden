import type { ConditionalValue } from '../types/system';
import type { ConfigSlotRecipeRuntimeFn, RecipeVariantMap } from '../types/recipe';

export type TableVariant = {
  variant?: "dense" | "plain"
  size?: "sm" | "md"
}

export type TableVariantProps = {
  [K in keyof TableVariant]?: ConditionalValue<TableVariant[K]>
}

export type TableVariantMap = RecipeVariantMap<TableVariant>

export type TableSlot = "root" | "body" | "cell" | "footer" | "head" | "header" | "row" | "caption"

export type TableRecipe = ConfigSlotRecipeRuntimeFn<TableSlot, TableVariantProps, TableVariantMap>

export declare const table: TableRecipe;