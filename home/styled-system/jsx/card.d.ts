/* eslint-disable */
import type { FunctionComponent } from 'react'
import type { CardProperties } from '../patterns/card';
import type { HTMLStyledProps } from '../types/jsx';
import type { DistributiveOmit } from '../types/system-types';

export interface CardProps extends CardProperties, DistributiveOmit<HTMLStyledProps<'div'>, keyof CardProperties > {}

/** A card component that can be used to display content in a container with a border and a shadow. */
export declare const Card: FunctionComponent<CardProps>