/* eslint-disable */
import type { SystemStyleObject, ConditionalValue } from '../types/index';
import type { Properties } from '../types/csstype';
import type { SystemProperties } from '../types/style-props';
import type { DistributiveOmit } from '../types/system-types';
import type { Tokens } from '../tokens/index';

export interface LStackProperties {
   
}


interface LStackStyles extends LStackProperties, DistributiveOmit<SystemStyleObject, keyof LStackProperties > {}

interface LStackPatternFn {
  (styles?: LStackStyles): string
  raw: (styles?: LStackStyles) => SystemStyleObject
}

/**
 * A VStack with full width aligned left.


 */
export declare const LStack: LStackPatternFn;
