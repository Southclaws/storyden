import type { ConditionalValue } from '../types/system';
import type { SlotRecipeRuntimeFn, RecipeVariantMap } from '../types/recipe';

export type FloatingPanelVariant = {}

export type FloatingPanelVariantProps = {
  [K in keyof FloatingPanelVariant]?: ConditionalValue<FloatingPanelVariant[K]>
}

export type FloatingPanelVariantMap = RecipeVariantMap<FloatingPanelVariant>

export type FloatingPanelSlot = "trigger" | "positioner" | "content" | "header" | "body" | "title" | "resizeTrigger" | "dragTrigger" | "stageTrigger" | "closeTrigger" | "control"

export type FloatingPanelRecipe = SlotRecipeRuntimeFn<FloatingPanelSlot, FloatingPanelVariantProps, FloatingPanelVariantMap>

export declare const floatingPanel: FloatingPanelRecipe;