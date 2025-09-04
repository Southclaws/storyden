/* eslint-disable */
import type { ConditionalValue } from '../types/index';
import type { DistributiveOmit, Pretty } from '../types/system-types';

interface ClipboardVariant {
  
}

type ClipboardVariantMap = {
  [key in keyof ClipboardVariant]: Array<ClipboardVariant[key]>
}

type ClipboardSlot = "root" | "control" | "trigger" | "indicator" | "input" | "label"

export type ClipboardVariantProps = {
  [key in keyof ClipboardVariant]?: ConditionalValue<ClipboardVariant[key]> | undefined
}

export interface ClipboardRecipe {
  __slot: ClipboardSlot
  __type: ClipboardVariantProps
  (props?: ClipboardVariantProps): Pretty<Record<ClipboardSlot, string>>
  raw: (props?: ClipboardVariantProps) => ClipboardVariantProps
  variantMap: ClipboardVariantMap
  variantKeys: Array<keyof ClipboardVariant>
  splitVariantProps<Props extends ClipboardVariantProps>(props: Props): [ClipboardVariantProps, Pretty<DistributiveOmit<Props, keyof ClipboardVariantProps>>]
  getVariantProps: (props?: ClipboardVariantProps) => ClipboardVariantProps
}


export declare const clipboard: ClipboardRecipe