/* eslint-disable */
import type { ConditionalValue } from '../types/index';
import type { Pretty } from '../types/helpers';
import type { DistributiveOmit } from '../types/system-types';

interface AdmonitionVariant {
  kind: "neutral" | "success" | "failure"
}

type AdmonitionVariantMap = {
  [key in keyof AdmonitionVariant]: Array<AdmonitionVariant[key]>
}

export type AdmonitionVariantProps = {
  [key in keyof AdmonitionVariant]?: ConditionalValue<AdmonitionVariant[key]>
}

export interface AdmonitionRecipe {
  __type: AdmonitionVariantProps
  (props?: AdmonitionVariantProps): string
  raw: (props?: AdmonitionVariantProps) => AdmonitionVariantProps
  variantMap: AdmonitionVariantMap
  variantKeys: Array<keyof AdmonitionVariant>
  splitVariantProps<Props extends AdmonitionVariantProps>(props: Props): [AdmonitionVariantProps, Pretty<DistributiveOmit<Props, keyof AdmonitionVariantProps>>]
}


export declare const admonition: AdmonitionRecipe