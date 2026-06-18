import type { FunctionComponent } from 'react';
import type { FrostedGlassProperties } from '../patterns/frosted-glass';
import type { HTMLStyledProps } from '../types/jsx';
import type { DistributiveOmit } from '../types/system';

export interface FrostedGlassProps extends FrostedGlassProperties, DistributiveOmit<HTMLStyledProps<"div">, keyof FrostedGlassProperties> {}

export declare const FrostedGlass: FunctionComponent<FrostedGlassProps>