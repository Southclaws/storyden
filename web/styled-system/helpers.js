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

export function stableStringify(value) {
  if (value === null) return 'null'
  const t = typeof value
  if (t === 'boolean') return value ? 'true' : 'false'
  if (t === 'number') return Number.isFinite(value) ? String(value) : 'null'
  if (t === 'string') return JSON.stringify(value)
  if (Array.isArray(value)) {
    let out = '['
    for (let i = 0; i < value.length; i++) {
      if (i) out += ','
      out += stableStringify(value[i])
    }
    return out + ']'
  }
  if (t === 'object') {
    const keys = Object.keys(value).sort()
    let out = '{'
    for (let i = 0; i < keys.length; i++) {
      if (i) out += ','
      const key = keys[i]
      out += JSON.stringify(key) + ':' + stableStringify(value[key])
    }
    return out + '}'
  }
  return 'null'
}

export function toCssVar(path) {
  let out = ""
  for (const ch of path.replaceAll(".", "-")) {
    if (ch >= "A" && ch <= "Z") out += "-" + ch.toLowerCase()
    else if (/[a-z0-9_-]/.test(ch) || ch >= "\u0081") out += ch
    else out += "\\" + ch
  }
  return "var(--" + out + ")"
}

export function colorMix(tokens, path) {
  const colorPrefix = "colors."
  if (!path.startsWith(colorPrefix)) return

  const index = path.indexOf("/", colorPrefix.length)
  if (index === -1 || index === path.length - 1) return

  const colorPath = path.slice(0, index)
  if (tokens[colorPath] === undefined) return

  const rawOpacity = path.slice(index + 1)
  const opacity = tokens["opacity." + rawOpacity]
  const percent = opacity === undefined ? Number(rawOpacity) : Number(opacity) * 100
  if (Number.isNaN(percent)) return

  return "color-mix(in oklab, " + toCssVar(colorPath) + " " + percent + "%, transparent)"
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

function compoundVariantMatches(compound, variants) {
  for (const key in compound) {
    if (key === "css" || key === "className" || key === "classNames") continue
    const expected = compound[key]
    const actual = variants[key]
    if (Array.isArray(expected) ? !expected.includes(actual) : actual !== expected) return false
  }
  return true
}

export function getCompoundVariantClassNames(compoundVariants, variants, formatClassName) {
  const classes = []
  for (const compound of compoundVariants) {
    if (!compoundVariantMatches(compound, variants)) continue
    if (compound.className) classes.push(formatClassName ? formatClassName(compound.className) : compound.className)
  }
  return classes.join(" ")
}

export function getCompoundVariantCss(compoundVariants, variants) {
  let result = {}
  for (const variant of compoundVariants) {
    if (!compoundVariantMatches(variant, variants)) continue
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

function flatHashOrNull(args) {
  let h = 5381
  for (let a = 0; a < args.length; a++) {
    const obj = args[a]
    if (obj === null || typeof obj !== "object") { h = (h * 33) ^ 1; continue }
    for (const k in obj) {
      const v = obj[k]
      const tv = typeof v
      if (v !== null && tv === "object") return null
      for (let i = 0; i < k.length; i++) h = (h * 33) ^ k.charCodeAt(i)
      if (tv === "string") { for (let i = 0; i < v.length; i++) h = (h * 33) ^ v.charCodeAt(i) }
      else if (tv === "number") h = (h * 33) ^ (v | 0)
      else if (tv === "boolean") h = (h * 33) ^ (v ? 991 : 997)
      else h = (h * 33) ^ 2
    }
  }
  return h >>> 0
}

function flatArgsEqual(a, b) {
  if (a.length !== b.length) return false
  for (let i = 0; i < a.length; i++) {
    const oa = a[i]
    const ob = b[i]
    if (oa === ob) continue
    if (oa === null || ob === null || typeof oa !== "object" || typeof ob !== "object") return false
    let n = 0
    for (const k in oa) {
      if (oa[k] !== ob[k]) return false
      n++
    }
    if (n !== Object.keys(ob).length) return false
  }
  return true
}

export function memo(fn) {
  const cache = new Map()
  const stringCache = new Map()
  const seen = new WeakSet()
  const newNode = () => ({ objects: new WeakMap(), prims: new Map(), out: void 0, has: false })
  const root = newNode()
  let lastHash
  let lastKey
  let lastValue
  let hasLast = false
  let misses = 0

  const step = (node, v) => {
    if (v !== null && typeof v === "object") {
      let next = node.objects.get(v)
      if (next === void 0) { next = newNode(); node.objects.set(v, next) }
      return next
    }
    let next = node.prims.get(v)
    if (next === void 0) {
      if (node.prims.size > 64) node.prims.clear()
      next = newNode()
      node.prims.set(v, next)
    }
    return next
  }
  const walk = (node, v) => {
    if (Array.isArray(v)) {
      node = step(node, "\u0000[")
      for (let i = 0; i < v.length; i++) node = walk(node, v[i])
      return step(node, "\u0000]")
    }
    return step(node, v)
  }
  const readWalk = (node, v) => {
    if (node === void 0) return void 0
    if (Array.isArray(v)) {
      node = node.prims.get("\u0000[")
      for (let i = 0; i < v.length && node !== void 0; i++) node = readWalk(node, v[i])
      return node === void 0 ? void 0 : node.prims.get("\u0000]")
    }
    return v !== null && typeof v === "object" ? node.objects.get(v) : node.prims.get(v)
  }
  const markSeen = (v) => {
    if (Array.isArray(v)) {
      let all = true
      for (let i = 0; i < v.length; i++) if (!markSeen(v[i])) all = false
      return all
    }
    if (v === null || typeof v !== "object") return true
    if (seen.has(v)) return true
    seen.add(v)
    return false
  }

  return ((...args) => {
    let composed = false
    for (let i = 0; i < args.length; i++) if (Array.isArray(args[i])) { composed = true; break }

    if (composed) {
      let node = root
      for (let i = 0; i < args.length && node !== void 0; i++) node = readWalk(node, args[i])
      if (node !== void 0 && node.has) return node.out

      const composedKey = JSON.stringify(args)
      let composedOut =
        hasLast && lastHash === void 0 && composedKey === lastKey ? lastValue : stringCache.get(composedKey)
      if (composedOut === void 0) {
        composedOut = fn(...args)
        stringCache.set(composedKey, composedOut)
        if (stringCache.size > 500) stringCache.delete(stringCache.keys().next().value)
      }

      let reused = false
      if ((++misses & 3) === 0) {
        reused = true
        for (let i = 0; i < args.length; i++) if (!markSeen(args[i])) reused = false
      }
      if (reused) {
        let insert = root
        for (let i = 0; i < args.length; i++) insert = walk(insert, args[i])
        insert.out = composedOut
        insert.has = true
      }
      lastHash = void 0
      lastKey = composedKey
      lastValue = composedOut
      hasLast = true
      return composedOut
    }

    const hash = flatHashOrNull(args)
    if (hash !== null) {
      if (hasLast && lastHash === hash && flatArgsEqual(args, lastKey)) return lastValue
      let bucket = cache.get(hash)
      if (bucket) {
        for (let i = 0; i < bucket.length; i++) {
          if (flatArgsEqual(args, bucket[i].args)) {
            lastHash = hash
            lastKey = args
            lastValue = bucket[i].out
            hasLast = true
            return bucket[i].out
          }
        }
      }
      const out = fn(...args)
      if (!bucket) {
        bucket = []
        cache.set(hash, bucket)
      }
      bucket.push({ args, out })
      if (bucket.length > 8) bucket.shift()
      if (cache.size > 500) cache.delete(cache.keys().next().value)
      lastHash = hash
      lastKey = args
      lastValue = out
      hasLast = true
      return out
    }

    const key = JSON.stringify(args)
    if (hasLast && lastHash === void 0 && key === lastKey) return lastValue
    const cached = stringCache.get(key)
    if (cached !== void 0) {
      lastHash = void 0
      lastKey = key
      lastValue = cached
      hasLast = true
      return cached
    }
    const out = fn(...args)
    stringCache.set(key, out)
    if (stringCache.size > 500) stringCache.delete(stringCache.keys().next().value)
    lastHash = void 0
    lastKey = key
    lastValue = out
    hasLast = true
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

const WHITESPACE_REGEX = /[\n\s]+/g
const sanitizeStyleValue = (value) => typeof value === "string" ? value.replace(WHITESPACE_REGEX, " ") : value

export function resolveStyleArgs(styles, context) {
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

export function createSerializeCss(context) {
  const u = context.utility
  const c = context.conditions
  const hash = context.hash
  const fmt = (s) => u.prefix ? u.prefix + "-" + s : s
  const toClass = (paths, name) => {
    const parts = c.finalize(paths)
    parts.push(hash ? name : fmt(name))
    return hash ? fmt(u.toHash(parts, toHash)) : parts.join(":")
  }
  return weakMemo(memo(function serializeCss({ base, ...styles } = {}) {
    const obj = normalizeStyleObject(base ? Object.assign(styles, base) : styles, context)
    const set = new Set()
    walkObject(obj, (value, paths) => {
      if (value == null) return
      const important = isImportant(value)
      const [prop, ...all] = c.shift(paths)
      const cond = filterBaseConditions(all)
      const res = u.transform(prop, withoutSpace(withoutImportant(sanitizeStyleValue(value))))
      let name = toClass(cond, res.className)
      if (important) name += "!"
      set.add(name)
    })
    let out = ""
    for (const name of set) out += out ? " " + name : name
    return out
  }))
}

export function createMergeCss(context) {
  return function mergeCss() {
    return mergeProps(...resolveStyleArgs(arguments, context))
  }
}

export function createSerializeCssArgs(serializeCss, mergeCss) {
  return memo(function serializeCssArgs(...styles) {
    return serializeCss(mergeCss(...styles))
  })
}

export function createAssignCss(context) {
  return function assignCss() {
    const out = {}
    const resolved = resolveStyleArgs(arguments, context)
    for (let i = 0; i < resolved.length; i++) Object.assign(out, resolved[i])
    return out
  }
}

const HYPHENATE_PROPERTY_REGEX = /[A-Z]/g
const MS_PROPERTY_REGEX = /^ms-/

export function hypenateProperty(property) {
  return property.startsWith("--") ? property : property.replace(HYPHENATE_PROPERTY_REGEX, "-$&").replace(MS_PROPERTY_REGEX, "-ms-").toLowerCase()
}

const splitPropsKeyMapCache = new WeakMap()

function getSplitPropsKeyMap(keys) {
  let keyMap = splitPropsKeyMapCache.get(keys)
  if (keyMap) return keyMap
  keyMap = Object.create(null)
  for (let i = 0; i < keys.length; i++) keyMap[keys[i]] = true
  splitPropsKeyMapCache.set(keys, keyMap)
  return keyMap
}

function copySplitProp(source, target, key) {
  const desc = Object.getOwnPropertyDescriptor(source, key)
  if (desc?.get || desc?.set) {
    Object.defineProperty(target, key, desc)
    return
  }
  target[key] = source[key]
}

export function splitProps(props, ...keys) {
  const propKeys = Object.keys(props)
  const keyCount = keys.length

  if (keyCount === 1) {
    const matcher = keys[0]
    const picked = Object.create(null)
    const rest = Object.create(null)

    if (Array.isArray(matcher)) {
      const keyMap = getSplitPropsKeyMap(matcher)
      for (let i = 0; i < propKeys.length; i++) {
        const key = propKeys[i]
        if (keyMap[key] === true) {
          copySplitProp(props, picked, key)
        } else {
          copySplitProp(props, rest, key)
        }
      }
      return [picked, rest]
    }

    for (let i = 0; i < propKeys.length; i++) {
      const key = propKeys[i]
      if (matcher(key)) {
        copySplitProp(props, picked, key)
      } else {
        copySplitProp(props, rest, key)
      }
    }
    return [picked, rest]
  }

  const matchers = new Array(keyCount)
  for (let i = 0; i < keyCount; i++) {
    const matcher = keys[i]
    matchers[i] = Array.isArray(matcher) ? getSplitPropsKeyMap(matcher) : matcher
  }

  const out = new Array(keyCount + 1)
  for (let i = 0; i <= keyCount; i++) out[i] = Object.create(null)
  const rest = out[keyCount]

  for (let i = 0; i < propKeys.length; i++) {
    const key = propKeys[i]
    let matched = false
    for (let j = 0; j < keyCount; j++) {
      const matcher = matchers[j]
      if (typeof matcher === "function" ? matcher(key) : matcher[key] === true) {
        copySplitProp(props, out[j], key)
        matched = true
        break
      }
    }
    if (!matched) copySplitProp(props, rest, key)
  }

  return out
}

const htmlProps = ['htmlSize', 'htmlTranslate', 'htmlWidth', 'htmlHeight']

function convertHTMLProp(key) {
  return htmlProps.includes(key) ? key.replace('html', '').toLowerCase() : key
}

// `Object.assign` keeps `.keys` inside a pure expression, so an unused helper tree-shakes away.
export const normalizeHTMLProps = /* @__PURE__ */ Object.assign(
  function normalizeHTMLProps(props) {
    return Object.fromEntries(Object.entries(props).map(([key, value]) => [convertHTMLProp(key), value]))
  },
  { keys: htmlProps },
)

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

export function isImportant(value) {
  return typeof value === "string" ? /\s*!(important)?\s*$/i.test(value) : false
}

export function withoutImportant(value) {
  return (typeof value === "string" ? value.replace(/\s*!(important)?\s*$/i, "").trim() : value)
}