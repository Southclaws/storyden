/* eslint-disable */
import type { SystemStyleObject, ConditionalValue } from '../types/index';
import type { Properties } from '../types/csstype';
import type { SystemProperties } from '../types/style-props';
import type { DistributiveOmit } from '../types/system-types';
import type { Tokens } from '../tokens/index';

export interface WstackProperties {
   
}

interface WstackStyles extends WstackProperties, DistributiveOmit<SystemStyleObject, keyof WstackProperties > {}

interface WstackPatternFn {
  (styles?: WstackStyles): string
  raw: (styles?: WstackStyles) => SystemStyleObject
}

/**
 * A HStack with full width and spaced children.
 */
export declare const wstack: WstackPatternFn;
