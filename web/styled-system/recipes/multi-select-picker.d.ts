/* eslint-disable */
import type { ConditionalValue } from '../types/index';
import type { DistributiveOmit, Pretty } from '../types/system-types';

interface MultiSelectPickerVariant {
  /**
 * @default "sm"
 */
size: "sm" | "md" | "lg"
}

type MultiSelectPickerVariantMap = {
  [key in keyof MultiSelectPickerVariant]: Array<MultiSelectPickerVariant[key]>
}



export type MultiSelectPickerVariantProps = {
  [key in keyof MultiSelectPickerVariant]?: ConditionalValue<MultiSelectPickerVariant[key]> | undefined
}

export interface MultiSelectPickerRecipe {
  
  __type: MultiSelectPickerVariantProps
  (props?: MultiSelectPickerVariantProps): string
  raw: (props?: MultiSelectPickerVariantProps) => MultiSelectPickerVariantProps
  variantMap: MultiSelectPickerVariantMap
  variantKeys: Array<keyof MultiSelectPickerVariant>
  splitVariantProps<Props extends MultiSelectPickerVariantProps>(props: Props): [MultiSelectPickerVariantProps, Pretty<DistributiveOmit<Props, keyof MultiSelectPickerVariantProps>>]
  getVariantProps: (props?: MultiSelectPickerVariantProps) => MultiSelectPickerVariantProps
}


export declare const multiSelectPicker: MultiSelectPickerRecipe