import type { ConditionalValue } from '../types/system';
import type { SlotRecipeRuntimeFn, RecipeVariantMap } from '../types/recipe';

export type TooltipVariant = {}

export type TooltipVariantProps = {
  [K in keyof TooltipVariant]?: ConditionalValue<TooltipVariant[K]>
}

export type TooltipVariantMap = RecipeVariantMap<TooltipVariant>

export type TooltipSlot = "trigger" | "arrow" | "arrowTip" | "positioner" | "content"

export type TooltipRecipe = SlotRecipeRuntimeFn<TooltipSlot, TooltipVariantProps, TooltipVariantMap>

export declare const tooltip: TooltipRecipe;