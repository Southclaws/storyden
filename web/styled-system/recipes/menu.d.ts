import type { ConditionalValue } from '../types/system';
import type { ConfigSlotRecipeRuntimeFn, RecipeVariantMap } from '../types/recipe';

export type MenuVariant = {
  size?: "xs" | "sm" | "md" | "lg"
}

export type MenuVariantProps = {
  [K in keyof MenuVariant]?: ConditionalValue<MenuVariant[K]>
}

export type MenuVariantMap = RecipeVariantMap<MenuVariant>

export type MenuSlot = "arrow" | "arrowTip" | "content" | "contextTrigger" | "indicator" | "item" | "itemGroup" | "itemGroupLabel" | "itemIndicator" | "itemText" | "positioner" | "separator" | "trigger" | "triggerItem"

export type MenuRecipe = ConfigSlotRecipeRuntimeFn<MenuSlot, MenuVariantProps, MenuVariantMap>

export declare const menu: MenuRecipe;