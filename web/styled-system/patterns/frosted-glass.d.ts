import type { PatternRuntimeConfig } from '../types/pattern';
import type { SystemStyleObject } from '../types/system';

export interface FrostedGlassProperties {
  className?: string
}

type FrostedGlassRestStyles = Omit<SystemStyleObject, keyof FrostedGlassProperties>

interface FrostedGlassStyles extends FrostedGlassProperties, FrostedGlassRestStyles {}

interface FrostedGlassPatternFn {
  (styles?: FrostedGlassStyles): string
  raw: (styles?: FrostedGlassStyles) => SystemStyleObject
}

export declare function FrostedGlassRaw(styles?: FrostedGlassStyles): SystemStyleObject;

export declare const FrostedGlass: FrostedGlassPatternFn;