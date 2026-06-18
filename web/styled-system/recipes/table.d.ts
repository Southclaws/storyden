import type { ConditionalValue } from '../types/system';
import type { SlotRecipeRuntimeFn, RecipeVariantMap } from '../types/recipe';

export type TableVariant = {
  size?: "md" | "sm"
  variant?: "dense" | "plain"
}

export type TableVariantProps = {
  [K in keyof TableVariant]?: ConditionalValue<TableVariant[K]>
}

export type TableVariantMap = RecipeVariantMap<TableVariant>

export type TableSlot = "root" | "body" | "cell" | "footer" | "head" | "header" | "row" | "caption"

export type TableRecipe = SlotRecipeRuntimeFn<TableSlot, TableVariantProps, TableVariantMap>

export declare const table: TableRecipe;