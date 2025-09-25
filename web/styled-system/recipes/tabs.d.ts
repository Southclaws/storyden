/* eslint-disable */
import type { ConditionalValue } from '../types/index';
import type { DistributiveOmit, Pretty } from '../types/system-types';

interface TabsVariant {
  /**
 * @default "line"
 */
variant: "enclosed" | "line" | "outline"
/**
 * @default "md"
 */
size: "sm" | "md" | "lg"
}

type TabsVariantMap = {
  [key in keyof TabsVariant]: Array<TabsVariant[key]>
}

type TabsSlot = "root" | "list" | "trigger" | "content" | "indicator"

export type TabsVariantProps = {
  [key in keyof TabsVariant]?: TabsVariant[key] | undefined
}

export interface TabsRecipe {
  __slot: TabsSlot
  __type: TabsVariantProps
  (props?: TabsVariantProps): Pretty<Record<TabsSlot, string>>
  raw: (props?: TabsVariantProps) => TabsVariantProps
  variantMap: TabsVariantMap
  variantKeys: Array<keyof TabsVariant>
  splitVariantProps<Props extends TabsVariantProps>(props: Props): [TabsVariantProps, Pretty<DistributiveOmit<Props, keyof TabsVariantProps>>]
  getVariantProps: (props?: TabsVariantProps) => TabsVariantProps
}


export declare const tabs: TabsRecipe