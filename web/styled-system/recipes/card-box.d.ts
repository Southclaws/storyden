/* eslint-disable */
import type { ConditionalValue } from '../types/index';
import type { DistributiveOmit, Pretty } from '../types/system-types';

interface CardBoxVariant {
  /**
 * @default "default"
 */
kind: "default" | "edge"
}

type CardBoxVariantMap = {
  [key in keyof CardBoxVariant]: Array<CardBoxVariant[key]>
}



export type CardBoxVariantProps = {
  [key in keyof CardBoxVariant]?: ConditionalValue<CardBoxVariant[key]> | undefined
}

export interface CardBoxRecipe {
  
  __type: CardBoxVariantProps
  (props?: CardBoxVariantProps): string
  raw: (props?: CardBoxVariantProps) => CardBoxVariantProps
  variantMap: CardBoxVariantMap
  variantKeys: Array<keyof CardBoxVariant>
  splitVariantProps<Props extends CardBoxVariantProps>(props: Props): [CardBoxVariantProps, Pretty<DistributiveOmit<Props, keyof CardBoxVariantProps>>]
  getVariantProps: (props?: CardBoxVariantProps) => CardBoxVariantProps
}


export declare const cardBox: CardBoxRecipe