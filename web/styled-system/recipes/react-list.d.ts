import type { ConditionalValue } from '../types/system';
import type { ConfigSlotRecipeRuntimeFn, RecipeVariantMap } from '../types/recipe';

export type ReactListVariant = {}

export type ReactListVariantProps = {
  [K in keyof ReactListVariant]?: ConditionalValue<ReactListVariant[K]>
}

export type ReactListVariantMap = RecipeVariantMap<ReactListVariant>

export type ReactListSlot = "root" | "trigger" | "count" | "picker"

export type ReactListRecipe = ConfigSlotRecipeRuntimeFn<ReactListSlot, ReactListVariantProps, ReactListVariantMap>

export declare const reactList: ReactListRecipe;