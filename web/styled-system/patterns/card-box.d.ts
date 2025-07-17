/* eslint-disable */
import type { SystemStyleObject, ConditionalValue } from '../types/index';
import type { Properties } from '../types/csstype';
import type { SystemProperties } from '../types/style-props';
import type { DistributiveOmit } from '../types/system-types';
import type { Tokens } from '../tokens/index';

export interface CardBoxProperties {
   kind?: ConditionalValue<"edge" | "default">
	display?: SystemProperties["display"]
}

interface CardBoxStyles extends CardBoxProperties, DistributiveOmit<SystemStyleObject, keyof CardBoxProperties > {}

interface CardBoxPatternFn {
  (styles?: CardBoxStyles): string
  raw: (styles?: CardBoxStyles) => SystemStyleObject
}

/**
 * A card component that can be used to display content in a container with a border and a shadow.
 */
export declare const CardBox: CardBoxPatternFn;
