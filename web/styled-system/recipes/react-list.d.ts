import type { ConditionalValue } from '../types/system';
import type { SlotRecipeRuntimeFn, RecipeVariantMap } from '../types/recipe';

export type ReactListVariant = {}

export type ReactListVariantProps = {
  [K in keyof ReactListVariant]?: ConditionalValue<ReactListVariant[K]>
}

export type ReactListVariantMap = RecipeVariantMap<ReactListVariant>

export type ReactListSlot = "root" | "reaction" | "count" | "picker"

export type ReactListRecipe = SlotRecipeRuntimeFn<ReactListSlot, ReactListVariantProps, ReactListVariantMap>

export declare const reactList: ReactListRecipe;