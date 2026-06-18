import type { ConditionalValue } from '../types/system';
import type { SlotRecipeRuntimeFn, RecipeVariantMap } from '../types/recipe';

export type ColorPickerVariant = {}

export type ColorPickerVariantProps = {
  [K in keyof ColorPickerVariant]?: ConditionalValue<ColorPickerVariant[K]>
}

export type ColorPickerVariantMap = RecipeVariantMap<ColorPickerVariant>

export type ColorPickerSlot = "root" | "label" | "control" | "trigger" | "positioner" | "content" | "area" | "areaThumb" | "valueText" | "areaBackground" | "channelSlider" | "channelSliderLabel" | "channelSliderTrack" | "channelSliderThumb" | "channelSliderValueText" | "channelInput" | "transparencyGrid" | "swatchGroup" | "swatchTrigger" | "swatchIndicator" | "swatch" | "eyeDropperTrigger" | "formatTrigger" | "formatSelect" | "view"

export type ColorPickerRecipe = SlotRecipeRuntimeFn<ColorPickerSlot, ColorPickerVariantProps, ColorPickerVariantMap>

export declare const colorPicker: ColorPickerRecipe;