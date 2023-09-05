/* eslint-disable */
import type { FunctionComponent } from 'react'
import type { FlexProperties } from '../patterns/flex';
import type { HTMLStyledProps } from '../types/jsx';
import type { DistributiveOmit } from '../types/system-types';

export type FlexProps = FlexProperties & DistributiveOmit<HTMLStyledProps<'div'>, keyof FlexProperties >


export declare const Flex: FunctionComponent<FlexProps>