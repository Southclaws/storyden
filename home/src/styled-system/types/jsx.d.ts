import type { ElementType, JSX } from 'react';
import type { RecipeDefinition, RecipeSelection, RecipeVariantRecord } from './recipe';
import type { Assign, JsxHTMLProps, JsxStyleProps } from './system';

interface AnyProps {
  [k: string]: unknown
}

export type DataAttrs = Record<`data-${string}`, unknown>

export interface UnstyledProps {
  unstyled?: boolean | undefined
}

export interface AsProps {
  as?: ElementType | undefined
}

export type ComponentProps<T extends ElementType> = T extends keyof JSX.IntrinsicElements
  ? JSX.IntrinsicElements[T]
  : T extends { (props: infer Props): any }
    ? Props
    : T extends abstract new (props: infer Props) => any
      ? Props
      : {}

type BaseComponentProps<T extends ElementType> = ComponentProps<T> & UnstyledProps & AsProps

export type StyledComponentProps<T extends ElementType, P extends AnyProps = {}> = JsxHTMLProps<
  BaseComponentProps<T>,
  Assign<JsxStyleProps, P>
>

export interface StyledComponent<T extends ElementType, P extends AnyProps = {}> {
  (props: StyledComponentProps<T, P>): JSX.Element
  displayName?: string | undefined
}

interface RuntimeRecipeFn {
  __type: any
}

export interface JsxFactoryOptions<TProps extends AnyProps, F extends string = string> {
  dataAttr?: boolean
  defaultProps?: Partial<TProps> & DataAttrs
  shouldForwardProp?: (prop: string, variantKeys: string[]) => boolean
  forwardProps?: readonly F[]
}

// Distributes over the few forwarded keys; never intersects the large key unions.
type ForwardedStyleKeys<T extends ElementType, F extends string> = F extends keyof JsxStyleProps
  ? F extends keyof ComponentProps<T> ? F : never
  : never

/** Props `S`, with each forwarded prop that shares a style prop's name typed from the component. */
export type WithForwardedProps<S, T extends ElementType, F extends string> = [ForwardedStyleKeys<T, F>] extends [never]
  ? S
  : Assign<S, Pick<ComponentProps<T>, ForwardedStyleKeys<T, F>>>

export type JsxRecipeProps<T extends ElementType, P extends AnyProps> = JsxHTMLProps<BaseComponentProps<T>, P>

export type JsxElement<T extends ElementType, P extends AnyProps> = T extends StyledComponent<infer A, infer B>
  ? StyledComponent<A, Assign<B, P>>
  : StyledComponent<T, P>

export interface JsxFactory {
  <T extends ElementType>(component: T): StyledComponent<T, {}>
  <T extends ElementType, P extends RecipeVariantRecord = {}, F extends string = never>(component: T, recipe: RecipeDefinition<P>, options?: JsxFactoryOptions<JsxRecipeProps<T, RecipeSelection<P>>, F>): JsxElement<T, WithForwardedProps<RecipeSelection<P>, T, F>>
  <T extends ElementType, P extends RuntimeRecipeFn, F extends string = never>(component: T, recipeFn: P, options?: JsxFactoryOptions<JsxRecipeProps<T, P["__type"]>, F>): JsxElement<T, WithForwardedProps<P["__type"], T, F>>
}

export type JsxElements = {
  [K in keyof JSX.IntrinsicElements]: StyledComponent<K, {}>
}

export type Styled = JsxFactory & JsxElements

export type HTMLStyledProps<T extends ElementType> = JsxHTMLProps<BaseComponentProps<T>, JsxStyleProps>

export type StyledVariantProps<T extends StyledComponent<any, any>> = T extends StyledComponent<any, infer Props> ? Props : never