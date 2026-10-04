import type { FunctionComponent } from 'react';
import type { VstackProperties } from '../patterns/vstack';
import type { HTMLStyledProps } from '../types/jsx';
import type { DistributiveOmit } from '../types/system';

export interface VstackProps extends VstackProperties, DistributiveOmit<HTMLStyledProps<"div">, keyof VstackProperties> {}

export declare const VStack: FunctionComponent<VstackProps>