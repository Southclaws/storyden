/* eslint-disable */
import type { ConditionalValue } from '../types/index';
import type { DistributiveOmit, Pretty } from '../types/system-types';

interface ColorPickerVariant {
  
}

type ColorPickerVariantMap = {
  [key in keyof ColorPickerVariant]: Array<ColorPickerVariant[key]>
}

export type ColorPickerVariantProps = {
  [key in keyof ColorPickerVariant]?: ConditionalValue<ColorPickerVariant[key]> | undefined
}

export interface ColorPickerRecipe {
  __type: ColorPickerVariantProps
  (props?: ColorPickerVariantProps): Pretty<Record<"root" | "label" | "control" | "trigger" | "positioner" | "content" | "area" | "areaThumb" | "areaBackground" | "channelSlider" | "channelSliderTrack" | "channelSliderThumb" | "channelInput" | "transparencyGrid" | "swatchGroup" | "swatchTrigger" | "swatchIndicator" | "swatch" | "eyeDropperTrigger" | "formatTrigger" | "formatSelect" | "valueText" | "view", string>>
  raw: (props?: ColorPickerVariantProps) => ColorPickerVariantProps
  variantMap: ColorPickerVariantMap
  variantKeys: Array<keyof ColorPickerVariant>
  splitVariantProps<Props extends ColorPickerVariantProps>(props: Props): [ColorPickerVariantProps, Pretty<DistributiveOmit<Props, keyof ColorPickerVariantProps>>]
  getVariantProps: (props?: ColorPickerVariantProps) => ColorPickerVariantProps
}


export declare const colorPicker: ColorPickerRecipe