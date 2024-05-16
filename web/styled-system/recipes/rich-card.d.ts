/* eslint-disable */
import type { ConditionalValue } from '../types/index';
import type { DistributiveOmit, Pretty } from '../types/system-types';

interface RichCardVariant {
  mediaDisplay: "with" | "without"
shape: "box" | "row"
size: "default" | "small"
}

type RichCardVariantMap = {
  [key in keyof RichCardVariant]: Array<RichCardVariant[key]>
}

export type RichCardVariantProps = {
  [key in keyof RichCardVariant]?: RichCardVariant[key] | undefined
}

export interface RichCardRecipe {
  __type: RichCardVariantProps
  (props?: RichCardVariantProps): Pretty<Record<"root" | "mediaBackdropContainer" | "mediaBackdrop" | "contentContainer" | "mediaContainer" | "textArea" | "footer" | "title" | "text" | "media" | "mediaMissing" | "controlsOverlayContainer" | "controls", string>>
  raw: (props?: RichCardVariantProps) => RichCardVariantProps
  variantMap: RichCardVariantMap
  variantKeys: Array<keyof RichCardVariant>
  splitVariantProps<Props extends RichCardVariantProps>(props: Props): [RichCardVariantProps, Pretty<DistributiveOmit<Props, keyof RichCardVariantProps>>]
  getVariantProps: (props?: RichCardVariantProps) => RichCardVariantProps
}


export declare const richCard: RichCardRecipe