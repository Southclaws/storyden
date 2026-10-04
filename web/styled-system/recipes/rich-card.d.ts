import type { ConditionalValue } from '../types/system';
import type { SlotRecipeRuntimeFn, RecipeVariantMap } from '../types/recipe';

export type RichCardVariant = {
  backgroundColor?: "accent" | "default" | "emphasized"
  shape?: "box" | "fill" | "responsive" | "row"
}

export type RichCardVariantProps = {
  [K in keyof RichCardVariant]?: ConditionalValue<RichCardVariant[K]>
}

export type RichCardVariantMap = RecipeVariantMap<RichCardVariant>

export type RichCardSlot = "container" | "root" | "headerContainer" | "menuContainer" | "titleContainer" | "titleIcon" | "title" | "contentContainer" | "mediaContainer" | "footerContainer" | "mediaBackdropContainer" | "mediaBackdrop" | "textArea" | "text" | "media" | "mediaMissing"

export type RichCardRecipe = SlotRecipeRuntimeFn<RichCardSlot, RichCardVariantProps, RichCardVariantMap>

export declare const richCard: RichCardRecipe;