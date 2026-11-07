import type { ConditionalValue } from '../types/system';
import type { ConfigSlotRecipeRuntimeFn, RecipeVariantMap } from '../types/recipe';

export type DragTreeVariant = {
  variant?: "clamped" | "scrollable"
}

export type DragTreeVariantProps = {
  [K in keyof DragTreeVariant]?: ConditionalValue<DragTreeVariant[K]>
}

export type DragTreeVariantMap = RecipeVariantMap<DragTreeVariant>

export type DragTreeSlot = "branch" | "branchContent" | "branchControl" | "branchIndentGuide" | "branchIndicator" | "branchText" | "branchTrigger" | "item" | "itemIndicator" | "itemText" | "label" | "nodeCheckbox" | "nodeRenameInput" | "root" | "tree" | "actions" | "link" | "row" | "dropIndicator" | "dropIndicatorLine" | "instructions" | "movePreview" | "moveTitle" | "moveHelp" | "moveError"

export type DragTreeRecipe = ConfigSlotRecipeRuntimeFn<DragTreeSlot, DragTreeVariantProps, DragTreeVariantMap>

export declare const dragTree: DragTreeRecipe;