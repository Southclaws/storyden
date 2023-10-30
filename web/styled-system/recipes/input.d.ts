/* eslint-disable */
import type { ConditionalValue } from '../types/index';
import type { Pretty } from '../types/helpers';
import type { DistributiveOmit } from '../types/system-types';

interface InputVariant {
  size: "2xs" | "xs" | "sm" | "md" | "lg" | "xl" | "2xl"
}

type InputVariantMap = {
  [key in keyof InputVariant]: Array<InputVariant[key]>
}

export type InputVariantProps = {
  [key in keyof InputVariant]?: ConditionalValue<InputVariant[key]>
}

export interface InputRecipe {
  __type: InputVariantProps
  (props?: InputVariantProps): string
  raw: (props?: InputVariantProps) => InputVariantProps
  variantMap: InputVariantMap
  variantKeys: Array<keyof InputVariant>
  splitVariantProps<Props extends InputVariantProps>(props: Props): [InputVariantProps, Pretty<DistributiveOmit<Props, keyof InputVariantProps>>]
}


export declare const input: InputRecipe