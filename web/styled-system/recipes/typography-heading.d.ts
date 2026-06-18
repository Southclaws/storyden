import type { ConditionalValue } from '../types/system';
import type { RecipeRuntimeFn, RecipeVariantMap } from '../types/recipe';

export type TypographyHeadingVariant = {
  size?: "2xl" | "lg" | "md" | "sm" | "xl" | "xs"
}

export type TypographyHeadingVariantProps = {
  [K in keyof TypographyHeadingVariant]?: ConditionalValue<TypographyHeadingVariant[K]>
}

export type TypographyHeadingVariantMap = RecipeVariantMap<TypographyHeadingVariant>

export type TypographyHeadingRecipe = RecipeRuntimeFn<TypographyHeadingVariantProps, TypographyHeadingVariantMap>

export declare const typographyHeading: TypographyHeadingRecipe;