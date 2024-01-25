/* eslint-disable */
import type { ConditionalValue } from '../types/index';
import type { Pretty } from '../types/helpers';
import type { DistributiveOmit } from '../types/system-types';

interface CardVariant {
  shape: "box" | "row"
}

type CardVariantMap = {
  [key in keyof CardVariant]: Array<CardVariant[key]>
}

export type CardVariantProps = {
  [key in keyof CardVariant]?: ConditionalValue<CardVariant[key]>
}

export interface CardRecipe {
  __type: CardVariantProps
  (props?: CardVariantProps): Pretty<Record<"root" | "textArea" | "title" | "text" | "media" | "mediaContainer" | "mediaBackdrop" | "mediaBackdropContainer", string>>
  raw: (props?: CardVariantProps) => CardVariantProps
  variantMap: CardVariantMap
  variantKeys: Array<keyof CardVariant>
  splitVariantProps<Props extends CardVariantProps>(props: Props): [CardVariantProps, Pretty<DistributiveOmit<Props, keyof CardVariantProps>>]
}


export declare const card: CardRecipe