import type { ConditionalValue } from '../types/system';
import type { SlotRecipeRuntimeFn, RecipeVariantMap } from '../types/recipe';

export type DragTreeVariant = {
  variant?: "clamped" | "scrollable"
}

export type DragTreeVariantProps = {
  [K in keyof DragTreeVariant]?: ConditionalValue<DragTreeVariant[K]>
}

export type DragTreeVariantMap = RecipeVariantMap<DragTreeVariant>

export type DragTreeSlot = "branch" | "branchContent" | "branchControl" | "branchIndentGuide" | "branchIndicator" | "branchText" | "branchTrigger" | "item" | "itemIndicator" | "itemText" | "label" | "nodeCheckbox" | "nodeRenameInput" | "root" | "tree" | "actions" | "link" | "row" | "dropIndicator" | "dropIndicatorLine"

export type DragTreeRecipe = SlotRecipeRuntimeFn<DragTreeSlot, DragTreeVariantProps, DragTreeVariantMap>

export declare const dragTree: DragTreeRecipe;