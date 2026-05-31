/* eslint-disable */
import type { FunctionComponent } from 'react'
import type { FloatingProperties } from '../patterns/floating';
import type { HTMLStyledProps } from '../types/jsx';
import type { DistributiveOmit } from '../types/system-types';

export interface FloatingProps extends FloatingProperties, DistributiveOmit<HTMLStyledProps<'div'>, keyof FloatingProperties > {}

/**
 * Floating overlay elements.
 */
export declare const Floating: FunctionComponent<FloatingProps>