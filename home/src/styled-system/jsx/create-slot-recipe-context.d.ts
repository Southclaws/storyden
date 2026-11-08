import type { RecipeSelection, SlotRecipeDefinition, SlotRecipeRuntimeFn, SlotRecipeVariantRecord } from '../types/recipe';
import type { Assign, JsxHTMLProps, JsxStyleProps } from '../types/system';
import type { AsProps, ComponentProps, DataAttrs, JsxFactoryOptions, UnstyledProps, WithForwardedProps } from '../types/jsx';
import type { ElementType, JSX, Provider } from 'react';

export type AnySlotRecipeDefinition = SlotRecipeDefinition<string, SlotRecipeVariantRecord<string>>;
export interface RuntimeSlotRecipeFn {
  __type: unknown;
  __slot: string;
  (props?: never): unknown;
}
export type SlotRecipeContextInput = SlotRecipeRuntimeFn<string, any, any, any> | RuntimeSlotRecipeFn | AnySlotRecipeDefinition;
export type SlotNameOf<R extends SlotRecipeContextInput> = R extends RuntimeSlotRecipeFn ? R['__slot'] : R extends SlotRecipeRuntimeFn<infer S, any, any, any> ? S : R extends SlotRecipeDefinition<infer S, any> ? S : string;
export type SlotRecipePropsOf<R extends SlotRecipeContextInput> = R extends RuntimeSlotRecipeFn ? R['__type'] : R extends SlotRecipeRuntimeFn<any, infer P, any, any> ? P : R extends SlotRecipeDefinition<any, infer T> ? RecipeSelection<T> : never;
export interface WithProviderOptions<P = {}> {
  defaultProps?: (Partial<P> & DataAttrs) | undefined;
}
export type SlotRecipeProviderProps<
  T extends ElementType,
  R extends SlotRecipeContextInput,
  F extends string
> = JsxHTMLProps<ComponentProps<T> & UnstyledProps & AsProps, WithForwardedProps<Assign<SlotRecipePropsOf<R>, JsxStyleProps>, T, F>>;
export type SlotRecipeProviderComponent<
  T extends ElementType,
  R extends SlotRecipeContextInput,
  F extends string = never
> = (props: SlotRecipeProviderProps<T, R, F>) => JSX.Element;
export type SlotRecipeRootProviderComponent<
  T extends ElementType,
  R extends SlotRecipeContextInput
> = (props: ComponentProps<T> & UnstyledProps & SlotRecipePropsOf<R>) => JSX.Element;
export type SlotRecipeConsumerComponent<
  T extends ElementType,
  F extends string = never
> = (props: JsxHTMLProps<ComponentProps<T> & UnstyledProps & AsProps, WithForwardedProps<JsxStyleProps, T, F>>) => JSX.Element;
export interface SlotRecipeContext<R extends SlotRecipeContextInput> {
  withRootProvider: <T extends ElementType>(Component: T, options?: WithProviderOptions<ComponentProps<T>> | undefined) => SlotRecipeRootProviderComponent<T, R>;
  withProvider: <
    T extends ElementType,
    F extends string = never
  >(Component: T, slot: SlotNameOf<R>, options?: JsxFactoryOptions<ComponentProps<T>, F> | undefined) => SlotRecipeProviderComponent<T, R, F>;
  withContext: <
    T extends ElementType,
    F extends string = never
  >(Component: T, slot: SlotNameOf<R>, options?: JsxFactoryOptions<ComponentProps<T>, F> | undefined) => SlotRecipeConsumerComponent<T, F>;
  PropsProvider: Provider<Partial<SlotRecipePropsOf<R>> & DataAttrs>;
  usePropsContext: () => SlotRecipePropsOf<R> | undefined;
}
export declare function createSlotRecipeContext<R extends SlotRecipeContextInput>(recipeInput: R): SlotRecipeContext<R>;