/* eslint-disable */
import type { ConditionalValue } from '../types/index';
import type { DistributiveOmit, Pretty } from '../types/system-types';

interface TypographyHeadingVariant {
  /**
 * @default "md"
 */
size: "xs" | "sm" | "md" | "lg" | "xl" | "2xl"
}

type TypographyHeadingVariantMap = {
  [key in keyof TypographyHeadingVariant]: Array<TypographyHeadingVariant[key]>
}



export type TypographyHeadingVariantProps = {
  [key in keyof TypographyHeadingVariant]?: ConditionalValue<TypographyHeadingVariant[key]> | undefined
}

export interface TypographyHeadingRecipe {
  
  __type: TypographyHeadingVariantProps
  (props?: TypographyHeadingVariantProps): string
  raw: (props?: TypographyHeadingVariantProps) => TypographyHeadingVariantProps
  variantMap: TypographyHeadingVariantMap
  variantKeys: Array<keyof TypographyHeadingVariant>
  splitVariantProps<Props extends TypographyHeadingVariantProps>(props: Props): [TypographyHeadingVariantProps, Pretty<DistributiveOmit<Props, keyof TypographyHeadingVariantProps>>]
  getVariantProps: (props?: TypographyHeadingVariantProps) => TypographyHeadingVariantProps
}


export declare const typographyHeading: TypographyHeadingRecipe