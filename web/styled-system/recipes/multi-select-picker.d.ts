import type { ConditionalValue } from '../types/system';
import type { ConfigRecipeRuntimeFn, RecipeVariantMap } from '../types/recipe';

export type MultiSelectPickerVariant = {
  size?: "sm" | "md" | "lg"
}

export type MultiSelectPickerVariantProps = {
  [K in keyof MultiSelectPickerVariant]?: ConditionalValue<MultiSelectPickerVariant[K]>
}

export type MultiSelectPickerVariantMap = RecipeVariantMap<MultiSelectPickerVariant>

export type MultiSelectPickerRecipe = ConfigRecipeRuntimeFn<MultiSelectPickerVariantProps, MultiSelectPickerVariantMap>

export declare const multiSelectPicker: MultiSelectPickerRecipe;