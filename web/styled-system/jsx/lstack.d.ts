/* eslint-disable */
import type { FunctionComponent } from 'react'
import type { LStackProperties } from '../patterns/lstack';
import type { HTMLStyledProps } from '../types/jsx';
import type { DistributiveOmit } from '../types/system-types';

export interface LStackProps extends LStackProperties, DistributiveOmit<HTMLStyledProps<'div'>, keyof LStackProperties > {}

/** A VStack with full width aligned left. */
export declare const LStack: FunctionComponent<LStackProps>