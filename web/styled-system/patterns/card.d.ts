/* eslint-disable */
import type { SystemStyleObject, ConditionalValue } from '../types/index';
import type { Properties } from '../types/csstype';
import type { PropertyValue } from '../types/prop-type';
import type { DistributiveOmit } from '../types/system-types';
import type { Tokens } from '../tokens/index';

export interface CardProperties {
   kind?: ConditionalValue<"edge" | "default">
	display?: PropertyValue<'display'>
}


interface CardStyles extends CardProperties, DistributiveOmit<SystemStyleObject, keyof CardProperties > {}

interface CardPatternFn {
  (styles?: CardStyles): string
  raw: (styles?: CardStyles) => SystemStyleObject
}

/** A card component that can be used to display content in a container with a border and a shadow. */
export declare const Card: CardPatternFn;
