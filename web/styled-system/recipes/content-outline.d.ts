import type { ConditionalValue } from '../types/system';
import type { ConfigSlotRecipeRuntimeFn, RecipeVariantMap } from '../types/recipe';

export type ContentOutlineVariant = {
  placement?: "inline" | "gutter"
}

export type ContentOutlineVariantProps = {
  [K in keyof ContentOutlineVariant]?: ConditionalValue<ContentOutlineVariant[K]>
}

export type ContentOutlineVariantMap = RecipeVariantMap<ContentOutlineVariant>

export type ContentOutlineSlot = "root" | "body" | "rail" | "trigger" | "ticks" | "tick" | "content" | "viewport" | "navigation" | "replyLabel" | "item" | "itemContent" | "byline" | "label" | "timestamp" | "snippet"

export type ContentOutlineRecipe = ConfigSlotRecipeRuntimeFn<ContentOutlineSlot, ContentOutlineVariantProps, ContentOutlineVariantMap>

export declare const contentOutline: ContentOutlineRecipe;