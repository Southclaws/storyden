export type FirstThatWorksMember = string | number;

export type FirstThatWorksMemberOf<T> = Extract<T, FirstThatWorksMember>;

export type FirstThatWorksFn = <
  T = FirstThatWorksMember,
  A extends FirstThatWorksMemberOf<T> = FirstThatWorksMemberOf<T>,
  B extends FirstThatWorksMemberOf<T> = FirstThatWorksMemberOf<T>,
  C extends FirstThatWorksMemberOf<T> = never,
  D extends FirstThatWorksMemberOf<T> = never,
  E extends FirstThatWorksMemberOf<T> = never,
  F extends FirstThatWorksMemberOf<T> = never,
>(
  first: A,
  second: B,
  third?: C,
  fourth?: D,
  fifth?: E,
  sixth?: F,
) => T extends FirstThatWorksMember ? A | B | C | D | E | F : FirstThatWorksMemberOf<T>;

export declare const firstThatWorks: FirstThatWorksFn;