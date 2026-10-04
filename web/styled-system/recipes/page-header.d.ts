import type { ConditionalValue } from '../types/system';
import type { SlotRecipeRuntimeFn, RecipeVariantMap } from '../types/recipe';

export type PageHeaderVariant = {}

export type PageHeaderVariantProps = {
  [K in keyof PageHeaderVariant]?: ConditionalValue<PageHeaderVariant[K]>
}

export type PageHeaderVariantMap = RecipeVariantMap<PageHeaderVariant>

export type PageHeaderSlot = "root" | "navigation" | "row" | "heading" | "back" | "titleGroup" | "titleRow" | "actions"

export type PageHeaderRecipe = SlotRecipeRuntimeFn<PageHeaderSlot, PageHeaderVariantProps, PageHeaderVariantMap>

export declare const pageHeader: PageHeaderRecipe;