/* eslint-disable */
import type { ElementType, JSX, ComponentPropsWithRef, ComponentType, Component } from 'react'
import type { RecipeDefinition, RecipeSelection, RecipeVariantRecord } from './recipe';
import type { Assign, DistributiveOmit, DistributiveUnion, JsxHTMLProps, JsxStyleProps, Pretty } from './system-types';

interface Dict {
  [k: string]: unknown
}

export interface UnstyledProps {
  /**
   * Whether to remove recipe styles
   */
  unstyled?: boolean | undefined
}

export interface AsProps {
  /**
   * The element to render as
   */
  as?: ElementType | undefined
}

export type ComponentProps<T extends ElementType> = T extends ComponentType<infer P> | Component<infer P>
  ? JSX.LibraryManagedAttributes<T, P>
  : ComponentPropsWithRef<T>

export interface StyledComponent<T extends ElementType, P extends Dict = {}> {
  (props: JsxHTMLProps<ComponentProps<T> & UnstyledProps & AsProps, Assign<JsxStyleProps, P>>): JSX.Element
  displayName?: string | undefined
}

interface RecipeFn {
  __type: any
}

export interface JsxFactoryOptions<TProps extends Dict> {
  dataAttr?: boolean
  defaultProps?: Partial<TProps>
  shouldForwardProp?: (prop: string, variantKeys: string[]) => boolean
  forwardProps?: string[]
}

export type JsxRecipeProps<T extends ElementType, P extends Dict> = JsxHTMLProps<ComponentProps<T> & UnstyledProps & AsProps, P>;

export type JsxElement<T extends ElementType, P extends Dict> = T extends StyledComponent<infer A, infer B>
  ? StyledComponent<A, Pretty<DistributiveUnion<P, B>>>
  : StyledComponent<T, P>

export interface JsxFactory {
  <T extends ElementType>(component: T): StyledComponent<T, {}>
  <T extends ElementType, P extends RecipeVariantRecord>(component: T, recipe: RecipeDefinition<P>, options?: JsxFactoryOptions<JsxRecipeProps<T, RecipeSelection<P>>>): JsxElement<
    T,
    RecipeSelection<P>
  >
  <T extends ElementType, P extends RecipeFn>(component: T, recipeFn: P, options?: JsxFactoryOptions<JsxRecipeProps<T, P['__type']>>): JsxElement<T, P['__type']>
}

export type JsxElements = {
  [K in keyof JSX.IntrinsicElements]: StyledComponent<K, {}>
}

export type Styled = JsxFactory & JsxElements

export type HTMLStyledProps<T extends ElementType> = JsxHTMLProps<ComponentProps<T> & UnstyledProps & AsProps, JsxStyleProps>

export type StyledVariantProps<T extends StyledComponent<any, any>> = T extends StyledComponent<any, infer Props> ? Props : never