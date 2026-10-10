import type { SystemStyleObject } from '../types/system';

export type ViewTransitionStyleObject = { group?: SystemStyleObject; imagePair?: SystemStyleObject; old?: SystemStyleObject; new?: SystemStyleObject };

export type ViewTransitionFn = (options: ViewTransitionStyleObject) => string;

export declare const viewTransition: ViewTransitionFn;