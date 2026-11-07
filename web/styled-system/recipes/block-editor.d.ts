import type { ConditionalValue } from '../types/system';
import type { ConfigSlotRecipeRuntimeFn, RecipeVariantMap } from '../types/recipe';

export type BlockEditorVariant = {}

export type BlockEditorVariantProps = {
  [K in keyof BlockEditorVariant]?: ConditionalValue<BlockEditorVariant[K]>
}

export type BlockEditorVariantMap = RecipeVariantMap<BlockEditorVariant>

export type BlockEditorSlot = "root" | "gutter" | "handle" | "menuTrigger" | "menuAnchor" | "menuContent" | "content"

export type BlockEditorRecipe = ConfigSlotRecipeRuntimeFn<BlockEditorSlot, BlockEditorVariantProps, BlockEditorVariantMap>

export declare const blockEditor: BlockEditorRecipe;