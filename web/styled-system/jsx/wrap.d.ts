/* eslint-disable */
import type { FunctionComponent } from 'react'
import type { WrapProperties } from '../patterns/wrap';
import type { HTMLStyledProps } from '../types/jsx';
import type { DistributiveOmit } from '../types/system-types';

export type WrapProps = WrapProperties & DistributiveOmit<HTMLStyledProps<'div'>, keyof WrapProperties >


export declare const Wrap: FunctionComponent<WrapProps>