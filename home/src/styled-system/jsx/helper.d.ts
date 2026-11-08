import type { RecipeCreatorFn } from '../types/recipe';
import type { SystemStyleObject } from '../types/system';

export type Props = Record<string, unknown>;
export type ShouldForwardProp = (prop: string, variantKeys?: string[]) => boolean;
/** A `cva` or config recipe function, as the JSX runtime reads it. */
export interface JsxRecipeFn {
  (props?: Props, withCompoundVariants?: boolean): string;
  raw(props?: Props): SystemStyleObject;
  merge(other: JsxRecipeFn, cvaFn: RecipeCreatorFn): JsxRecipeFn;
  variantKeys: string[];
  __cva__?: boolean;
  __recipe__?: boolean;
  __name__?: string;
  __memoizedRaw__?: (props?: Props) => SystemStyleObject;
  __getCompoundVariantClasses__?: (props: Props) => string;
}
/** Markers the factory sets on the components it creates. */
export interface StyledMarkers {
  __cva__?: JsxRecipeFn;
  __base__?: unknown;
  __shouldForwardProps__?: ShouldForwardProp;
}
export declare const composeShouldForwardProps: (tag: StyledMarkers, shouldForwardProp: ShouldForwardProp) => ShouldForwardProp;
export declare const composeCvaFn: (cvaA: JsxRecipeFn | undefined, cvaB: JsxRecipeFn) => JsxRecipeFn;
export declare function mergeDefaultProps(defaultProps: Props, props: Props): Props;
export declare const getDisplayName: (Component: string | {
  displayName?: string;
  name?: string;
} | undefined) => string;
export type SplitJsxProps = [htmlProps: Props | undefined, forwardedProps: Props | undefined, variantProps: Props, propStyles: SystemStyleObject | undefined, cssStyles: SystemStyleObject | undefined, elementProps: Props | undefined];
export declare function splitJsxProps(props: Props, shouldForwardProp: (prop: string) => boolean, variantSet: Set<string>, isCssProperty: (prop: string) => boolean, skipClass?: boolean): SplitJsxProps;
export declare function serializeSplitStyles(propStyles: SystemStyleObject | undefined, cssStyles: SystemStyleObject | undefined, baseStyles?: SystemStyleObject): string;
export declare function splitStyleProps(styleProps: Props): [propStyles: SystemStyleObject | undefined, cssStyles: SystemStyleObject | undefined];