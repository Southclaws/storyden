import type { ConditionalValue } from '../types/system';
import type { SlotRecipeRuntimeFn, RecipeVariantMap } from '../types/recipe';

export type ContentOutlineVariant = {
  placement?: "gutter" | "inline"
}

export type ContentOutlineVariantProps = {
  [K in keyof ContentOutlineVariant]?: ConditionalValue<ContentOutlineVariant[K]>
}

export type ContentOutlineVariantMap = RecipeVariantMap<ContentOutlineVariant>

export type ContentOutlineSlot = "root" | "body" | "rail" | "trigger" | "ticks" | "tick" | "content" | "viewport" | "navigation" | "replyLabel" | "item" | "itemContent" | "byline" | "label" | "timestamp" | "snippet"

export type ContentOutlineRecipe = SlotRecipeRuntimeFn<ContentOutlineSlot, ContentOutlineVariantProps, ContentOutlineVariantMap>

export declare const contentOutline: ContentOutlineRecipe;