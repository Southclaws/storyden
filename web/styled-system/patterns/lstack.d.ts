/* eslint-disable */
import type { SystemStyleObject, ConditionalValue } from '../types/index';
import type { Properties } from '../types/csstype';
import type { SystemProperties } from '../types/style-props';
import type { DistributiveOmit } from '../types/system-types';
import type { Tokens } from '../tokens/index';

export interface LstackProperties {
   
}

interface LstackStyles extends LstackProperties, DistributiveOmit<SystemStyleObject, keyof LstackProperties > {}

interface LstackPatternFn {
  (styles?: LstackStyles): string
  raw: (styles?: LstackStyles) => SystemStyleObject
}

/**
 * A VStack with full width aligned left.
 */
export declare const lstack: LstackPatternFn;
