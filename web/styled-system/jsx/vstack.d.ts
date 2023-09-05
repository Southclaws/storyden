/* eslint-disable */
import type { FunctionComponent } from 'react'
import type { VstackProperties } from '../patterns/vstack';
import type { HTMLStyledProps } from '../types/jsx';
import type { DistributiveOmit } from '../types/system-types';

export type VstackProps = VstackProperties & DistributiveOmit<HTMLStyledProps<'div'>, keyof VstackProperties >


export declare const VStack: FunctionComponent<VstackProps>