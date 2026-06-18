import { getCompoundVariantCss, memo, mergeProps, splitProps, toVariantMap, uniq, withDefaults } from '../helpers';
import { css, mergeCss } from './css';

export const cva = (config) => {
  const defaults = (c) => ({ base: {}, variants: {}, defaultVariants: {}, compoundVariants: [], ...c })
  const { base, variants, defaultVariants, compoundVariants } = defaults(config)

  const getVariantProps = (props) => withDefaults(defaultVariants, props)

  const resolve = (props = {}) => {
    const computed = getVariantProps(props)
    const styles = [base]
    for (const key in computed) {
      const value = computed[key]
      if (variants[key]?.[value]) styles.push(variants[key][value])
    }
    styles.push(getCompoundVariantCss(compoundVariants, computed))
    return mergeCss(...styles)
  }

  const variantKeys = Object.keys(variants)
  const variantMap = toVariantMap(variants)

  const merge = (other) => {
    const override = defaults(other.config)
    const keys = uniq(other.variantKeys, variantKeys)
    return cva({
      base: mergeCss(base, override.base),
      variants: Object.fromEntries(keys.map((key) => [key, mergeCss(variants[key], override.variants[key])])),
      defaultVariants: mergeProps(defaultVariants, override.defaultVariants),
      compoundVariants: [...compoundVariants, ...override.compoundVariants],
    })
  }

  return Object.assign(memo(function cvaFn(props) {
    return css(resolve(props))
  }), {
    __cva__: true,
    variantMap,
    variantKeys,
    raw: resolve,
    config,
    merge,
    splitVariantProps(props) {
      return splitProps(props, variantKeys)
    },
    getVariantProps,
  })
}