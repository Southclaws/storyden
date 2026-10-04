import type { ConditionalValue } from '../types/system';
import type { RecipeRuntimeFn, RecipeVariantMap } from '../types/recipe';

export type MultiSelectPickerVariant = {
  size?: "lg" | "md" | "sm"
}

export type MultiSelectPickerVariantProps = {
  [K in keyof MultiSelectPickerVariant]?: ConditionalValue<MultiSelectPickerVariant[K]>
}

export type MultiSelectPickerVariantMap = RecipeVariantMap<MultiSelectPickerVariant>

export type MultiSelectPickerRecipe = RecipeRuntimeFn<MultiSelectPickerVariantProps, MultiSelectPickerVariantMap>

export declare const multiSelectPicker: MultiSelectPickerRecipe;