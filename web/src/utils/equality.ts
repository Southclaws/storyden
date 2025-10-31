/**
 * Deep equality comparison utility
 * Based on fast-deep-equal with structural comparison
 */
export function deepEqual(a: unknown, b: unknown): boolean {
  if (a === b) return true;

  if (a && b && typeof a === "object" && typeof b === "object") {
    if (Array.isArray(a) !== Array.isArray(b)) return false;

    if (Array.isArray(a)) {
      const length = a.length;
      if (length !== (b as unknown[]).length) return false;
      for (let i = length; i-- !== 0; ) {
        if (!deepEqual(a[i], (b as unknown[])[i])) return false;
      }
      return true;
    }

    if (a instanceof Date && b instanceof Date) {
      return a.getTime() === b.getTime();
    }

    if (a instanceof RegExp && b instanceof RegExp) {
      return a.toString() === b.toString();
    }

    const keys = Object.keys(a);
    const length = keys.length;

    if (length !== Object.keys(b).length) return false;

    for (let i = length; i-- !== 0; ) {
      const key = keys[i];
      if (key === undefined || !Object.prototype.hasOwnProperty.call(b, key))
        return false;
    }

    for (let i = length; i-- !== 0; ) {
      const key = keys[i];
      if (key === undefined) return false;
      if (
        !deepEqual(
          (a as Record<string, unknown>)[key],
          (b as Record<string, unknown>)[key],
        )
      ) {
        return false;
      }
    }

    return true;
  }

  return a !== a && b !== b; // NaN case
}
