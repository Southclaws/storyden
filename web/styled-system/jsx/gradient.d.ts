/* eslint-disable */
import type { FunctionComponent } from 'react'
import type { GradientProperties } from '../patterns/gradient';
import type { HTMLStyledProps } from '../types/jsx';
import type { DistributiveOmit } from '../types/system-types';

export interface GradientProps extends GradientProperties, DistributiveOmit<HTMLStyledProps<'div'>, keyof GradientProperties > {}

/** A gradient effect that can be used to create a gradient background for elements. */
export declare const Gradient: FunctionComponent<GradientProps>