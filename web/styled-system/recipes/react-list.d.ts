/* eslint-disable */
import type { ConditionalValue } from '../types/index';
import type { DistributiveOmit, Pretty } from '../types/system-types';

interface ReactListVariant {
  
}

type ReactListVariantMap = {
  [key in keyof ReactListVariant]: Array<ReactListVariant[key]>
}

type ReactListSlot = "root" | "reaction" | "count" | "picker"

export type ReactListVariantProps = {
  [key in keyof ReactListVariant]?: ConditionalValue<ReactListVariant[key]> | undefined
}

export interface ReactListRecipe {
  __slot: ReactListSlot
  __type: ReactListVariantProps
  (props?: ReactListVariantProps): Pretty<Record<ReactListSlot, string>>
  raw: (props?: ReactListVariantProps) => ReactListVariantProps
  variantMap: ReactListVariantMap
  variantKeys: Array<keyof ReactListVariant>
  splitVariantProps<Props extends ReactListVariantProps>(props: Props): [ReactListVariantProps, Pretty<DistributiveOmit<Props, keyof ReactListVariantProps>>]
  getVariantProps: (props?: ReactListVariantProps) => ReactListVariantProps
}


export declare const reactList: ReactListRecipe