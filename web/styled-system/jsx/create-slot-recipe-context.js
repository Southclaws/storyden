"use client";

import { createContext, useContext, createElement, forwardRef } from 'react';
import { cx, sva } from '../css/index';
import { styled } from './factory';
import { getDisplayName } from './helper';

function createSafeContext(contextName) {
  const Context = createContext(undefined)
  const useStyleContext = (componentName, slot) => {
    const context = useContext(Context)
    if (context === undefined) {
      const componentInfo = componentName ? `Component "${componentName}"` : 'A component'
      const slotInfo = slot ? ` (slot: "${slot}")` : ''
      throw new Error(`${componentInfo}${slotInfo} cannot access ${contextName} because it's missing its Provider.`)
    }
    return context
  }
  return [Context, useStyleContext]
}

function resolveSlotRecipe(recipe) {
  if (recipe == null) throw new Error('createSlotRecipeContext requires a slot recipe')
  if (typeof recipe.splitVariantProps === 'function') return recipe
  if (recipe.slots) return recipe
  throw new Error('createSlotRecipeContext requires a slot recipe')
}

export function createSlotRecipeContext(recipeInput) {
  const recipe = resolveSlotRecipe(recipeInput)
  const isRuntimeRecipe = typeof recipe.splitVariantProps === 'function'
  const isConfigRecipe = isRuntimeRecipe && recipe.__recipe__ !== undefined
  const recipeName = isRuntimeRecipe && recipe.__name__ ? recipe.__name__ : undefined
  const contextName = recipeName ? `createSlotRecipeContext("${recipeName}")` : 'createSlotRecipeContext'
  const [SlotStylesContext, useSlotStylesContext] = createSafeContext(contextName)
  const slotRecipeFn = isRuntimeRecipe ? recipe : sva(recipe.config ?? recipe)

  const resolveProps = (props, slotStyles) => {
    const { unstyled, ...restProps } = props
    if (unstyled) return restProps
    if (isConfigRecipe) return { ...restProps, className: cx(slotStyles, restProps.className) }
    return { ...slotStyles, ...restProps }
  }

  const withRootProvider = (Component, options) => {
    const WithRootProvider = (props) => {
      const [variantProps, otherProps] = slotRecipeFn.splitVariantProps(props)
      const resolvedSlots = isConfigRecipe ? slotRecipeFn(variantProps) : slotRecipeFn.raw(variantProps)
      const mergedProps = options?.defaultProps ? Object.assign({}, options.defaultProps, otherProps) : otherProps
      return createElement(SlotStylesContext.Provider, {
        value: resolvedSlots,
        children: createElement(Component, mergedProps),
      })
    }
    const componentName = getDisplayName(Component)
    WithRootProvider.displayName = `withRootProvider(${componentName})`
    return WithRootProvider
  }

  const withProvider = (Component, slot, options) => {
    const StyledComponent = styled(Component, {}, options)
    const WithProvider = forwardRef(function WithProvider(props, ref) {
      const [variantProps, restProps] = slotRecipeFn.splitVariantProps(props)
      const resolvedSlots = isConfigRecipe ? slotRecipeFn(variantProps) : slotRecipeFn.raw(variantProps)
      if (restProps.className == null && options?.defaultProps?.className) restProps.className = options.defaultProps.className
      const resolvedProps = resolveProps(restProps, resolvedSlots[slot])
      options?.forwardProps?.forEach((key) => {
        if (key in variantProps) resolvedProps[key] = variantProps[key]
      })
      return createElement(SlotStylesContext.Provider, {
        value: resolvedSlots,
        children: createElement(StyledComponent, {
          ...resolvedProps,
          'data-slot': slot,
          ref,
        }),
      })
    })
    const componentName = getDisplayName(Component)
    WithProvider.displayName = `withProvider(${componentName})`
    return WithProvider
  }

  const withContext = (Component, slot, options) => {
    const StyledComponent = styled(Component, {}, options)
    const componentName = getDisplayName(Component)
    const WithContext = forwardRef(function WithContext(props, ref) {
      const resolvedSlots = useSlotStylesContext(componentName, slot)
      const nextProps = props.className == null && options?.defaultProps?.className
        ? { ...props, className: options.defaultProps.className }
        : props
      const resolvedProps = resolveProps(nextProps, resolvedSlots[slot])
      return createElement(StyledComponent, {
        ...resolvedProps,
        'data-slot': slot,
        ref,
      })
    })
    WithContext.displayName = `withContext(${componentName})`
    return WithContext
  }

  return {
    withRootProvider,
    withProvider,
    withContext,
  }
}