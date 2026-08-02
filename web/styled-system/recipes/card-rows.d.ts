/* eslint-disable */
import type { ConditionalValue } from '../types/index';
import type { DistributiveOmit, Pretty } from '../types/system-types';

interface CardRowsVariant {
  
}

type CardRowsVariantMap = {
  [key in keyof CardRowsVariant]: Array<CardRowsVariant[key]>
}



export type CardRowsVariantProps = {
  [key in keyof CardRowsVariant]?: ConditionalValue<CardRowsVariant[key]> | undefined
}

export interface CardRowsRecipe {
  
  __type: CardRowsVariantProps
  (props?: CardRowsVariantProps): string
  raw: (props?: CardRowsVariantProps) => CardRowsVariantProps
  variantMap: CardRowsVariantMap
  variantKeys: Array<keyof CardRowsVariant>
  splitVariantProps<Props extends CardRowsVariantProps>(props: Props): [CardRowsVariantProps, Pretty<DistributiveOmit<Props, keyof CardRowsVariantProps>>]
  getVariantProps: (props?: CardRowsVariantProps) => CardRowsVariantProps
}


export declare const cardRows: CardRowsRecipe