import type { ConditionalValue } from '../types/system';
import type { SlotRecipeRuntimeFn, RecipeVariantMap } from '../types/recipe';

export type ClipboardVariant = {}

export type ClipboardVariantProps = {
  [K in keyof ClipboardVariant]?: ConditionalValue<ClipboardVariant[K]>
}

export type ClipboardVariantMap = RecipeVariantMap<ClipboardVariant>

export type ClipboardSlot = "root" | "control" | "trigger" | "indicator" | "input" | "label"

export type ClipboardRecipe = SlotRecipeRuntimeFn<ClipboardSlot, ClipboardVariantProps, ClipboardVariantMap>

export declare const clipboard: ClipboardRecipe;