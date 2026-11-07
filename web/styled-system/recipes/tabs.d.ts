import type { ConfigSlotRecipeRuntimeFn, RecipeVariantMap } from '../types/recipe';

export type TabsVariant = {
  variant?: "enclosed" | "line" | "outline"
  size?: "sm" | "md" | "lg"
}

export type TabsVariantProps = {
  [K in keyof TabsVariant]?: TabsVariant[K] | undefined
}

export type TabsVariantMap = RecipeVariantMap<TabsVariant>

export type TabsSlot = "root" | "list" | "trigger" | "content" | "indicator"

export type TabsRecipe = ConfigSlotRecipeRuntimeFn<TabsSlot, TabsVariantProps, TabsVariantMap>

export declare const tabs: TabsRecipe;