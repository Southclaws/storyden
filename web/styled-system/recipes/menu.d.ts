import type { ConditionalValue } from '../types/system';
import type { SlotRecipeRuntimeFn, RecipeVariantMap } from '../types/recipe';

export type MenuVariant = {
  size?: "lg" | "md" | "sm" | "xs"
}

export type MenuVariantProps = {
  [K in keyof MenuVariant]?: ConditionalValue<MenuVariant[K]>
}

export type MenuVariantMap = RecipeVariantMap<MenuVariant>

export type MenuSlot = "arrow" | "arrowTip" | "content" | "contextTrigger" | "indicator" | "item" | "itemGroup" | "itemGroupLabel" | "itemIndicator" | "itemText" | "positioner" | "separator" | "trigger" | "triggerItem"

export type MenuRecipe = SlotRecipeRuntimeFn<MenuSlot, MenuVariantProps, MenuVariantMap>

export declare const menu: MenuRecipe;