import type { ConditionalValue } from '../types/system';
import type { ConfigRecipeRuntimeFn, RecipeVariantMap } from '../types/recipe';

export type HeadingInputVariant = {}

export type HeadingInputVariantProps = {
  [K in keyof HeadingInputVariant]?: ConditionalValue<HeadingInputVariant[K]>
}

export type HeadingInputVariantMap = RecipeVariantMap<HeadingInputVariant>

export type HeadingInputRecipe = ConfigRecipeRuntimeFn<HeadingInputVariantProps, HeadingInputVariantMap>

export declare const headingInput: HeadingInputRecipe;