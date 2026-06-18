import type { PatternRuntimeConfig } from '../types/pattern';
import type { SystemStyleObject } from '../types/system';

export interface FloatingProperties {
  className?: string
}

type FloatingRestStyles = Omit<SystemStyleObject, keyof FloatingProperties>

interface FloatingStyles extends FloatingProperties, FloatingRestStyles {}

interface FloatingPatternFn {
  (styles?: FloatingStyles): string
  raw: (styles?: FloatingStyles) => SystemStyleObject
}

export declare function FloatingRaw(styles?: FloatingStyles): SystemStyleObject;

export declare const Floating: FloatingPatternFn;