/* eslint-disable */
import type { ConditionalValue } from '../types/index';
import type { DistributiveOmit, Pretty } from '../types/system-types';

interface AdmonitionVariant {
  kind: "neutral" | "success" | "failure"
}

type AdmonitionVariantMap = {
  [key in keyof AdmonitionVariant]: Array<AdmonitionVariant[key]>
}



export type AdmonitionVariantProps = {
  [key in keyof AdmonitionVariant]?: ConditionalValue<AdmonitionVariant[key]> | undefined
}

export interface AdmonitionRecipe {
  
  __type: AdmonitionVariantProps
  (props?: AdmonitionVariantProps): string
  raw: (props?: AdmonitionVariantProps) => AdmonitionVariantProps
  variantMap: AdmonitionVariantMap
  variantKeys: Array<keyof AdmonitionVariant>
  splitVariantProps<Props extends AdmonitionVariantProps>(props: Props): [AdmonitionVariantProps, Pretty<DistributiveOmit<Props, keyof AdmonitionVariantProps>>]
  getVariantProps: (props?: AdmonitionVariantProps) => AdmonitionVariantProps
}


export declare const admonition: AdmonitionRecipe