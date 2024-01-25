/* eslint-disable */
import type { ConditionalValue } from '../types/index';
import type { Pretty } from '../types/helpers';
import type { DistributiveOmit } from '../types/system-types';

interface HeadingInputVariant {
  
}

type HeadingInputVariantMap = {
  [key in keyof HeadingInputVariant]: Array<HeadingInputVariant[key]>
}

export type HeadingInputVariantProps = {
  [key in keyof HeadingInputVariant]?: ConditionalValue<HeadingInputVariant[key]>
}

export interface HeadingInputRecipe {
  __type: HeadingInputVariantProps
  (props?: HeadingInputVariantProps): string
  raw: (props?: HeadingInputVariantProps) => HeadingInputVariantProps
  variantMap: HeadingInputVariantMap
  variantKeys: Array<keyof HeadingInputVariant>
  splitVariantProps<Props extends HeadingInputVariantProps>(props: Props): [HeadingInputVariantProps, Pretty<DistributiveOmit<Props, keyof HeadingInputVariantProps>>]
}


export declare const headingInput: HeadingInputRecipe