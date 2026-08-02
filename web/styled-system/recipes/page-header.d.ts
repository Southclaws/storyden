/* eslint-disable */
import type { ConditionalValue } from '../types/index';
import type { DistributiveOmit, Pretty } from '../types/system-types';

interface PageHeaderVariant {
  
}

type PageHeaderVariantMap = {
  [key in keyof PageHeaderVariant]: Array<PageHeaderVariant[key]>
}

type PageHeaderSlot = "root" | "navigation" | "row" | "heading" | "titleGroup" | "actions"

export type PageHeaderVariantProps = {
  [key in keyof PageHeaderVariant]?: ConditionalValue<PageHeaderVariant[key]> | undefined
}

export interface PageHeaderRecipe {
  __slot: PageHeaderSlot
  __type: PageHeaderVariantProps
  (props?: PageHeaderVariantProps): Pretty<Record<PageHeaderSlot, string>>
  raw: (props?: PageHeaderVariantProps) => PageHeaderVariantProps
  variantMap: PageHeaderVariantMap
  variantKeys: Array<keyof PageHeaderVariant>
  splitVariantProps<Props extends PageHeaderVariantProps>(props: Props): [PageHeaderVariantProps, Pretty<DistributiveOmit<Props, keyof PageHeaderVariantProps>>]
  getVariantProps: (props?: PageHeaderVariantProps) => PageHeaderVariantProps
}


export declare const pageHeader: PageHeaderRecipe