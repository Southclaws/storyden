/* eslint-disable */
import type { ConditionalValue } from '../types/index';
import type { Pretty } from '../types/helpers';
import type { DistributiveOmit } from '../types/system-types';

interface SkeletonVariant {
  
}

type SkeletonVariantMap = {
  [key in keyof SkeletonVariant]: Array<SkeletonVariant[key]>
}

export type SkeletonVariantProps = {
  [key in keyof SkeletonVariant]?: ConditionalValue<SkeletonVariant[key]>
}

export interface SkeletonRecipe {
  __type: SkeletonVariantProps
  (props?: SkeletonVariantProps): string
  raw: (props?: SkeletonVariantProps) => SkeletonVariantProps
  variantMap: SkeletonVariantMap
  variantKeys: Array<keyof SkeletonVariant>
  splitVariantProps<Props extends SkeletonVariantProps>(props: Props): [SkeletonVariantProps, Pretty<DistributiveOmit<Props, keyof SkeletonVariantProps>>]
}


export declare const skeleton: SkeletonRecipe