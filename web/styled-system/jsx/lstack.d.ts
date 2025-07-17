/* eslint-disable */
import type { FunctionComponent } from 'react'
import type { LstackProperties } from '../patterns/lstack';
import type { HTMLStyledProps } from '../types/jsx';
import type { DistributiveOmit } from '../types/system-types';

export interface LstackProps extends LstackProperties, DistributiveOmit<HTMLStyledProps<'div'>, keyof LstackProperties > {}

/**
 * A VStack with full width aligned left.
 */
export declare const LStack: FunctionComponent<LstackProps>