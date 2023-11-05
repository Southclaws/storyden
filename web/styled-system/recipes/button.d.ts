/* eslint-disable */
import type { ConditionalValue } from '../types/index';
import type { Pretty } from '../types/helpers';
import type { DistributiveOmit } from '../types/system-types';

interface ButtonVariant {
  kind: "neutral" | "primary" | "secondary" | "destructive" | "blank"
size: "xs" | "sm" | "md" | "lg" | "xl" | "2xl"
}

type ButtonVariantMap = {
  [key in keyof ButtonVariant]: Array<ButtonVariant[key]>
}

export type ButtonVariantProps = {
  [key in keyof ButtonVariant]?: ButtonVariant[key]
}

export interface ButtonRecipe {
  __type: ButtonVariantProps
  (props?: ButtonVariantProps): string
  raw: (props?: ButtonVariantProps) => ButtonVariantProps
  variantMap: ButtonVariantMap
  variantKeys: Array<keyof ButtonVariant>
  splitVariantProps<Props extends ButtonVariantProps>(props: Props): [ButtonVariantProps, Pretty<DistributiveOmit<Props, keyof ButtonVariantProps>>]
}


export declare const button: ButtonRecipe