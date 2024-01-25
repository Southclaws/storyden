/* eslint-disable */
import type { SystemStyleObject, ConditionalValue } from '../types/index';
import type { Properties } from '../types/csstype';
import type { PropertyValue } from '../types/prop-type';
import type { DistributiveOmit } from '../types/system-types';
import type { Tokens } from '../tokens/index';

export interface NavigationProperties {
   
}


interface NavigationStyles extends NavigationProperties, DistributiveOmit<SystemStyleObject, keyof NavigationProperties > {}

interface NavigationPatternFn {
  (styles?: NavigationStyles): string
  raw: (styles?: NavigationStyles) => SystemStyleObject
}

/** Navigation overlay elements. */
export declare const Navigation: NavigationPatternFn;
