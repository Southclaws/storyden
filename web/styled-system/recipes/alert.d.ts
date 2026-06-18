import type { ConditionalValue } from '../types/system';
import type { SlotRecipeRuntimeFn, RecipeVariantMap } from '../types/recipe';

export type AlertVariant = {}

export type AlertVariantProps = {
  [K in keyof AlertVariant]?: ConditionalValue<AlertVariant[K]>
}

export type AlertVariantMap = RecipeVariantMap<AlertVariant>

export type AlertSlot = "root" | "content" | "description" | "icon" | "title"

export type AlertRecipe = SlotRecipeRuntimeFn<AlertSlot, AlertVariantProps, AlertVariantMap>

export declare const alert: AlertRecipe;