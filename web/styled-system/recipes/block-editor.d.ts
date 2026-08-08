/* eslint-disable */
import type { ConditionalValue } from '../types/index';
import type { DistributiveOmit, Pretty } from '../types/system-types';

interface BlockEditorVariant {
  
}

type BlockEditorVariantMap = {
  [key in keyof BlockEditorVariant]: Array<BlockEditorVariant[key]>
}

type BlockEditorSlot = "root" | "gutter" | "handle" | "menuTrigger" | "menuAnchor" | "menuContent" | "content"

export type BlockEditorVariantProps = {
  [key in keyof BlockEditorVariant]?: ConditionalValue<BlockEditorVariant[key]> | undefined
}

export interface BlockEditorRecipe {
  __slot: BlockEditorSlot
  __type: BlockEditorVariantProps
  (props?: BlockEditorVariantProps): Pretty<Record<BlockEditorSlot, string>>
  raw: (props?: BlockEditorVariantProps) => BlockEditorVariantProps
  variantMap: BlockEditorVariantMap
  variantKeys: Array<keyof BlockEditorVariant>
  splitVariantProps<Props extends BlockEditorVariantProps>(props: Props): [BlockEditorVariantProps, Pretty<DistributiveOmit<Props, keyof BlockEditorVariantProps>>]
  getVariantProps: (props?: BlockEditorVariantProps) => BlockEditorVariantProps
}


export declare const blockEditor: BlockEditorRecipe