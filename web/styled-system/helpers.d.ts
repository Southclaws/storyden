export declare function isObject(v: unknown): v is Record<string, unknown>;

export declare function isBaseCondition(v: string): boolean;

export declare function filterBaseConditions(c: string[]): string[];

export declare function toHash(v: string): string;

export declare function compact<T extends Record<string, unknown>>(v: T): Partial<T>;

export declare function withDefaults(defaults: Record<string, any>, props: Record<string, any>): Record<string, any>;

export declare function toVariantMap(variants: Record<string, any>): Record<string, any>;

export declare function getCompoundVariantClassNames(compoundVariants: Array<Record<string, any>>, variants: Record<string, any>, formatClassName?: (className: string) => string): string;

export declare function getCompoundVariantCss(compoundVariants: Array<Record<string, any>>, variants: Record<string, any>): Record<string, any>;

export declare function getSlotCompoundVariant(compoundVariants: Array<Record<string, any>>, slot: string): Array<Record<string, any>>;

export declare function getSlotRecipes(recipe: Record<string, any>): Record<string, any>;

export declare function toResponsiveObject(values: any[], breakpoints: string[]): Record<string, any>;

export declare function walkObject(target: unknown, fn: (value: any, path: string[]) => any, options?: Record<string, any>): any;

export declare function mapObject(obj: unknown, fn: (value: any) => any): any;

export declare function normalizeStyleObject(styles: Record<string, any>, context: Record<string, any>, shorthand?: boolean): Record<string, any>;

export declare function memo<T extends (...args: any[]) => any>(fn: T): T;

export declare function weakMemo<T extends (arg: any) => any>(fn: T): T;

export declare function mergeProps(...src: Array<Record<string, any> | undefined>): Record<string, any>;

export declare function createCssRuntime(context: Record<string, any>): { serializeCss: (...styles: any[]) => string; mergeCss: (...styles: any[]) => any; assignCss: (...styles: any[]) => any };

export declare function hypenateProperty(property: string): string;

export declare function splitProps<T extends Record<string, any>>(props: T, ...keys: Array<Array<keyof T> | ((key: keyof T) => boolean)>): any[];

export declare function normalizeHTMLProps(props: Record<string, any>): Record<string, any>
export declare namespace normalizeHTMLProps {
  export const keys: string[]
}

export declare function uniq<T>(...items: Array<T[] | undefined>): T[];

export declare function withoutSpace<T extends string | number | boolean>(str: T): T;