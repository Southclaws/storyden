/* eslint-disable */
import type { ConditionalValue } from '../types/index';
import type { DistributiveOmit, Pretty } from '../types/system-types';

interface GroupVariant {
  /**
 * @default "horizontal"
 */
orientation: "horizontal" | "vertical"
attached: boolean
grow: boolean
}

type GroupVariantMap = {
  [key in keyof GroupVariant]: Array<GroupVariant[key]>
}



export type GroupVariantProps = {
  [key in keyof GroupVariant]?: GroupVariant[key] | undefined
}

export interface GroupRecipe {
  
  __type: GroupVariantProps
  (props?: GroupVariantProps): string
  raw: (props?: GroupVariantProps) => GroupVariantProps
  variantMap: GroupVariantMap
  variantKeys: Array<keyof GroupVariant>
  splitVariantProps<Props extends GroupVariantProps>(props: Props): [GroupVariantProps, Pretty<DistributiveOmit<Props, keyof GroupVariantProps>>]
  getVariantProps: (props?: GroupVariantProps) => GroupVariantProps
}


export declare const group: GroupRecipe