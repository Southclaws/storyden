/* eslint-disable */
import type { ConditionalValue } from '../types/index';
import type { DistributiveOmit, Pretty } from '../types/system-types';

interface CardVariant {
  mediaDisplay: "with" | "without"
shape: "box" | "row"
size: "default" | "small"
}

type CardVariantMap = {
  [key in keyof CardVariant]: Array<CardVariant[key]>
}

export type CardVariantProps = {
  [key in keyof CardVariant]?: CardVariant[key] | undefined
}

export interface CardRecipe {
  __type: CardVariantProps
  (props?: CardVariantProps): Pretty<Record<"root" | "mediaBackdropContainer" | "mediaBackdrop" | "contentContainer" | "mediaContainer" | "textArea" | "footer" | "title" | "text" | "media" | "mediaMissing" | "controlsOverlayContainer" | "controls", string>>
  raw: (props?: CardVariantProps) => CardVariantProps
  variantMap: CardVariantMap
  variantKeys: Array<keyof CardVariant>
  splitVariantProps<Props extends CardVariantProps>(props: Props): [CardVariantProps, Pretty<DistributiveOmit<Props, keyof CardVariantProps>>]
}


export declare const card: CardRecipe