import type { ConditionalValue } from '../types/system';
import type { RecipeRuntimeFn, RecipeVariantMap } from '../types/recipe';

export type HeadingInputVariant = {}

export type HeadingInputVariantProps = {
  [K in keyof HeadingInputVariant]?: ConditionalValue<HeadingInputVariant[K]>
}

export type HeadingInputVariantMap = RecipeVariantMap<HeadingInputVariant>

export type HeadingInputRecipe = RecipeRuntimeFn<HeadingInputVariantProps, HeadingInputVariantMap>

export declare const headingInput: HeadingInputRecipe;