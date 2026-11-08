import type { ConfigSlotRecipeRuntimeFn, RecipeVariantMap } from '../types/recipe';

export type RichCardVariant = {
  indicator?: boolean
  titleWeight?: "semibold" | "bold"
  hasMedia?: boolean
  backgroundColor?: "default" | "emphasized" | "accent"
  shape?: "minimal" | "row" | "responsive" | "box" | "fill"
}

export type RichCardVariantProps = {
  [K in keyof RichCardVariant]?: RichCardVariant[K] | undefined
}

export type RichCardVariantMap = RecipeVariantMap<RichCardVariant>

export type RichCardSlot = "container" | "root" | "headerContainer" | "menuContainer" | "titleContainer" | "titleIcon" | "indicator" | "title" | "contentContainer" | "mediaContainer" | "footerContainer" | "mediaBackdropContainer" | "mediaBackdrop" | "textArea" | "text" | "media" | "mediaMissing"

export type RichCardRecipe = ConfigSlotRecipeRuntimeFn<RichCardSlot, RichCardVariantProps, RichCardVariantMap>

export declare const richCard: RichCardRecipe;