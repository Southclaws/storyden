import type { ConditionalValue } from '../types/system';
import type { ConfigSlotRecipeRuntimeFn, RecipeVariantMap } from '../types/recipe';

export type AlertVariant = {
  tone?: "danger" | "info" | "success" | "warning"
}

export type AlertVariantProps = {
  [K in keyof AlertVariant]?: ConditionalValue<AlertVariant[K]>
}

export type AlertVariantMap = RecipeVariantMap<AlertVariant>

export type AlertSlot = "root" | "content" | "description" | "icon" | "title"

export type AlertRecipe = ConfigSlotRecipeRuntimeFn<AlertSlot, AlertVariantProps, AlertVariantMap>

export declare const alert: AlertRecipe;