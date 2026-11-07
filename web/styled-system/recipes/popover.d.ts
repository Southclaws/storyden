import type { ConditionalValue } from '../types/system';
import type { ConfigSlotRecipeRuntimeFn, RecipeVariantMap } from '../types/recipe';

export type PopoverVariant = {}

export type PopoverVariantProps = {
  [K in keyof PopoverVariant]?: ConditionalValue<PopoverVariant[K]>
}

export type PopoverVariantMap = RecipeVariantMap<PopoverVariant>

export type PopoverSlot = "arrow" | "arrowTip" | "anchor" | "trigger" | "indicator" | "positioner" | "content" | "title" | "description" | "closeTrigger"

export type PopoverRecipe = ConfigSlotRecipeRuntimeFn<PopoverSlot, PopoverVariantProps, PopoverVariantMap>

export declare const popover: PopoverRecipe;