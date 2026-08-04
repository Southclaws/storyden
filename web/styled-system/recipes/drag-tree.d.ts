/* eslint-disable */
import type { ConditionalValue } from '../types/index';
import type { DistributiveOmit, Pretty } from '../types/system-types';

interface DragTreeVariant {
  /**
 * @default "clamped"
 */
variant: "clamped" | "scrollable"
}

type DragTreeVariantMap = {
  [key in keyof DragTreeVariant]: Array<DragTreeVariant[key]>
}

type DragTreeSlot = "branch" | "branchContent" | "branchControl" | "branchIndentGuide" | "branchIndicator" | "branchText" | "branchTrigger" | "item" | "itemIndicator" | "itemText" | "label" | "nodeCheckbox" | "nodeRenameInput" | "root" | "tree" | "actions" | "link" | "row" | "dropIndicator" | "dropIndicatorLine"

export type DragTreeVariantProps = {
  [key in keyof DragTreeVariant]?: ConditionalValue<DragTreeVariant[key]> | undefined
}

export interface DragTreeRecipe {
  __slot: DragTreeSlot
  __type: DragTreeVariantProps
  (props?: DragTreeVariantProps): Pretty<Record<DragTreeSlot, string>>
  raw: (props?: DragTreeVariantProps) => DragTreeVariantProps
  variantMap: DragTreeVariantMap
  variantKeys: Array<keyof DragTreeVariant>
  splitVariantProps<Props extends DragTreeVariantProps>(props: Props): [DragTreeVariantProps, Pretty<DistributiveOmit<Props, keyof DragTreeVariantProps>>]
  getVariantProps: (props?: DragTreeVariantProps) => DragTreeVariantProps
}


export declare const dragTree: DragTreeRecipe