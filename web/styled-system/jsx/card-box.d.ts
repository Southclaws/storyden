/* eslint-disable */
import type { FunctionComponent } from 'react'
import type { CardBoxProperties } from '../patterns/card-box';
import type { HTMLStyledProps } from '../types/jsx';
import type { DistributiveOmit } from '../types/system-types';

export interface CardBoxProps extends CardBoxProperties, DistributiveOmit<HTMLStyledProps<'div'>, keyof CardBoxProperties > {}

/**
 * A card component that can be used to display content in a container with a border and a shadow.
 */
export declare const CardBox: FunctionComponent<CardBoxProps>