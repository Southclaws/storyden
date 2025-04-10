/* eslint-disable */
import type { SystemStyleObject, ConditionalValue } from '../types/index';
import type { Properties } from '../types/csstype';
import type { SystemProperties } from '../types/style-props';
import type { DistributiveOmit } from '../types/system-types';
import type { Tokens } from '../tokens/index';

export interface LinkButtonProperties {
   
}


interface LinkButtonStyles extends LinkButtonProperties, DistributiveOmit<SystemStyleObject, keyof LinkButtonProperties > {}

interface LinkButtonPatternFn {
  (styles?: LinkButtonStyles): string
  raw: (styles?: LinkButtonStyles) => SystemStyleObject
}

/**
 * Link button


 */
export declare const linkButton: LinkButtonPatternFn;
