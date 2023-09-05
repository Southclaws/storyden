/* eslint-disable */
import type { FunctionComponent } from 'react'
import type { CenterProperties } from '../patterns/center';
import type { HTMLStyledProps } from '../types/jsx';
import type { DistributiveOmit } from '../types/system-types';

export type CenterProps = CenterProperties & DistributiveOmit<HTMLStyledProps<'div'>, keyof CenterProperties >


export declare const Center: FunctionComponent<CenterProps>