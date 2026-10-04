import type { ConditionalValue } from '../types/system';
import type { SlotRecipeRuntimeFn, RecipeVariantMap } from '../types/recipe';

export type SectionNavigationVariant = {}

export type SectionNavigationVariantProps = {
  [K in keyof SectionNavigationVariant]?: ConditionalValue<SectionNavigationVariant[K]>
}

export type SectionNavigationVariantMap = RecipeVariantMap<SectionNavigationVariant>

export type SectionNavigationSlot = "root" | "trigger" | "triggerLabel" | "triggerIndicator" | "positioner" | "content" | "navigation" | "section" | "sectionLabel" | "items" | "item" | "link" | "linkIndicator"

export type SectionNavigationRecipe = SlotRecipeRuntimeFn<SectionNavigationSlot, SectionNavigationVariantProps, SectionNavigationVariantMap>

export declare const sectionNavigation: SectionNavigationRecipe;