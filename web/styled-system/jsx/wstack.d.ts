import type { FunctionComponent } from 'react';
import type { WstackProperties } from '../patterns/wstack';
import type { HTMLStyledProps } from '../types/jsx';
import type { DistributiveOmit } from '../types/system';

export interface WstackProps extends WstackProperties, DistributiveOmit<HTMLStyledProps<"div">, keyof WstackProperties> {}

export declare const WStack: FunctionComponent<WstackProps>