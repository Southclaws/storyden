import type { ConditionalValue } from '../types/system';
import type { SlotRecipeRuntimeFn, RecipeVariantMap } from '../types/recipe';

export type TreeViewVariant = {
  variant?: "clamped" | "scrollable"
}

export type TreeViewVariantProps = {
  [K in keyof TreeViewVariant]?: ConditionalValue<TreeViewVariant[K]>
}

export type TreeViewVariantMap = RecipeVariantMap<TreeViewVariant>

export type TreeViewSlot = "branch" | "branchContent" | "branchControl" | "branchIndentGuide" | "branchIndicator" | "branchText" | "branchTrigger" | "item" | "itemIndicator" | "itemText" | "label" | "nodeCheckbox" | "nodeRenameInput" | "root" | "tree"

export type TreeViewRecipe = SlotRecipeRuntimeFn<TreeViewSlot, TreeViewVariantProps, TreeViewVariantMap>

export declare const treeView: TreeViewRecipe;