import type { FunctionComponent } from 'react';
import type { CqProperties } from '../patterns/cq';
import type { HTMLStyledProps } from '../types/jsx';
import type { DistributiveOmit } from '../types/system';

export interface CqProps extends CqProperties, DistributiveOmit<HTMLStyledProps<"div">, keyof CqProperties> {}

export declare const Cq: FunctionComponent<CqProps>