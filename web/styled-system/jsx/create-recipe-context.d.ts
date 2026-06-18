import type { RecipeDefinition, RecipeRuntimeFn, RecipeSelection, RecipeVariantRecord } from '../types/recipe';
import type { Assign, JsxHTMLProps, JsxStyleProps } from '../types/system';
import type { AsProps, ComponentProps, DataAttrs, JsxFactoryOptions } from '../types/jsx';
import type { ElementType, JSX, Provider } from 'react';

interface UnstyledProps {
  unstyled?: boolean | undefined
}

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

type RecipeContextComponentProps<T extends ElementType, R extends RecipeContextRecipe> = JsxHTMLProps<
  ComponentProps<T> & UnstyledProps & AsProps & DataAttrs,
  Assign<RecipePropsOf<R>, JsxStyleProps>
>

type RecipeContextComponent<T extends ElementType, R extends RecipeContextRecipe> = (
  props: RecipeContextComponentProps<T, R>
) => JSX.Element

export interface RecipeContext<R extends RecipeContextRecipe> {
  withContext: <T extends ElementType>(
    Component: T,
    options?: JsxFactoryOptions<ComponentProps<T>> | undefined
  ) => RecipeContextComponent<T, R>
  PropsProvider: Provider<Partial<RecipePropsOf<R>> & DataAttrs>
  usePropsContext: () => RecipePropsOf<R> | undefined
}

export declare function createRecipeContext<R extends RecipeContextRecipe>(recipe: R): RecipeContext<R>