import type { PatternRuntimeConfig } from '../types/pattern';
import type { SystemStyleObject } from '../types/system';

export interface FrostedGlassProperties {}

type FrostedGlassRestStyles = Omit<SystemStyleObject, keyof FrostedGlassProperties>

interface FrostedGlassStyles extends FrostedGlassProperties, FrostedGlassRestStyles {}

interface FrostedGlassPatternFn {
  (styles?: FrostedGlassStyles): string
  raw: (styles?: FrostedGlassStyles) => SystemStyleObject
  propKeys: Array<keyof FrostedGlassProperties>
}

export declare function frostedGlassRaw(styles?: FrostedGlassStyles): SystemStyleObject;

export declare const frostedGlass: FrostedGlassPatternFn;