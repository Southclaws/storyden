/* eslint-disable */
import type { FunctionComponent } from 'react'
import type { WstackProperties } from '../patterns/wstack';
import type { HTMLStyledProps } from '../types/jsx';
import type { DistributiveOmit } from '../types/system-types';

export interface WstackProps extends WstackProperties, DistributiveOmit<HTMLStyledProps<'div'>, keyof WstackProperties > {}

/**
 * A HStack with full width and spaced children.
 */
export declare const WStack: FunctionComponent<WstackProps>