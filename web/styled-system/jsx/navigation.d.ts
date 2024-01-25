/* eslint-disable */
import type { FunctionComponent } from 'react'
import type { NavigationProperties } from '../patterns/navigation';
import type { HTMLStyledProps } from '../types/jsx';
import type { DistributiveOmit } from '../types/system-types';

export interface NavigationProps extends NavigationProperties, DistributiveOmit<HTMLStyledProps<'div'>, keyof NavigationProperties > {}

/** Navigation overlay elements. */
export declare const Navigation: FunctionComponent<NavigationProps>