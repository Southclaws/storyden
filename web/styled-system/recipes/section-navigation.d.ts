import type { ConditionalValue } from '../types/system';
import type { ConfigSlotRecipeRuntimeFn, RecipeVariantMap } from '../types/recipe';

export type SectionNavigationVariant = {}

export type SectionNavigationVariantProps = {
  [K in keyof SectionNavigationVariant]?: ConditionalValue<SectionNavigationVariant[K]>
}

export type SectionNavigationVariantMap = RecipeVariantMap<SectionNavigationVariant>

export type SectionNavigationSlot = "root" | "trigger" | "triggerLabel" | "triggerIndicator" | "positioner" | "content" | "navigation" | "section" | "sectionLabel" | "items" | "item" | "link" | "linkIndicator"

export type SectionNavigationRecipe = ConfigSlotRecipeRuntimeFn<SectionNavigationSlot, SectionNavigationVariantProps, SectionNavigationVariantMap>

export declare const sectionNavigation: SectionNavigationRecipe;