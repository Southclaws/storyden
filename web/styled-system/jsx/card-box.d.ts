import type { FunctionComponent } from 'react';
import type { CardBoxProperties } from '../patterns/card-box';
import type { HTMLStyledProps } from '../types/jsx';
import type { DistributiveOmit } from '../types/system';

export interface CardBoxProps extends CardBoxProperties, DistributiveOmit<HTMLStyledProps<"div">, keyof CardBoxProperties> {}

export declare const CardBox: FunctionComponent<CardBoxProps>