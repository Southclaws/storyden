/* eslint-disable */
import type { ConditionalValue } from '../types/index';
import type { DistributiveOmit, Pretty } from '../types/system-types';

interface RichCardVariant {
  /**
 * @default "row"
 */
shape: "row" | "responsive" | "box" | "fill"
}

type RichCardVariantMap = {
  [key in keyof RichCardVariant]: Array<RichCardVariant[key]>
}

export type RichCardVariantProps = {
  [key in keyof RichCardVariant]?: ConditionalValue<RichCardVariant[key]> | undefined
}

export interface RichCardRecipe {
  __type: RichCardVariantProps
  (props?: RichCardVariantProps): Pretty<Record<"root" | "headerContainer" | "menuContainer" | "titleContainer" | "contentContainer" | "mediaContainer" | "footerContainer" | "mediaBackdropContainer" | "mediaBackdrop" | "textArea" | "text" | "media" | "mediaMissing", string>>
  raw: (props?: RichCardVariantProps) => RichCardVariantProps
  variantMap: RichCardVariantMap
  variantKeys: Array<keyof RichCardVariant>
  splitVariantProps<Props extends RichCardVariantProps>(props: Props): [RichCardVariantProps, Pretty<DistributiveOmit<Props, keyof RichCardVariantProps>>]
  getVariantProps: (props?: RichCardVariantProps) => RichCardVariantProps
}


export declare const richCard: RichCardRecipe