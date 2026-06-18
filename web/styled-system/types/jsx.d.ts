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

type BaseComponentProps<T extends ElementType> = ComponentProps<T> & UnstyledProps & AsProps & DataAttrs

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

export interface JsxFactoryOptions<TProps extends AnyProps> {
  dataAttr?: boolean
  defaultProps?: Partial<TProps> & DataAttrs
  shouldForwardProp?: (prop: string, variantKeys: string[]) => boolean
  forwardProps?: string[]
}

export type JsxRecipeProps<T extends ElementType, P extends AnyProps> = JsxHTMLProps<BaseComponentProps<T>, P>

export type JsxElement<T extends ElementType, P extends AnyProps> = T extends StyledComponent<infer A, infer B>
  ? StyledComponent<A, Assign<B, P>>
  : StyledComponent<T, P>

export interface JsxFactory {
  <T extends ElementType>(component: T): StyledComponent<T, {}>
  <T extends ElementType, P extends RecipeVariantRecord = {}>(component: T, recipe: RecipeDefinition<P>, options?: JsxFactoryOptions<JsxRecipeProps<T, RecipeSelection<P>>>): JsxElement<T, RecipeSelection<P>>
  <T extends ElementType, P extends RuntimeRecipeFn>(component: T, recipeFn: P, options?: JsxFactoryOptions<JsxRecipeProps<T, P["__type"]>>): JsxElement<T, P["__type"]>
}

export type JsxElements = {
  [K in keyof JSX.IntrinsicElements]: StyledComponent<K, {}>
}

export type Styled = JsxFactory & JsxElements

export type HTMLStyledProps<T extends ElementType> = JsxHTMLProps<BaseComponentProps<T>, JsxStyleProps>

export type StyledVariantProps<T extends StyledComponent<any, any>> = T extends StyledComponent<any, infer Props> ? Props : never