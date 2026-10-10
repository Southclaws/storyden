"use client";

import { createContext, useContext, createElement, forwardRef } from 'react';
import { styled } from './factory';
import { getDisplayName } from './helper';

function resolveRecipe(recipe) {
  if (recipe == null) throw new Error('createRecipeContext requires a recipe')
  if (recipe.__recipe__ === true || recipe.__cva__ === true) return recipe
  if (recipe.base || recipe.variants || recipe.defaultVariants || recipe.compoundVariants) return recipe
  throw new Error('createRecipeContext requires a recipe')
}

export function createRecipeContext(recipeInput) {
  const recipe = resolveRecipe(recipeInput)
  const PropsContext = createContext(undefined)
  const usePropsContext = () => useContext(PropsContext)

  const withContext = (Component, options) => {
    const StyledComponent = styled(Component, recipe, options)
    const componentName = getDisplayName(Component)

    const WithContext = forwardRef(function WithContext(inProps, ref) {
      const propsContext = usePropsContext()
      const props = propsContext ? Object.assign({}, propsContext, inProps) : inProps
      return createElement(StyledComponent, { ...props, ref })
    })

    WithContext.displayName = `withContext(${componentName})`
    return WithContext
  }

  return {
    withContext,
    PropsProvider: PropsContext.Provider,
    usePropsContext,
  }
}