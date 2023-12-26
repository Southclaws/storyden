/* eslint-disable */
import type { ConditionalValue } from '../types/index';
import type { Pretty } from '../types/helpers';
import type { DistributiveOmit } from '../types/system-types';

interface TitleInputVariant {
  
}

type TitleInputVariantMap = {
  [key in keyof TitleInputVariant]: Array<TitleInputVariant[key]>
}

export type TitleInputVariantProps = {
  [key in keyof TitleInputVariant]?: ConditionalValue<TitleInputVariant[key]>
}

export interface TitleInputRecipe {
  __type: TitleInputVariantProps
  (props?: TitleInputVariantProps): string
  raw: (props?: TitleInputVariantProps) => TitleInputVariantProps
  variantMap: TitleInputVariantMap
  variantKeys: Array<keyof TitleInputVariant>
  splitVariantProps<Props extends TitleInputVariantProps>(props: Props): [TitleInputVariantProps, Pretty<DistributiveOmit<Props, keyof TitleInputVariantProps>>]
}


export declare const titleInput: TitleInputRecipe