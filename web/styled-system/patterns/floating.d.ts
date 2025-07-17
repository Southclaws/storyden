/* eslint-disable */
import type { SystemStyleObject, ConditionalValue } from '../types/index';
import type { Properties } from '../types/csstype';
import type { SystemProperties } from '../types/style-props';
import type { DistributiveOmit } from '../types/system-types';
import type { Tokens } from '../tokens/index';

export interface FloatingProperties {
   
}

interface FloatingStyles extends FloatingProperties, DistributiveOmit<SystemStyleObject, keyof FloatingProperties > {}

interface FloatingPatternFn {
  (styles?: FloatingStyles): string
  raw: (styles?: FloatingStyles) => SystemStyleObject
}

/**
 * Floating overlay elements.
 */
export declare const Floating: FloatingPatternFn;
