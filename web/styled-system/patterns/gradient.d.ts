/* eslint-disable */
import type { SystemStyleObject, ConditionalValue } from '../types/index';
import type { Properties } from '../types/csstype';
import type { PropertyValue } from '../types/prop-type';
import type { DistributiveOmit } from '../types/system-types';
import type { Tokens } from '../tokens/index';

export interface GradientProperties {
   
}


interface GradientStyles extends GradientProperties, DistributiveOmit<SystemStyleObject, keyof GradientProperties > {}

interface GradientPatternFn {
  (styles?: GradientStyles): string
  raw: (styles?: GradientStyles) => SystemStyleObject
}

/** A gradient effect that can be used to create a gradient background for elements. */
export declare const Gradient: GradientPatternFn;
