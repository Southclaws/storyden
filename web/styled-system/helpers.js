export function isObject(v) {
  return typeof v === "object" && v != null && !Array.isArray(v)
}

const HAS_OWN = Object.prototype.hasOwnProperty

export function isBaseCondition(v) {
  return v === "base"
}

export function filterBaseConditions(c) {
  const out = []
  for (let i = 0; i < c.length; i++) {
    if (!isBaseCondition(c[i])) out.push(c[i])
  }
  return out
}

export function toHash(v) {
  let h = 5381
  for (let i = v.length; i; ) h = (h * 33) ^ v.charCodeAt(--i)
  let x = h >>> 0, out = ''
  for (; x > 52; x = (x / 52) | 0) {
    const c = x % 52
    out = String.fromCharCode(c + (c > 25 ? 39 : 97)) + out
  }
  const c = x % 52
  return String.fromCharCode(c + (c > 25 ? 39 : 97)) + out
}

export function compact(v) {
  const out = Object.create(null)
  if (!v) return out
  for (const k in v) {
    if (v[k] !== void 0) out[k] = v[k]
  }
  return out
}

export function withDefaults(defaults, props) {
  return { ...defaults, ...compact(props) }
}

export function toVariantMap(variants) {
  const map = {}
  for (const key in variants) map[key] = Object.keys(variants[key])
  return map
}

export function getCompoundVariantClassNames(compoundVariants, variants, formatClassName) {
  const classes = []
  outer: for (const compound of compoundVariants) {
    for (const key in compound) {
      if (key === "css" || key === "className" || key === "classNames") continue
      const expected = compound[key]
      const actual = variants[key]
      if (Array.isArray(expected)) {
        if (!expected.includes(actual)) continue outer
      } else if (actual !== expected) {
        continue outer
      }
    }
    if (compound.className) classes.push(formatClassName ? formatClassName(compound.className) : compound.className)
  }
  return classes.join(" ")
}

export function getCompoundVariantCss(compoundVariants, variants) {
  let result = {}
  outer: for (const variant of compoundVariants) {
    for (const key in variant) {
      if (key === "css" || key === "className" || key === "classNames") continue
      const expected = variant[key]
      const actual = variants[key]
      if (Array.isArray(expected)) {
        if (!expected.includes(actual)) continue outer
      } else if (actual !== expected) {
        continue outer
      }
    }
    result = mergeProps(result, variant.css)
  }
  return result
}

export function getSlotCompoundVariant(compoundVariants, slot) {
  const result = []
  for (const variant of compoundVariants) {
    const css = variant.css?.[slot]
    const className = variant.classNames?.[slot] ?? variant.className
    if (!css && !className) continue
    const next = css ? { css } : {}
    for (const key in variant) {
      if (key === "css" || key === "className" || key === "classNames") continue
      next[key] = variant[key]
    }
    if (className) next.className = className
    result.push(next)
  }
  return result
}

export function getSlotRecipes(recipe) {
  const result = {}
  const slots = recipe.slots ?? []
  for (const slot of slots) {
    result[slot] = {
      className: recipe.className ? recipe.className + "__" + slot : slot,
      base: recipe.base?.[slot] ?? {},
      variants: {},
      defaultVariants: recipe.defaultVariants ?? {},
      compoundVariants: getSlotCompoundVariant(recipe.compoundVariants ?? [], slot),
    }
  }
  const variants = recipe.variants ?? {}
  for (const variantsKey in variants) {
    const variantGroup = variants[variantsKey]
    for (const slot of slots) {
      const group = result[slot].variants[variantsKey] = {}
      for (const variantKey in variantGroup) {
        group[variantKey] = variantGroup[variantKey][slot] ?? {}
      }
    }
  }
  return result
}

export function toResponsiveObject(values, breakpoints) {
  const out = Object.create(null)
  for (let i = 0; i < values.length; i++) {
    if (values[i] != null) out[breakpoints[i]] = values[i]
  }
  return out
}

export function walkObject(target, fn, options) {
  options ||= {}
  const { stop, getKey } = options
  const inner = (value, path = []) => {
    if (!value || typeof value !== "object") return fn(value, path)
    if (stop?.(value, path)) return fn(value, path)
    const out = Array.isArray(value) ? [] : Object.create(null)
    for (const prop in value) {
      if (!HAS_OWN.call(value, prop)) continue
      const key = getKey?.(prop, value[prop]) ?? prop
      path.push(key)
      const next = inner(value[prop], path)
      path.pop()
      if (next != null) out[key] = next
    }
    return out
  }
  return inner(target)
}

export function mapObject(obj, fn) {
  return Array.isArray(obj) ? obj.map(fn) : isObject(obj) ? walkObject(obj, fn) : fn(obj)
}

