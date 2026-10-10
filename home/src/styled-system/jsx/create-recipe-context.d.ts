import type { RecipeDefinition, RecipeRuntimeFn, RecipeSelection, RecipeVariantRecord } from '../types/recipe';
import type { Assign, JsxHTMLProps, JsxStyleProps } from '../types/system';
import type { AsProps, ComponentProps, DataAttrs, JsxFactoryOptions, UnstyledProps, WithForwardedProps } from '../types/jsx';
import type { ElementType, JSX, Provider } from 'react';

type AnyRecipeDefinition = RecipeDefinition<RecipeVariantRecord>

interface RuntimeRecipeFn {
  __type: any
  (props?: any): string
}

type RecipeContextRecipe = RecipeRuntimeFn<any, any> | RuntimeRecipeFn | AnyRecipeDefinition

type RecipePropsOf<R extends RecipeContextRecipe> = R extends RuntimeRecipeFn
  ? R['__type']
  : R extends RecipeRuntimeFn<infer P, any>
    ? P
    : R extends RecipeDefinition<infer T>
      ? RecipeSelection<T>
      : never

type RecipeContextComponentProps<T extends ElementType, R extends RecipeContextRecipe, F extends string> = JsxHTMLProps<
  ComponentProps<T> & UnstyledProps & AsProps,
  WithForwardedProps<Assign<RecipePropsOf<R>, JsxStyleProps>, T, F>
>

type RecipeContextComponent<T extends ElementType, R extends RecipeContextRecipe, F extends string = never> = (
  props: RecipeContextComponentProps<T, R, F>
) => JSX.Element

export interface RecipeContext<R extends RecipeContextRecipe> {
  withContext: <T extends ElementType, F extends string = never>(
    Component: T,
    options?: JsxFactoryOptions<ComponentProps<T>, F> | undefined
  ) => RecipeContextComponent<T, R, F>
  PropsProvider: Provider<Partial<RecipePropsOf<R>> & DataAttrs>
  usePropsContext: () => RecipePropsOf<R> | undefined
}

export declare function createRecipeContext<R extends RecipeContextRecipe>(recipe: R): RecipeContext<R>