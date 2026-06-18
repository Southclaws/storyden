import type { ConditionalValue } from '../types/system';
import type { RecipeRuntimeFn, RecipeVariantMap } from '../types/recipe';

export type RichCardVariant = {
  backgroundColor?: "accent" | "default" | "emphasized"
  shape?: "box" | "fill" | "responsive" | "row"
}

export type RichCardVariantProps = {
  [K in keyof RichCardVariant]?: ConditionalValue<RichCardVariant[K]>
}

export type RichCardVariantMap = RecipeVariantMap<RichCardVariant>

export type RichCardRecipe = RecipeRuntimeFn<RichCardVariantProps, RichCardVariantMap>

export declare const richCard: RichCardRecipe;