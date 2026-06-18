import type { ConditionalValue } from '../types/system';
import type { SlotRecipeRuntimeFn, RecipeVariantMap } from '../types/recipe';

export type TabsVariant = {
  size?: "lg" | "md" | "sm"
  variant?: "enclosed" | "line" | "outline"
}

export type TabsVariantProps = {
  [K in keyof TabsVariant]?: TabsVariant[K]
}

export type TabsVariantMap = RecipeVariantMap<TabsVariant>

export type TabsSlot = "root" | "list" | "trigger" | "content" | "indicator"

export type TabsRecipe = SlotRecipeRuntimeFn<TabsSlot, TabsVariantProps, TabsVariantMap>

export declare const tabs: TabsRecipe;