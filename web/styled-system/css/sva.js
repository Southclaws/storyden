import { getSlotRecipes, memo, splitProps, toVariantMap, withDefaults } from '../helpers';
import { cva } from './cva';
import { cx } from './cx';

export const sva = (config) => {
  const slotRecipes = getSlotRecipes(config)
  const slots = []
  for (const slot in slotRecipes) slots.push([slot, cva(slotRecipes[slot])])

  const defaultVariants = config.defaultVariants ?? {}

  const classNameMap = {}
  if (config.className) {
    for (const [slot, slotFn] of slots) classNameMap[slot] = slotFn.config.className
  }

  const variants = config.variants ?? {}
  const variantKeys = Object.keys(variants)
  const variantMap = toVariantMap(variants)

  const svaFn = (props) => {
    const result = {}
    for (const [slot, slotFn] of slots) result[slot] = cx(slotFn(props), classNameMap[slot])
    return result
  }

  const raw = (props) => {
    const result = {}
    for (const [slot, slotFn] of slots) result[slot] = slotFn.raw(props)
    return result
  }

  return Object.assign(memo(svaFn), {
    __cva__: false,
    raw,
    config,
    variantMap,
    variantKeys,
    classNameMap,
    splitVariantProps(props) {
      return splitProps(props, variantKeys)
    },
    getVariantProps(props) {
      return withDefaults(defaultVariants, props)
    },
  })
}