import type { FunctionComponent } from 'react';
import type { CardProperties } from '../patterns/card';
import type { HTMLStyledProps } from '../types/jsx';
import type { DistributiveOmit } from '../types/system';

export interface CardProps extends CardProperties, DistributiveOmit<HTMLStyledProps<"div">, keyof CardProperties> {}

export declare const Card: FunctionComponent<CardProps>