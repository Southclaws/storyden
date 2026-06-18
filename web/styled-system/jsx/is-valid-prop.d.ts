import type { DistributiveOmit, JsxStyleProps } from '../types/system';

declare const isCssProperty: (value: string) => boolean

type CssPropKey = keyof JsxStyleProps
type OmittedCssProps<T> = DistributiveOmit<T, CssPropKey>

declare const splitCssProps: <T>(props: T) => [JsxStyleProps, OmittedCssProps<T>]

export { isCssProperty, splitCssProps }