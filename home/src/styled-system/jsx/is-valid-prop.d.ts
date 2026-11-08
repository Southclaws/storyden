import type { DistributiveOmit, JsxStyleProps } from '../types/system';

export type CssPropKey = keyof JsxStyleProps;
export type OmittedCssProps<T> = DistributiveOmit<T, CssPropKey>;
export declare function isCssProperty(value: string): boolean;
export declare function splitCssProps<T>(props: T): [JsxStyleProps, OmittedCssProps<T>];