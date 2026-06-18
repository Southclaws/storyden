import type { FunctionComponent } from 'react';
import type { LstackProperties } from '../patterns/lstack';
import type { HTMLStyledProps } from '../types/jsx';
import type { DistributiveOmit } from '../types/system';

export interface LstackProps extends LstackProperties, DistributiveOmit<HTMLStyledProps<"div">, keyof LstackProperties> {}

export declare const LStack: FunctionComponent<LstackProps>