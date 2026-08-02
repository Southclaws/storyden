/* eslint-disable */
import type { ConditionalValue } from '../types/index';
import type { DistributiveOmit, Pretty } from '../types/system-types';

interface SectionNavigationVariant {
  
}

type SectionNavigationVariantMap = {
  [key in keyof SectionNavigationVariant]: Array<SectionNavigationVariant[key]>
}

type SectionNavigationSlot = "root" | "trigger" | "triggerLabel" | "triggerIndicator" | "positioner" | "content" | "navigation" | "section" | "sectionLabel" | "items" | "item" | "link" | "linkIndicator"

export type SectionNavigationVariantProps = {
  [key in keyof SectionNavigationVariant]?: ConditionalValue<SectionNavigationVariant[key]> | undefined
}

export interface SectionNavigationRecipe {
  __slot: SectionNavigationSlot
  __type: SectionNavigationVariantProps
  (props?: SectionNavigationVariantProps): Pretty<Record<SectionNavigationSlot, string>>
  raw: (props?: SectionNavigationVariantProps) => SectionNavigationVariantProps
  variantMap: SectionNavigationVariantMap
  variantKeys: Array<keyof SectionNavigationVariant>
  splitVariantProps<Props extends SectionNavigationVariantProps>(props: Props): [SectionNavigationVariantProps, Pretty<DistributiveOmit<Props, keyof SectionNavigationVariantProps>>]
  getVariantProps: (props?: SectionNavigationVariantProps) => SectionNavigationVariantProps
}


export declare const sectionNavigation: SectionNavigationRecipe