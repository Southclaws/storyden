/* eslint-disable */
import type { SystemStyleObject, ConditionalValue } from '../types/index';
import type { Properties } from '../types/csstype';
import type { SystemProperties } from '../types/style-props';
import type { DistributiveOmit } from '../types/system-types';
import type { Tokens } from '../tokens/index';

export interface FrostedGlassProperties {
   
}

interface FrostedGlassStyles extends FrostedGlassProperties, DistributiveOmit<SystemStyleObject, keyof FrostedGlassProperties > {}

interface FrostedGlassPatternFn {
  (styles?: FrostedGlassStyles): string
  raw: (styles?: FrostedGlassStyles) => SystemStyleObject
}

/**
 * A frosted glass effect for overlays, modals, menus, etc. This is most prominently used on the navigation overlays and menus.
 */
export declare const FrostedGlass: FrostedGlassPatternFn;
