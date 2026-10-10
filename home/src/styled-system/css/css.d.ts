import type { SystemStyleObject } from '../types/system';

type Styles = SystemStyleObject | undefined | null | false
type StyleList = Styles | StyleList[]

interface CssRawFunction {
  (styles: Styles): SystemStyleObject
  (styles: StyleList[]): SystemStyleObject
  (...styles: StyleList[]): SystemStyleObject
  (styles: Styles): SystemStyleObject
}

interface CssFunction {
  (styles: Styles): string
  (styles: StyleList[]): string
  (...styles: StyleList[]): string
  (styles: Styles): string

  raw: CssRawFunction
}

export declare const css: CssFunction;