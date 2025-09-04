/* eslint-disable */
import type { ConditionalValue } from '../types/index';
import type { DistributiveOmit, Pretty } from '../types/system-types';

interface MenuVariant {
  /**
 * @default "xs"
 */
size: "xs" | "sm" | "md" | "lg"
}

type MenuVariantMap = {
  [key in keyof MenuVariant]: Array<MenuVariant[key]>
}

type MenuSlot = "arrow" | "arrowTip" | "content" | "contextTrigger" | "indicator" | "item" | "itemGroup" | "itemGroupLabel" | "itemIndicator" | "itemText" | "positioner" | "separator" | "trigger" | "triggerItem"

export type MenuVariantProps = {
  [key in keyof MenuVariant]?: ConditionalValue<MenuVariant[key]> | undefined
}

export interface MenuRecipe {
  __slot: MenuSlot
  __type: MenuVariantProps
  (props?: MenuVariantProps): Pretty<Record<MenuSlot, string>>
  raw: (props?: MenuVariantProps) => MenuVariantProps
  variantMap: MenuVariantMap
  variantKeys: Array<keyof MenuVariant>
  splitVariantProps<Props extends MenuVariantProps>(props: Props): [MenuVariantProps, Pretty<DistributiveOmit<Props, keyof MenuVariantProps>>]
  getVariantProps: (props?: MenuVariantProps) => MenuVariantProps
}


export declare const menu: MenuRecipe