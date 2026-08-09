/* eslint-disable */
import type { ConditionalValue } from '../types/index';
import type { DistributiveOmit, Pretty } from '../types/system-types';

interface CardGridVariant {
  
}

type CardGridVariantMap = {
  [key in keyof CardGridVariant]: Array<CardGridVariant[key]>
}

type CardGridSlot = "container" | "grid"

export type CardGridVariantProps = {
  [key in keyof CardGridVariant]?: ConditionalValue<CardGridVariant[key]> | undefined
}

export interface CardGridRecipe {
  __slot: CardGridSlot
  __type: CardGridVariantProps
  (props?: CardGridVariantProps): Pretty<Record<CardGridSlot, string>>
  raw: (props?: CardGridVariantProps) => CardGridVariantProps
  variantMap: CardGridVariantMap
  variantKeys: Array<keyof CardGridVariant>
  splitVariantProps<Props extends CardGridVariantProps>(props: Props): [CardGridVariantProps, Pretty<DistributiveOmit<Props, keyof CardGridVariantProps>>]
  getVariantProps: (props?: CardGridVariantProps) => CardGridVariantProps
}


export declare const cardGrid: CardGridRecipe