export function normalizeStyleObject(styles, context, shorthand) {
  const { utility, conditions } = context
  const { hasShorthand, resolveShorthand } = utility
  shorthand = shorthand !== false
  return walkObject(styles, (value) => {
    if (Array.isArray(value)) return toResponsiveObject(value, conditions.breakpoints.keys)
    return value
  }, {
    stop: Array.isArray,
    getKey: shorthand ? (prop) => hasShorthand ? resolveShorthand(prop) : prop : void 0
  })
}

export function memo(fn) {
  const cache = new Map()
  return ((...args) => {
    const key = JSON.stringify(args)
    if (cache.has(key)) {
      const out = cache.get(key)
      cache.delete(key)
      cache.set(key, out)
      return out
    }
    const out = fn(...args)
    cache.set(key, out)
    if (cache.size > 500) cache.delete(cache.keys().next().value)
    return out
  })
}

export function weakMemo(fn) {
  const cache = new WeakMap()
  return ((arg) => {
    if (!arg || typeof arg !== "object") return fn(arg)
    if (cache.has(arg)) return cache.get(arg)
    const out = fn(arg)
    cache.set(arg, out)
    return out
  })
}

export function mergeProps(...src) {
  const out = Object.create(null)
  for (const obj of src) {
    if (!obj) continue
    for (const k in obj) {
      if (!HAS_OWN.call(obj, k) || k === "__proto__" || k === "constructor" || k === "prototype") continue
      const prev = out[k]
      const next = obj[k]
      out[k] = isObject(prev) && isObject(next) ? mergeProps(prev, next) : next
    }
  }
  return out
}

export function createCssRuntime(context) {
  const { utility: u, hash, conditions: c } = context
  const fmt = (s) => u.prefix ? u.prefix + "-" + s : s
  const toClass = (paths, name) => {
    const parts = c.finalize(paths)
    parts.push(hash ? name : fmt(name))
    return hash ? fmt(u.toHash(parts, toHash)) : parts.join(":")
  }
  const serializeCss = weakMemo(function serializeCss({ base, ...styles } = {}) {
    const obj = normalizeStyleObject(base ? Object.assign(styles, base) : styles, context)
    const set = new Set()
    walkObject(obj, (value, paths) => {
      if (value == null) return
      const [prop, ...all] = c.shift(paths)
      const cond = filterBaseConditions(all)
      const res = u.transform(prop, withoutSpace(value))
      set.add(toClass(cond, res.className))
    })
    let out = ""
    for (const name of set) out += out ? " " + name : name
    return out
  })
  const resolve = (styles) => {
    const out = []
    const visit = (items) => {
      for (let i = 0; i < items.length; i++) {
        const style = items[i]
        if (Array.isArray(style)) {
          visit(style)
          continue
        }
        if (!isObject(style)) continue
        for (const key in style) {
          if (style[key] !== void 0) {
            out.push(style)
            break
          }
        }
      }
    }
    visit(styles)
    if (out.length < 2) return out
    for (let i = 0; i < out.length; i++) out[i] = normalizeStyleObject(out[i], context)
    return out
  }
  const mergeCss = function() {
    return mergeProps(...resolve(arguments))
  }
  const assignCss = function() {
    const out = {}
    const resolved = resolve(arguments)
    for (let i = 0; i < resolved.length; i++) Object.assign(out, resolved[i])
    return out
  }
  return { serializeCss, mergeCss, assignCss }
}

export function hypenateProperty(property) {
  return property.startsWith("--") ? property : property.replace(/[A-Z]/g, "-$&").replace(/^ms-/, "-ms-").toLowerCase()
}

export function splitProps(props, ...keys) {
  const desc = Object.getOwnPropertyDescriptors(props)
  const all = Object.keys(desc)
  const split = (ks) => {
    const out = Object.create(null)
    for (let i = 0; i < ks.length; i++) {
      const k = ks[i]
      if (desc[k]) {
        Object.defineProperty(out, k, desc[k])
        delete desc[k]
      }
    }
    return out
  }
  const out = []
  for (const key of keys) {
    if (Array.isArray(key)) {
      out.push(split(key))
      continue
    }
    const picked = []
    for (let i = 0; i < all.length; i++) {
      if (key(all[i])) picked.push(all[i])
    }
    out.push(split(picked))
  }
  out.push(split(all))
  return out
}

const htmlProps = ['htmlSize', 'htmlTranslate', 'htmlWidth', 'htmlHeight']

function convertHTMLProp(key) {
  return htmlProps.includes(key) ? key.replace('html', '').toLowerCase() : key
}

export function normalizeHTMLProps(props) {
  return Object.fromEntries(Object.entries(props).map(([key, value]) => [convertHTMLProp(key), value]))
}

normalizeHTMLProps.keys = htmlProps

export function uniq(...items) {
  const set = new Set()
  for (const values of items) {
    if (!values) continue
    for (let i = 0; i < values.length; i++) set.add(values[i])
  }
  return Array.from(set)
}

export function withoutSpace(str) {
  return (typeof str === "string" && str.indexOf(" ") >= 0 ? str.replaceAll(" ", "_") : str)
}