import type { RecipeSelection, SlotRecipeDefinition, SlotRecipeRuntimeFn, SlotRecipeVariantRecord } from '../types/recipe';
import type { Assign, JsxHTMLProps, JsxStyleProps } from '../types/system';
import type { AsProps, ComponentProps, DataAttrs, JsxFactoryOptions } from '../types/jsx';
import type { ElementType, JSX } from 'react';

interface UnstyledProps {
  unstyled?: boolean | undefined
}

type AnySlotRecipeDefinition = SlotRecipeDefinition<string, SlotRecipeVariantRecord<string>>

interface RuntimeSlotRecipeFn {
  __type: any
  __slot: string
  (props?: any): any
}

type SlotRecipeContextInput = SlotRecipeRuntimeFn<string, any, any> | RuntimeSlotRecipeFn | AnySlotRecipeDefinition

type SlotNameOf<R extends SlotRecipeContextInput> = R extends RuntimeSlotRecipeFn
  ? R['__slot']
  : R extends SlotRecipeRuntimeFn<infer S, any, any>
    ? S
    : R extends SlotRecipeDefinition<infer S, any>
      ? S
      : string

type SlotRecipePropsOf<R extends SlotRecipeContextInput> = R extends RuntimeSlotRecipeFn
  ? R['__type']
  : R extends SlotRecipeRuntimeFn<any, infer P, any>
    ? P
    : R extends SlotRecipeDefinition<any, infer T>
      ? RecipeSelection<T>
      : never

interface WithProviderOptions<P = {}> {
  defaultProps?: (Partial<P> & DataAttrs) | undefined
}

type SlotRecipeProviderProps<T extends ElementType, R extends SlotRecipeContextInput> = JsxHTMLProps<
  ComponentProps<T> & UnstyledProps & AsProps & DataAttrs,
  Assign<SlotRecipePropsOf<R>, JsxStyleProps>
>

type SlotRecipeProviderComponent<T extends ElementType, R extends SlotRecipeContextInput> = (
  props: SlotRecipeProviderProps<T, R>
) => JSX.Element

type SlotRecipeRootProviderComponent<T extends ElementType, R extends SlotRecipeContextInput> = (
  props: ComponentProps<T> & UnstyledProps & DataAttrs & SlotRecipePropsOf<R>
) => JSX.Element

type SlotRecipeConsumerComponent<T extends ElementType> = (
  props: JsxHTMLProps<ComponentProps<T> & UnstyledProps & AsProps & DataAttrs, JsxStyleProps>
) => JSX.Element

export interface SlotRecipeContext<R extends SlotRecipeContextInput> {
  withRootProvider: <T extends ElementType>(
    Component: T,
    options?: WithProviderOptions<ComponentProps<T>> | undefined
  ) => SlotRecipeRootProviderComponent<T, R>
  withProvider: <T extends ElementType>(
    Component: T,
    slot: SlotNameOf<R>,
    options?: JsxFactoryOptions<ComponentProps<T>> | undefined
  ) => SlotRecipeProviderComponent<T, R>
  withContext: <T extends ElementType>(
    Component: T,
    slot: SlotNameOf<R>,
    options?: JsxFactoryOptions<ComponentProps<T>> | undefined
  ) => SlotRecipeConsumerComponent<T>
}

export declare function createSlotRecipeContext<R extends SlotRecipeContextInput>(recipe: R): SlotRecipeContext<R>