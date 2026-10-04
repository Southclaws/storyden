import type { ConditionalValue } from '../types/system';
import type { SlotRecipeRuntimeFn, RecipeVariantMap } from '../types/recipe';

export type BlockEditorVariant = {}

export type BlockEditorVariantProps = {
  [K in keyof BlockEditorVariant]?: ConditionalValue<BlockEditorVariant[K]>
}

export type BlockEditorVariantMap = RecipeVariantMap<BlockEditorVariant>

export type BlockEditorSlot = "root" | "gutter" | "handle" | "menuTrigger" | "menuAnchor" | "menuContent" | "content"

export type BlockEditorRecipe = SlotRecipeRuntimeFn<BlockEditorSlot, BlockEditorVariantProps, BlockEditorVariantMap>

export declare const blockEditor: BlockEditorRecipe;