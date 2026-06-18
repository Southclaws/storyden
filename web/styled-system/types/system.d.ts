import type { TokenValue } from './tokens';

export type Pretty<T> = { [K in keyof T]: T[K] } & {}

export type DistributiveOmit<T, K extends keyof any> = T extends unknown ? Omit<T, K> : never

export type DistributiveUnion<T, U> = {
  [K in keyof T]: K extends keyof U ? U[K] | T[K] : T[K]
} & DistributiveOmit<U, keyof T>

export type Assign<T, U> = {
  [K in keyof T]: K extends keyof U ? U[K] : T[K]
} & U

export interface Conditions {
  "2xl": string
  "2xlDown": string
  "2xlOnly": string
  "_active": string
  "_after": string
  "_atValue": string
  "_autofill": string
  "_backdrop": string
  "_before": string
  "_checked": string
  "_closed": string
  "_collapsed": string
  "_complete": string
  "_containerLarge": string
  "_containerMedium": string
  "_containerSmall": string
  "_current": string
  "_currentPage": string
  "_currentStep": string
  "_dark": string
  "_default": string
  "_disabled": string
  "_dragging": string
  "_empty": string
  "_enabled": string
  "_even": string
  "_expanded": string
  "_file": string
  "_first": string
  "_firstLetter": string
  "_firstLine": string
  "_firstOfType": string
  "_focus": string
  "_focusVisible": string
  "_focusWithin": string
  "_fullscreen": string
  "_grabbed": string
  "_groupActive": string
  "_groupChecked": string
  "_groupDisabled": string
  "_groupExpanded": string
  "_groupFocus": string
  "_groupFocusVisible": string
  "_groupFocusWithin": string
  "_groupHover": string
  "_groupInvalid": string
  "_hidden": string
  "_highContrast": string
  "_highlighted": string
  "_horizontal": string
  "_hover": string
  "_icon": string
  "_inRange": string
  "_incomplete": string
  "_indeterminate": string
  "_invalid": string
  "_invertedColors": string
  "_landscape": string
  "_last": string
  "_lastOfType": string
  "_lessContrast": string
  "_light": string
  "_loading": string
  "_ltr": string
  "_marker": string
  "_moreContrast": string
  "_motionReduce": string
  "_motionSafe": string
  "_noscript": string
  "_now": string
  "_odd": string
  "_off": string
  "_on": string
  "_only": string
  "_onlyOfType": string
  "_open": string
  "_optional": string
  "_osDark": string
  "_osLight": string
  "_outOfRange": string
  "_overValue": string
  "_peerActive": string
  "_peerChecked": string
  "_peerDisabled": string
  "_peerExpanded": string
  "_peerFocus": string
  "_peerFocusVisible": string
  "_peerFocusWithin": string
  "_peerHover": string
  "_peerInvalid": string
  "_peerPlaceholderShown": string
  "_placeholder": string
  "_placeholderShown": string
  "_portrait": string
  "_pressed": string
  "_print": string
  "_rangeEnd": string
  "_rangeStart": string
  "_readOnly": string
  "_readWrite": string
  "_required": string
  "_rtl": string
  "_scrollbar": string
  "_scrollbarThumb": string
  "_scrollbarTrack": string
  "_selected": string
  "_selection": string
  "_starting": string
  "_target": string
  "_today": string
  "_topmost": string
  "_unavailable": string
  "_underValue": string
  "_valid": string
  "_vertical": string
  "_visited": string
  "base": string
  "lg": string
  "lgDown": string
  "lgOnly": string
  "lgTo2xl": string
  "lgToXl": string
  "md": string
  "mdDown": string
  "mdOnly": string
  "mdTo2xl": string
  "mdToLg": string
  "mdToXl": string
  "sm": string
  "smDown": string
  "smOnly": string
  "smTo2xl": string
  "smToLg": string
  "smToMd": string
  "smToXl": string
  "xl": string
  "xlDown": string
  "xlOnly": string
  "xlTo2xl": string
}

export interface Breakpoints {
  "base": string
  "sm": string
  "md": string
  "lg": string
  "xl": string
  "2xl": string
}

export type ContainerName = AnyString
export type ContainerValue = ContainerName | `${ContainerName} / inline-size` | `${ContainerName} / size` | AnyString

export type Condition = keyof Conditions

export type ConditionalValue<T> =
  | T
  | Array<T | null>
  | { [K in Condition]?: ConditionalValue<T> }

export type AnyString = string & {}

export type AnyNumber = number & {}

export type CssVars = `var(--${string})`

type WithColorOpacityModifier<T> = [T] extends [string] ? `${T}/${string}` & { __colorOpacityModifier?: true } : never

type ImportantMark = "!" | "!important"
type WhitespaceImportant = ` ${ImportantMark}`
type Important = ImportantMark | WhitespaceImportant
type WithImportant<T> = [T] extends [string] ? `${T}${Important}` & { __important?: true } : never

export type WithEscapeHatch<T> = T | `[${string}]` | WithColorOpacityModifier<T> | WithImportant<T>

export type OnlyKnown<Value> = Value extends boolean ? Value : Value extends `${infer _}` ? Value : never

export type AlignContentValue = WithEscapeHatch<Globals | CssVars | OnlyKnown<"baseline" | "center" | "end" | "flex-end" | "flex-start" | "normal" | "space-around" | "space-between" | "space-evenly" | "start" | "stretch">>

export type AlignItemsValue = WithEscapeHatch<Globals | CssVars | OnlyKnown<"anchor-center" | "baseline" | "center" | "end" | "flex-end" | "flex-start" | "normal" | "self-end" | "self-start" | "start" | "stretch">>

export type AlignSelfValue = WithEscapeHatch<Globals | CssVars | OnlyKnown<"anchor-center" | "auto" | "baseline" | "center" | "end" | "flex-end" | "flex-start" | "normal" | "self-end" | "self-start" | "start" | "stretch">>

export type AnimationCompositionValue = WithEscapeHatch<Globals | CssVars | OnlyKnown<"accumulate" | "add" | "replace">>

export type AnimationDirectionValue = WithEscapeHatch<Globals | CssVars | OnlyKnown<"alternate" | "alternate-reverse" | "normal" | "reverse">>

export type AnimationFillModeValue = WithEscapeHatch<Globals | CssVars | OnlyKnown<"backwards" | "both" | "forwards" | "none">>

export type AnimationIterationCountValue = string | number | CssVars | AnyString

export type AnimationPlayStateValue = string | number | CssVars | AnyString

export type AnimationRangeEndValue = string | number | CssVars | AnyString

export type AnimationRangeStartValue = string | number | CssVars | AnyString

export type AnimationRangeValue = string | number | CssVars | AnyString

export type AnimationStateValue = string | number | CssVars | AnyString

export type AnimationTimelineValue = string | number | CssVars | AnyString

export type AnimationsValue = WithEscapeHatch<Globals | TokenValue<"animations"> | CssVars>

export type AppearanceValue = WithEscapeHatch<Globals | CssVars | OnlyKnown<"auto" | "button" | "checkbox" | "listbox" | "menulist" | "menulist-button" | "meter" | "none" | "progress-bar" | "radio" | "searchfield" | "textarea" | "textfield">>

export type AspectRatiosValue = WithEscapeHatch<AutoGlobals | TokenValue<"aspectRatios"> | CssVars>

export type AssetsValue = WithEscapeHatch<Globals | TokenValue<"assets"> | CssVars>

export type BackdropBrightnessValue = string | number | CssVars | AnyString

export type BackdropContrastValue = string | number | CssVars | AnyString

export type BackdropFilterValue = WithEscapeHatch<Globals | "auto" | CssVars>

export type BackdropGrayscaleValue = string | number | CssVars | AnyString

export type BackdropHueRotateValue = string | number | CssVars | AnyString

export type BackdropInvertValue = string | number | CssVars | AnyString

export type BackdropOpacityValue = string | number | CssVars | AnyString

export type BackdropSaturateValue = string | number | CssVars | AnyString

export type BackdropSepiaValue = string | number | CssVars | AnyString

export type BackfaceVisibilityValue = WithEscapeHatch<Globals | CssVars | OnlyKnown<"hidden" | "visible">>

export type BackgroundAttachmentValue = WithEscapeHatch<Globals | CssVars | OnlyKnown<"fixed" | "local" | "scroll">>

export type BackgroundBlendModeValue = string | number | CssVars | AnyString

export type BackgroundClipValue = WithEscapeHatch<Globals | CssVars | OnlyKnown<"border-area" | "border-box" | "content-box" | "padding-box" | "text">>

export type BackgroundConicValue = string | number | CssVars | AnyString

export type BackgroundGradientValue = WithEscapeHatch<Globals | "to-b" | "to-bl" | "to-br" | "to-l" | "to-r" | "to-t" | "to-tl" | "to-tr" | CssVars>

export type BackgroundLinearValue = WithEscapeHatch<Globals | "to-b" | "to-bl" | "to-br" | "to-l" | "to-r" | "to-t" | "to-tl" | "to-tr" | CssVars>

export type BackgroundOriginValue = string | number | CssVars | AnyString

export type BackgroundPositionValue = string | number | CssVars | AnyString

export type BackgroundPositionXValue = string | number | CssVars | AnyString

export type BackgroundPositionYValue = string | number | CssVars | AnyString

export type BackgroundRepeatValue = string | number | CssVars | AnyString

export type BackgroundSizeValue = string | number | CssVars | AnyString

export type BlockSizeValue = WithEscapeHatch<Globals | "0" | "0.5" | "1" | "1.5" | "1/2" | "1/3" | "1/4" | "1/5" | "1/6" | "10" | "10.5" | "11" | "12" | "14" | "16" | "2" | "2.5" | "2/3" | "2/4" | "2/5" | "2/6" | "20" | "24" | "28" | "2xl" | "2xs" | "3" | "3.5" | "3/4" | "3/5" | "3/6" | "32" | "36" | "3xl" | "4" | "4.5" | "4/5" | "4/6" | "40" | "44" | "48" | "4xl" | "5" | "5.5" | "5/6" | "52" | "56" | "5xl" | "6" | "60" | "64" | "6xl" | "7" | "7.5" | "72" | "7xl" | "8" | "80" | "8xl" | "9" | "9.5" | "96" | "auto" | "breakpoint-2xl" | "breakpoint-lg" | "breakpoint-md" | "breakpoint-sm" | "breakpoint-xl" | "dvh" | "fit" | "full" | "lg" | "lvh" | "max" | "md" | "min" | "prose" | "safeBottom" | "screen" | "scrollGutter" | "sm" | "svh" | "xl" | "xs" | CssVars>

export type BlursValue = WithEscapeHatch<Globals | TokenValue<"blurs"> | CssVars>

export type BorderCollapseValue = WithEscapeHatch<Globals | CssVars | OnlyKnown<"collapse" | "separate">>

export type BorderSpacingValue = WithEscapeHatch<Globals | "-1" | "-10" | "-11" | "-12" | "-14" | "-16" | "-2" | "-20" | "-24" | "-28" | "-3" | "-32" | "-36" | "-4" | "-40" | "-44" | "-48" | "-5" | "-52" | "-56" | "-6" | "-60" | "-64" | "-7" | "-72" | "-8" | "-80" | "-9" | "-96" | "-safeBottom" | "-safeTop" | "-scrollGutter" | "0" | "0.-5" | "0.5" | "1" | "1.-5" | "1.5" | "10" | "10.-5" | "10.5" | "11" | "12" | "14" | "16" | "2" | "2.-5" | "2.5" | "20" | "24" | "28" | "3" | "3.-5" | "3.5" | "32" | "36" | "4" | "4.-5" | "4.5" | "40" | "44" | "48" | "5" | "5.-5" | "5.5" | "52" | "56" | "6" | "60" | "64" | "7" | "7.-5" | "7.5" | "72" | "8" | "80" | "9" | "9.-5" | "9.5" | "96" | "auto" | "safeBottom" | "safeTop" | "scrollGutter" | CssVars>

export type BorderStyleValue = string | number | CssVars | AnyString

export type BorderStylesValue = WithEscapeHatch<Globals | TokenValue<"borderStyles"> | CssVars>

export type BorderWidthsValue = WithEscapeHatch<Globals | TokenValue<"borderWidths"> | CssVars>

export type BordersValue = WithEscapeHatch<Globals | TokenValue<"borders"> | CssVars>

export type BoxDecorationBreakValue = WithEscapeHatch<Globals | CssVars | OnlyKnown<"clone" | "slice">>

export type BoxSizeValue = WithEscapeHatch<Globals | "0" | "0.5" | "1" | "1.5" | "1/12" | "1/2" | "1/3" | "1/4" | "1/5" | "1/6" | "10" | "10.5" | "10/12" | "11" | "11/12" | "12" | "14" | "16" | "2" | "2.5" | "2/12" | "2/3" | "2/4" | "2/5" | "2/6" | "20" | "24" | "28" | "2xl" | "2xs" | "3" | "3.5" | "3/12" | "3/4" | "3/5" | "3/6" | "32" | "36" | "3xl" | "4" | "4.5" | "4/12" | "4/5" | "4/6" | "40" | "44" | "48" | "4xl" | "5" | "5.5" | "5/12" | "5/6" | "52" | "56" | "5xl" | "6" | "6/12" | "60" | "64" | "6xl" | "7" | "7.5" | "7/12" | "72" | "7xl" | "8" | "8/12" | "80" | "8xl" | "9" | "9.5" | "9/12" | "96" | "auto" | "breakpoint-2xl" | "breakpoint-lg" | "breakpoint-md" | "breakpoint-sm" | "breakpoint-xl" | "fit" | "full" | "lg" | "max" | "md" | "min" | "prose" | "safeBottom" | "screen" | "scrollGutter" | "sm" | "xl" | "xs" | CssVars>

export type BoxSizingValue = WithEscapeHatch<Globals | CssVars | OnlyKnown<"border-box" | "content-box">>

export type BreakpointsValue = WithEscapeHatch<Globals | TokenValue<"breakpoints"> | CssVars>

export type BrightnessValue = string | number | CssVars | AnyString

export type ClipPathValue = string | number | CssVars | AnyString

export type ColorPaletteValue = WithEscapeHatch<Globals | "accent" | "accent.dark" | "accent.dark.text" | "accent.light" | "accent.light.text" | "amber" | "amber.dark" | "amber.light" | "backgroundGradientH" | "backgroundGradientV" | "black" | "blue" | "blue.dark" | "blue.light" | "cardBackgroundGradient" | "conicGradient" | "current" | "gray" | "gray.dark" | "gray.light" | "green" | "green.dark" | "green.light" | "neutral" | "neutral.dark" | "neutral.light" | "orange" | "orange.dark" | "orange.light" | "pink" | "pink.dark" | "pink.light" | "red" | "red.dark" | "red.light" | "slate" | "slate.dark" | "slate.light" | "tomato" | "tomato.dark" | "tomato.light" | "transparent" | "white" | CssVars>

export type ColorsValue = WithEscapeHatch<ColorGlobals | TokenValue<"colors"> | CssVars>

export type ContainerNamesValue = WithEscapeHatch<Globals | TokenValue<"containerNames"> | CssVars>

export type ContainerTypeValue = string | number | CssVars | AnyString

export type ContainerValue = string | number | CssVars | AnyString

export type ContrastValue = string | number | CssVars | AnyString

export type CursorValue = WithEscapeHatch<Globals | TokenValue<"cursor"> | CssVars>

export type DebugValue = WithEscapeHatch<Globals | boolean | CssVars>

export type DisplayValue = WithEscapeHatch<Globals | CssVars | OnlyKnown<"-ms-flexbox" | "-ms-grid" | "-ms-inline-flexbox" | "-ms-inline-grid" | "-webkit-flex" | "-webkit-inline-flex" | "block" | "contents" | "flex" | "flow" | "flow-root" | "grid" | "inline" | "inline-block" | "inline-flex" | "inline-grid" | "inline-list-item" | "inline-table" | "list-item" | "none" | "ruby" | "ruby-base" | "ruby-base-container" | "ruby-text" | "ruby-text-container" | "run-in" | "table" | "table-caption" | "table-cell" | "table-column" | "table-column-group" | "table-footer-group" | "table-header-group" | "table-row" | "table-row-group">>

export type DropShadowsValue = WithEscapeHatch<Globals | TokenValue<"dropShadows"> | CssVars>

export type DurationsValue = WithEscapeHatch<Globals | TokenValue<"durations"> | CssVars>

export type EasingsValue = WithEscapeHatch<Globals | TokenValue<"easings"> | CssVars>

export type FilterValue = WithEscapeHatch<Globals | "auto" | CssVars>

export type FlexBasisValue = WithEscapeHatch<Globals | "0" | "0.5" | "1" | "1.5" | "1/12" | "1/2" | "1/3" | "1/4" | "1/5" | "1/6" | "10" | "10.5" | "10/12" | "11" | "11/12" | "12" | "14" | "16" | "2" | "2.5" | "2/12" | "2/3" | "2/4" | "2/5" | "2/6" | "20" | "24" | "28" | "2xl" | "2xs" | "3" | "3.5" | "3/12" | "3/4" | "3/5" | "3/6" | "32" | "36" | "3xl" | "4" | "4.5" | "4/12" | "4/5" | "4/6" | "40" | "44" | "48" | "4xl" | "5" | "5.5" | "5/12" | "5/6" | "52" | "56" | "5xl" | "6" | "6/12" | "60" | "64" | "6xl" | "7" | "7.5" | "7/12" | "72" | "7xl" | "8" | "8/12" | "80" | "8xl" | "9" | "9.5" | "9/12" | "96" | "breakpoint-2xl" | "breakpoint-lg" | "breakpoint-md" | "breakpoint-sm" | "breakpoint-xl" | "fit" | "full" | "lg" | "max" | "md" | "min" | "prose" | "safeBottom" | "scrollGutter" | "sm" | "xl" | "xs" | CssVars>

export type FlexDirectionValue = WithEscapeHatch<Globals | CssVars | OnlyKnown<"column" | "column-reverse" | "row" | "row-reverse">>

export type FlexGrowValue = string | number | CssVars | AnyString

export type FlexShrinkValue = string | number | CssVars | AnyString

export type FlexValue = WithEscapeHatch<Globals | "1" | "auto" | "initial" | "none" | CssVars>

export type FloatValue = WithEscapeHatch<Globals | CssVars | OnlyKnown<"end" | "start">>

export type FocusRingValue = WithEscapeHatch<Globals | "inside" | "mixed" | "none" | "outside" | CssVars>

export type FocusVisibleRingValue = WithEscapeHatch<Globals | "inside" | "mixed" | "none" | "outside" | CssVars>

export type FontFeatureSettingsValue = string | number | CssVars | AnyString

export type FontKerningValue = WithEscapeHatch<Globals | CssVars | OnlyKnown<"auto" | "none" | "normal">>

export type FontPaletteValue = string | number | CssVars | AnyString

export type FontSizeAdjustValue = string | number | CssVars | AnyString

export type FontSizesValue = WithEscapeHatch<Globals | TokenValue<"fontSizes"> | CssVars>

export type FontSmoothingValue = WithEscapeHatch<Globals | "antialiased" | "subpixel-antialiased" | CssVars>

export type FontVariantAlternatesValue = string | number | CssVars | AnyString

export type FontVariantCapsValue = string | number | CssVars | AnyString

export type FontVariantNumericValue = string | number | CssVars | AnyString

export type FontVariantValue = string | number | CssVars | AnyString

export type FontVariationSettingsValue = string | number | CssVars | AnyString

export type FontWeightsValue = WithEscapeHatch<Globals | TokenValue<"fontWeights"> | CssVars>

export type FontsValue = WithEscapeHatch<Globals | TokenValue<"fonts"> | CssVars>

export type GradientFromPositionValue = string | number | CssVars | AnyString

export type GradientToPositionValue = string | number | CssVars | AnyString

export type GradientViaPositionValue = string | number | CssVars | AnyString

export type GradientsValue = WithEscapeHatch<Globals | TokenValue<"gradients"> | CssVars>

export type GrayscaleValue = string | number | CssVars | AnyString

export type GridAutoColumnsValue = WithEscapeHatch<Globals | "fr" | "max" | "min" | CssVars>

export type GridAutoFlowValue = string | number | CssVars | AnyString

export type GridAutoRowsValue = WithEscapeHatch<Globals | "fr" | "max" | "min" | CssVars>

export type GridColumnEndValue = string | number | CssVars | AnyString

export type GridColumnStartValue = string | number | CssVars | AnyString

export type GridColumnValue = string | number | CssVars | AnyString

export type GridRowValue = string | number | CssVars | AnyString

export type GridTemplateColumnsValue = string | number | CssVars | AnyString

export type GridTemplateRowsValue = string | number | CssVars | AnyString

export type HeightValue = WithEscapeHatch<Globals | "0" | "0.5" | "1" | "1.5" | "1/2" | "1/3" | "1/4" | "1/5" | "1/6" | "10" | "10.5" | "11" | "12" | "14" | "16" | "2" | "2.5" | "2/3" | "2/4" | "2/5" | "2/6" | "20" | "24" | "28" | "2xl" | "2xs" | "3" | "3.5" | "3/4" | "3/5" | "3/6" | "32" | "36" | "3xl" | "4" | "4.5" | "4/5" | "4/6" | "40" | "44" | "48" | "4xl" | "5" | "5.5" | "5/6" | "52" | "56" | "5xl" | "6" | "60" | "64" | "6xl" | "7" | "7.5" | "72" | "7xl" | "8" | "80" | "8xl" | "9" | "9.5" | "96" | "auto" | "breakpoint-2xl" | "breakpoint-lg" | "breakpoint-md" | "breakpoint-sm" | "breakpoint-xl" | "dvh" | "fit" | "full" | "lg" | "lvh" | "max" | "md" | "min" | "prose" | "safeBottom" | "screen" | "scrollGutter" | "sm" | "svh" | "xl" | "xs" | CssVars>

export type HueRotateValue = string | number | CssVars | AnyString

export type HyphensValue = string | number | CssVars | AnyString

export type InlineSizeValue = WithEscapeHatch<Globals | "0" | "0.5" | "1" | "1.5" | "1/12" | "1/2" | "1/3" | "1/4" | "1/5" | "1/6" | "10" | "10.5" | "10/12" | "11" | "11/12" | "12" | "14" | "16" | "2" | "2.5" | "2/12" | "2/3" | "2/4" | "2/5" | "2/6" | "20" | "24" | "28" | "2xl" | "2xs" | "3" | "3.5" | "3/12" | "3/4" | "3/5" | "3/6" | "32" | "36" | "3xl" | "4" | "4.5" | "4/12" | "4/5" | "4/6" | "40" | "44" | "48" | "4xl" | "5" | "5.5" | "5/12" | "5/6" | "52" | "56" | "5xl" | "6" | "6/12" | "60" | "64" | "6xl" | "7" | "7.5" | "7/12" | "72" | "7xl" | "8" | "8/12" | "80" | "8xl" | "9" | "9.5" | "9/12" | "96" | "auto" | "breakpoint-2xl" | "breakpoint-lg" | "breakpoint-md" | "breakpoint-sm" | "breakpoint-xl" | "fit" | "full" | "lg" | "max" | "md" | "min" | "prose" | "safeBottom" | "screen" | "scrollGutter" | "sm" | "xl" | "xs" | CssVars>

export type InsetValue = WithEscapeHatch<Globals | "-1" | "-10" | "-11" | "-12" | "-14" | "-16" | "-2" | "-20" | "-24" | "-28" | "-3" | "-32" | "-36" | "-4" | "-40" | "-44" | "-48" | "-5" | "-52" | "-56" | "-6" | "-60" | "-64" | "-7" | "-72" | "-8" | "-80" | "-9" | "-96" | "-safeBottom" | "-safeTop" | "-scrollGutter" | "0" | "0.-5" | "0.5" | "1" | "1.-5" | "1.5" | "10" | "10.-5" | "10.5" | "11" | "12" | "14" | "16" | "2" | "2.-5" | "2.5" | "20" | "24" | "28" | "3" | "3.-5" | "3.5" | "32" | "36" | "4" | "4.-5" | "4.5" | "40" | "44" | "48" | "5" | "5.-5" | "5.5" | "52" | "56" | "6" | "60" | "64" | "7" | "7.-5" | "7.5" | "72" | "8" | "80" | "9" | "9.-5" | "9.5" | "96" | "auto" | "safeBottom" | "safeTop" | "scrollGutter" | CssVars>

export type InvertValue = string | number | CssVars | AnyString

export type JustifyContentValue = string | number | CssVars | AnyString

export type KeyframesValue = WithEscapeHatch<Globals | TokenValue<"keyframes"> | CssVars>

export type LetterSpacingsValue = WithEscapeHatch<Globals | TokenValue<"letterSpacings"> | CssVars>

export type LineClampValue = string | number | CssVars | AnyString

export type LineHeightsValue = WithEscapeHatch<Globals | TokenValue<"lineHeights"> | CssVars>

export type ListStylePositionValue = string | number | CssVars | AnyString

export type ListStyleTypeValue = string | number | CssVars | AnyString

export type ListStyleValue = string | number | CssVars | AnyString

export type MarginBlockEndValue = WithEscapeHatch<Globals | "-1" | "-10" | "-11" | "-12" | "-14" | "-16" | "-2" | "-20" | "-24" | "-28" | "-3" | "-32" | "-36" | "-4" | "-40" | "-44" | "-48" | "-5" | "-52" | "-56" | "-6" | "-60" | "-64" | "-7" | "-72" | "-8" | "-80" | "-9" | "-96" | "-safeBottom" | "-safeTop" | "-scrollGutter" | "0" | "0.-5" | "0.5" | "1" | "1.-5" | "1.5" | "10" | "10.-5" | "10.5" | "11" | "12" | "14" | "16" | "2" | "2.-5" | "2.5" | "20" | "24" | "28" | "3" | "3.-5" | "3.5" | "32" | "36" | "4" | "4.-5" | "4.5" | "40" | "44" | "48" | "5" | "5.-5" | "5.5" | "52" | "56" | "6" | "60" | "64" | "7" | "7.-5" | "7.5" | "72" | "8" | "80" | "9" | "9.-5" | "9.5" | "96" | "auto" | "safeBottom" | "safeTop" | "scrollGutter" | CssVars>

export type MarginBlockStartValue = WithEscapeHatch<Globals | "-1" | "-10" | "-11" | "-12" | "-14" | "-16" | "-2" | "-20" | "-24" | "-28" | "-3" | "-32" | "-36" | "-4" | "-40" | "-44" | "-48" | "-5" | "-52" | "-56" | "-6" | "-60" | "-64" | "-7" | "-72" | "-8" | "-80" | "-9" | "-96" | "-safeBottom" | "-safeTop" | "-scrollGutter" | "0" | "0.-5" | "0.5" | "1" | "1.-5" | "1.5" | "10" | "10.-5" | "10.5" | "11" | "12" | "14" | "16" | "2" | "2.-5" | "2.5" | "20" | "24" | "28" | "3" | "3.-5" | "3.5" | "32" | "36" | "4" | "4.-5" | "4.5" | "40" | "44" | "48" | "5" | "5.-5" | "5.5" | "52" | "56" | "6" | "60" | "64" | "7" | "7.-5" | "7.5" | "72" | "8" | "80" | "9" | "9.-5" | "9.5" | "96" | "auto" | "safeBottom" | "safeTop" | "scrollGutter" | CssVars>

export type MarginBlockValue = WithEscapeHatch<Globals | "-1" | "-10" | "-11" | "-12" | "-14" | "-16" | "-2" | "-20" | "-24" | "-28" | "-3" | "-32" | "-36" | "-4" | "-40" | "-44" | "-48" | "-5" | "-52" | "-56" | "-6" | "-60" | "-64" | "-7" | "-72" | "-8" | "-80" | "-9" | "-96" | "-safeBottom" | "-safeTop" | "-scrollGutter" | "0" | "0.-5" | "0.5" | "1" | "1.-5" | "1.5" | "10" | "10.-5" | "10.5" | "11" | "12" | "14" | "16" | "2" | "2.-5" | "2.5" | "20" | "24" | "28" | "3" | "3.-5" | "3.5" | "32" | "36" | "4" | "4.-5" | "4.5" | "40" | "44" | "48" | "5" | "5.-5" | "5.5" | "52" | "56" | "6" | "60" | "64" | "7" | "7.-5" | "7.5" | "72" | "8" | "80" | "9" | "9.-5" | "9.5" | "96" | "auto" | "safeBottom" | "safeTop" | "scrollGutter" | CssVars>

export type MarginBottomValue = WithEscapeHatch<Globals | "-1" | "-10" | "-11" | "-12" | "-14" | "-16" | "-2" | "-20" | "-24" | "-28" | "-3" | "-32" | "-36" | "-4" | "-40" | "-44" | "-48" | "-5" | "-52" | "-56" | "-6" | "-60" | "-64" | "-7" | "-72" | "-8" | "-80" | "-9" | "-96" | "-safeBottom" | "-safeTop" | "-scrollGutter" | "0" | "0.-5" | "0.5" | "1" | "1.-5" | "1.5" | "10" | "10.-5" | "10.5" | "11" | "12" | "14" | "16" | "2" | "2.-5" | "2.5" | "20" | "24" | "28" | "3" | "3.-5" | "3.5" | "32" | "36" | "4" | "4.-5" | "4.5" | "40" | "44" | "48" | "5" | "5.-5" | "5.5" | "52" | "56" | "6" | "60" | "64" | "7" | "7.-5" | "7.5" | "72" | "8" | "80" | "9" | "9.-5" | "9.5" | "96" | "auto" | "safeBottom" | "safeTop" | "scrollGutter" | CssVars>

export type MarginInlineEndValue = WithEscapeHatch<Globals | "-1" | "-10" | "-11" | "-12" | "-14" | "-16" | "-2" | "-20" | "-24" | "-28" | "-3" | "-32" | "-36" | "-4" | "-40" | "-44" | "-48" | "-5" | "-52" | "-56" | "-6" | "-60" | "-64" | "-7" | "-72" | "-8" | "-80" | "-9" | "-96" | "-safeBottom" | "-safeTop" | "-scrollGutter" | "0" | "0.-5" | "0.5" | "1" | "1.-5" | "1.5" | "10" | "10.-5" | "10.5" | "11" | "12" | "14" | "16" | "2" | "2.-5" | "2.5" | "20" | "24" | "28" | "3" | "3.-5" | "3.5" | "32" | "36" | "4" | "4.-5" | "4.5" | "40" | "44" | "48" | "5" | "5.-5" | "5.5" | "52" | "56" | "6" | "60" | "64" | "7" | "7.-5" | "7.5" | "72" | "8" | "80" | "9" | "9.-5" | "9.5" | "96" | "auto" | "safeBottom" | "safeTop" | "scrollGutter" | CssVars>

export type MarginInlineStartValue = WithEscapeHatch<Globals | "-1" | "-10" | "-11" | "-12" | "-14" | "-16" | "-2" | "-20" | "-24" | "-28" | "-3" | "-32" | "-36" | "-4" | "-40" | "-44" | "-48" | "-5" | "-52" | "-56" | "-6" | "-60" | "-64" | "-7" | "-72" | "-8" | "-80" | "-9" | "-96" | "-safeBottom" | "-safeTop" | "-scrollGutter" | "0" | "0.-5" | "0.5" | "1" | "1.-5" | "1.5" | "10" | "10.-5" | "10.5" | "11" | "12" | "14" | "16" | "2" | "2.-5" | "2.5" | "20" | "24" | "28" | "3" | "3.-5" | "3.5" | "32" | "36" | "4" | "4.-5" | "4.5" | "40" | "44" | "48" | "5" | "5.-5" | "5.5" | "52" | "56" | "6" | "60" | "64" | "7" | "7.-5" | "7.5" | "72" | "8" | "80" | "9" | "9.-5" | "9.5" | "96" | "auto" | "safeBottom" | "safeTop" | "scrollGutter" | CssVars>

export type MarginInlineValue = WithEscapeHatch<Globals | "-1" | "-10" | "-11" | "-12" | "-14" | "-16" | "-2" | "-20" | "-24" | "-28" | "-3" | "-32" | "-36" | "-4" | "-40" | "-44" | "-48" | "-5" | "-52" | "-56" | "-6" | "-60" | "-64" | "-7" | "-72" | "-8" | "-80" | "-9" | "-96" | "-safeBottom" | "-safeTop" | "-scrollGutter" | "0" | "0.-5" | "0.5" | "1" | "1.-5" | "1.5" | "10" | "10.-5" | "10.5" | "11" | "12" | "14" | "16" | "2" | "2.-5" | "2.5" | "20" | "24" | "28" | "3" | "3.-5" | "3.5" | "32" | "36" | "4" | "4.-5" | "4.5" | "40" | "44" | "48" | "5" | "5.-5" | "5.5" | "52" | "56" | "6" | "60" | "64" | "7" | "7.-5" | "7.5" | "72" | "8" | "80" | "9" | "9.-5" | "9.5" | "96" | "auto" | "safeBottom" | "safeTop" | "scrollGutter" | CssVars>

export type MarginLeftValue = WithEscapeHatch<Globals | "-1" | "-10" | "-11" | "-12" | "-14" | "-16" | "-2" | "-20" | "-24" | "-28" | "-3" | "-32" | "-36" | "-4" | "-40" | "-44" | "-48" | "-5" | "-52" | "-56" | "-6" | "-60" | "-64" | "-7" | "-72" | "-8" | "-80" | "-9" | "-96" | "-safeBottom" | "-safeTop" | "-scrollGutter" | "0" | "0.-5" | "0.5" | "1" | "1.-5" | "1.5" | "10" | "10.-5" | "10.5" | "11" | "12" | "14" | "16" | "2" | "2.-5" | "2.5" | "20" | "24" | "28" | "3" | "3.-5" | "3.5" | "32" | "36" | "4" | "4.-5" | "4.5" | "40" | "44" | "48" | "5" | "5.-5" | "5.5" | "52" | "56" | "6" | "60" | "64" | "7" | "7.-5" | "7.5" | "72" | "8" | "80" | "9" | "9.-5" | "9.5" | "96" | "auto" | "safeBottom" | "safeTop" | "scrollGutter" | CssVars>

export type MarginRightValue = WithEscapeHatch<Globals | "-1" | "-10" | "-11" | "-12" | "-14" | "-16" | "-2" | "-20" | "-24" | "-28" | "-3" | "-32" | "-36" | "-4" | "-40" | "-44" | "-48" | "-5" | "-52" | "-56" | "-6" | "-60" | "-64" | "-7" | "-72" | "-8" | "-80" | "-9" | "-96" | "-safeBottom" | "-safeTop" | "-scrollGutter" | "0" | "0.-5" | "0.5" | "1" | "1.-5" | "1.5" | "10" | "10.-5" | "10.5" | "11" | "12" | "14" | "16" | "2" | "2.-5" | "2.5" | "20" | "24" | "28" | "3" | "3.-5" | "3.5" | "32" | "36" | "4" | "4.-5" | "4.5" | "40" | "44" | "48" | "5" | "5.-5" | "5.5" | "52" | "56" | "6" | "60" | "64" | "7" | "7.-5" | "7.5" | "72" | "8" | "80" | "9" | "9.-5" | "9.5" | "96" | "auto" | "safeBottom" | "safeTop" | "scrollGutter" | CssVars>

export type MarginTopValue = WithEscapeHatch<Globals | "-1" | "-10" | "-11" | "-12" | "-14" | "-16" | "-2" | "-20" | "-24" | "-28" | "-3" | "-32" | "-36" | "-4" | "-40" | "-44" | "-48" | "-5" | "-52" | "-56" | "-6" | "-60" | "-64" | "-7" | "-72" | "-8" | "-80" | "-9" | "-96" | "-safeBottom" | "-safeTop" | "-scrollGutter" | "0" | "0.-5" | "0.5" | "1" | "1.-5" | "1.5" | "10" | "10.-5" | "10.5" | "11" | "12" | "14" | "16" | "2" | "2.-5" | "2.5" | "20" | "24" | "28" | "3" | "3.-5" | "3.5" | "32" | "36" | "4" | "4.-5" | "4.5" | "40" | "44" | "48" | "5" | "5.-5" | "5.5" | "52" | "56" | "6" | "60" | "64" | "7" | "7.-5" | "7.5" | "72" | "8" | "80" | "9" | "9.-5" | "9.5" | "96" | "auto" | "safeBottom" | "safeTop" | "scrollGutter" | CssVars>

export type MarginValue = WithEscapeHatch<Globals | "-1" | "-10" | "-11" | "-12" | "-14" | "-16" | "-2" | "-20" | "-24" | "-28" | "-3" | "-32" | "-36" | "-4" | "-40" | "-44" | "-48" | "-5" | "-52" | "-56" | "-6" | "-60" | "-64" | "-7" | "-72" | "-8" | "-80" | "-9" | "-96" | "-safeBottom" | "-safeTop" | "-scrollGutter" | "0" | "0.-5" | "0.5" | "1" | "1.-5" | "1.5" | "10" | "10.-5" | "10.5" | "11" | "12" | "14" | "16" | "2" | "2.-5" | "2.5" | "20" | "24" | "28" | "3" | "3.-5" | "3.5" | "32" | "36" | "4" | "4.-5" | "4.5" | "40" | "44" | "48" | "5" | "5.-5" | "5.5" | "52" | "56" | "6" | "60" | "64" | "7" | "7.-5" | "7.5" | "72" | "8" | "80" | "9" | "9.-5" | "9.5" | "96" | "auto" | "safeBottom" | "safeTop" | "scrollGutter" | CssVars>

export type MaskImageValue = string | number | CssVars | AnyString

export type MaskSizeValue = string | number | CssVars | AnyString

export type MaskValue = string | number | CssVars | AnyString

export type MaxBlockSizeValue = WithEscapeHatch<Globals | "0" | "0.5" | "1" | "1.5" | "1/2" | "1/3" | "1/4" | "1/5" | "1/6" | "10" | "10.5" | "11" | "12" | "14" | "16" | "2" | "2.5" | "2/3" | "2/4" | "2/5" | "2/6" | "20" | "24" | "28" | "2xl" | "2xs" | "3" | "3.5" | "3/4" | "3/5" | "3/6" | "32" | "36" | "3xl" | "4" | "4.5" | "4/5" | "4/6" | "40" | "44" | "48" | "4xl" | "5" | "5.5" | "5/6" | "52" | "56" | "5xl" | "6" | "60" | "64" | "6xl" | "7" | "7.5" | "72" | "7xl" | "8" | "80" | "8xl" | "9" | "9.5" | "96" | "auto" | "breakpoint-2xl" | "breakpoint-lg" | "breakpoint-md" | "breakpoint-sm" | "breakpoint-xl" | "dvh" | "fit" | "full" | "lg" | "lvh" | "max" | "md" | "min" | "prose" | "safeBottom" | "screen" | "scrollGutter" | "sm" | "svh" | "xl" | "xs" | CssVars>

export type MaxHeightValue = WithEscapeHatch<Globals | "0" | "0.5" | "1" | "1.5" | "1/2" | "1/3" | "1/4" | "1/5" | "1/6" | "10" | "10.5" | "11" | "12" | "14" | "16" | "2" | "2.5" | "2/3" | "2/4" | "2/5" | "2/6" | "20" | "24" | "28" | "2xl" | "2xs" | "3" | "3.5" | "3/4" | "3/5" | "3/6" | "32" | "36" | "3xl" | "4" | "4.5" | "4/5" | "4/6" | "40" | "44" | "48" | "4xl" | "5" | "5.5" | "5/6" | "52" | "56" | "5xl" | "6" | "60" | "64" | "6xl" | "7" | "7.5" | "72" | "7xl" | "8" | "80" | "8xl" | "9" | "9.5" | "96" | "auto" | "breakpoint-2xl" | "breakpoint-lg" | "breakpoint-md" | "breakpoint-sm" | "breakpoint-xl" | "dvh" | "fit" | "full" | "lg" | "lvh" | "max" | "md" | "min" | "prose" | "safeBottom" | "screen" | "scrollGutter" | "sm" | "svh" | "xl" | "xs" | CssVars>

export type MaxInlineSizeValue = WithEscapeHatch<Globals | "0" | "0.5" | "1" | "1.5" | "1/12" | "1/2" | "1/3" | "1/4" | "1/5" | "1/6" | "10" | "10.5" | "10/12" | "11" | "11/12" | "12" | "14" | "16" | "2" | "2.5" | "2/12" | "2/3" | "2/4" | "2/5" | "2/6" | "20" | "24" | "28" | "2xl" | "2xs" | "3" | "3.5" | "3/12" | "3/4" | "3/5" | "3/6" | "32" | "36" | "3xl" | "4" | "4.5" | "4/12" | "4/5" | "4/6" | "40" | "44" | "48" | "4xl" | "5" | "5.5" | "5/12" | "5/6" | "52" | "56" | "5xl" | "6" | "6/12" | "60" | "64" | "6xl" | "7" | "7.5" | "7/12" | "72" | "7xl" | "8" | "8/12" | "80" | "8xl" | "9" | "9.5" | "9/12" | "96" | "auto" | "breakpoint-2xl" | "breakpoint-lg" | "breakpoint-md" | "breakpoint-sm" | "breakpoint-xl" | "fit" | "full" | "lg" | "max" | "md" | "min" | "prose" | "safeBottom" | "screen" | "scrollGutter" | "sm" | "xl" | "xs" | CssVars>

export type MaxWidthValue = WithEscapeHatch<Globals | "0" | "0.5" | "1" | "1.5" | "1/12" | "1/2" | "1/3" | "1/4" | "1/5" | "1/6" | "10" | "10.5" | "10/12" | "11" | "11/12" | "12" | "14" | "16" | "2" | "2.5" | "2/12" | "2/3" | "2/4" | "2/5" | "2/6" | "20" | "24" | "28" | "2xl" | "2xs" | "3" | "3.5" | "3/12" | "3/4" | "3/5" | "3/6" | "32" | "36" | "3xl" | "4" | "4.5" | "4/12" | "4/5" | "4/6" | "40" | "44" | "48" | "4xl" | "5" | "5.5" | "5/12" | "5/6" | "52" | "56" | "5xl" | "6" | "6/12" | "60" | "64" | "6xl" | "7" | "7.5" | "7/12" | "72" | "7xl" | "8" | "8/12" | "80" | "8xl" | "9" | "9.5" | "9/12" | "96" | "auto" | "breakpoint-2xl" | "breakpoint-lg" | "breakpoint-md" | "breakpoint-sm" | "breakpoint-xl" | "fit" | "full" | "lg" | "max" | "md" | "min" | "prose" | "safeBottom" | "screen" | "scrollGutter" | "sm" | "xl" | "xs" | CssVars>

export type MinBlockSizeValue = WithEscapeHatch<Globals | "0" | "0.5" | "1" | "1.5" | "1/2" | "1/3" | "1/4" | "1/5" | "1/6" | "10" | "10.5" | "11" | "12" | "14" | "16" | "2" | "2.5" | "2/3" | "2/4" | "2/5" | "2/6" | "20" | "24" | "28" | "2xl" | "2xs" | "3" | "3.5" | "3/4" | "3/5" | "3/6" | "32" | "36" | "3xl" | "4" | "4.5" | "4/5" | "4/6" | "40" | "44" | "48" | "4xl" | "5" | "5.5" | "5/6" | "52" | "56" | "5xl" | "6" | "60" | "64" | "6xl" | "7" | "7.5" | "72" | "7xl" | "8" | "80" | "8xl" | "9" | "9.5" | "96" | "auto" | "breakpoint-2xl" | "breakpoint-lg" | "breakpoint-md" | "breakpoint-sm" | "breakpoint-xl" | "dvh" | "fit" | "full" | "lg" | "lvh" | "max" | "md" | "min" | "prose" | "safeBottom" | "screen" | "scrollGutter" | "sm" | "svh" | "xl" | "xs" | CssVars>

export type MinHeightValue = WithEscapeHatch<Globals | "0" | "0.5" | "1" | "1.5" | "1/2" | "1/3" | "1/4" | "1/5" | "1/6" | "10" | "10.5" | "11" | "12" | "14" | "16" | "2" | "2.5" | "2/3" | "2/4" | "2/5" | "2/6" | "20" | "24" | "28" | "2xl" | "2xs" | "3" | "3.5" | "3/4" | "3/5" | "3/6" | "32" | "36" | "3xl" | "4" | "4.5" | "4/5" | "4/6" | "40" | "44" | "48" | "4xl" | "5" | "5.5" | "5/6" | "52" | "56" | "5xl" | "6" | "60" | "64" | "6xl" | "7" | "7.5" | "72" | "7xl" | "8" | "80" | "8xl" | "9" | "9.5" | "96" | "auto" | "breakpoint-2xl" | "breakpoint-lg" | "breakpoint-md" | "breakpoint-sm" | "breakpoint-xl" | "dvh" | "fit" | "full" | "lg" | "lvh" | "max" | "md" | "min" | "prose" | "safeBottom" | "screen" | "scrollGutter" | "sm" | "svh" | "xl" | "xs" | CssVars>

export type MinInlineSizeValue = WithEscapeHatch<Globals | "0" | "0.5" | "1" | "1.5" | "1/12" | "1/2" | "1/3" | "1/4" | "1/5" | "1/6" | "10" | "10.5" | "10/12" | "11" | "11/12" | "12" | "14" | "16" | "2" | "2.5" | "2/12" | "2/3" | "2/4" | "2/5" | "2/6" | "20" | "24" | "28" | "2xl" | "2xs" | "3" | "3.5" | "3/12" | "3/4" | "3/5" | "3/6" | "32" | "36" | "3xl" | "4" | "4.5" | "4/12" | "4/5" | "4/6" | "40" | "44" | "48" | "4xl" | "5" | "5.5" | "5/12" | "5/6" | "52" | "56" | "5xl" | "6" | "6/12" | "60" | "64" | "6xl" | "7" | "7.5" | "7/12" | "72" | "7xl" | "8" | "8/12" | "80" | "8xl" | "9" | "9.5" | "9/12" | "96" | "auto" | "breakpoint-2xl" | "breakpoint-lg" | "breakpoint-md" | "breakpoint-sm" | "breakpoint-xl" | "fit" | "full" | "lg" | "max" | "md" | "min" | "prose" | "safeBottom" | "screen" | "scrollGutter" | "sm" | "xl" | "xs" | CssVars>

export type MinWidthValue = WithEscapeHatch<Globals | "0" | "0.5" | "1" | "1.5" | "1/12" | "1/2" | "1/3" | "1/4" | "1/5" | "1/6" | "10" | "10.5" | "10/12" | "11" | "11/12" | "12" | "14" | "16" | "2" | "2.5" | "2/12" | "2/3" | "2/4" | "2/5" | "2/6" | "20" | "24" | "28" | "2xl" | "2xs" | "3" | "3.5" | "3/12" | "3/4" | "3/5" | "3/6" | "32" | "36" | "3xl" | "4" | "4.5" | "4/12" | "4/5" | "4/6" | "40" | "44" | "48" | "4xl" | "5" | "5.5" | "5/12" | "5/6" | "52" | "56" | "5xl" | "6" | "6/12" | "60" | "64" | "6xl" | "7" | "7.5" | "7/12" | "72" | "7xl" | "8" | "8/12" | "80" | "8xl" | "9" | "9.5" | "9/12" | "96" | "auto" | "breakpoint-2xl" | "breakpoint-lg" | "breakpoint-md" | "breakpoint-sm" | "breakpoint-xl" | "fit" | "full" | "lg" | "max" | "md" | "min" | "prose" | "safeBottom" | "screen" | "scrollGutter" | "sm" | "xl" | "xs" | CssVars>

export type MixBlendModeValue = WithEscapeHatch<Globals | CssVars | OnlyKnown<"color" | "color-burn" | "color-dodge" | "darken" | "difference" | "exclusion" | "hard-light" | "hue" | "lighten" | "luminosity" | "multiply" | "normal" | "overlay" | "plus-darker" | "plus-lighter" | "saturation" | "screen" | "soft-light">>

export type ObjectFitValue = WithEscapeHatch<Globals | CssVars | OnlyKnown<"contain" | "cover" | "fill" | "none" | "scale-down">>

export type ObjectPositionValue = string | number | CssVars | AnyString

export type OpacityValue = WithEscapeHatch<Globals | TokenValue<"opacity"> | CssVars>

export type OverflowAnchorValue = string | number | CssVars | AnyString

export type OverflowBlockValue = WithEscapeHatch<Globals | CssVars | OnlyKnown<"auto" | "clip" | "hidden" | "scroll" | "visible">>

export type OverflowClipBoxValue = string | number | CssVars | AnyString

export type OverflowClipMarginValue = string | number | CssVars | AnyString

export type OverflowInlineValue = WithEscapeHatch<Globals | CssVars | OnlyKnown<"auto" | "clip" | "hidden" | "scroll" | "visible">>

export type OverflowValue = WithEscapeHatch<Globals | CssVars | OnlyKnown<"-moz-hidden-unscrollable" | "auto" | "clip" | "hidden" | "overlay" | "scroll" | "visible">>

export type OverflowWrapValue = WithEscapeHatch<Globals | CssVars | OnlyKnown<"anywhere" | "break-word" | "normal">>

export type OverflowXValue = WithEscapeHatch<Globals | CssVars | OnlyKnown<"-moz-hidden-unscrollable" | "auto" | "clip" | "hidden" | "overlay" | "scroll" | "visible">>

export type OverflowYValue = WithEscapeHatch<Globals | CssVars | OnlyKnown<"-moz-hidden-unscrollable" | "auto" | "clip" | "hidden" | "overlay" | "scroll" | "visible">>

export type OverscrollBehaviorBlockValue = string | number | CssVars | AnyString

export type OverscrollBehaviorInlineValue = string | number | CssVars | AnyString

export type OverscrollBehaviorValue = string | number | CssVars | AnyString

export type OverscrollBehaviorXValue = string | number | CssVars | AnyString

export type OverscrollBehaviorYValue = string | number | CssVars | AnyString

export type PositionValue = WithEscapeHatch<Globals | CssVars | OnlyKnown<"-webkit-sticky" | "absolute" | "fixed" | "relative" | "static" | "sticky">>

export type RadiiValue = WithEscapeHatch<Globals | TokenValue<"radii"> | CssVars>

export type RotateValue = WithEscapeHatch<Globals | "auto" | "auto-3d" | CssVars>

export type SaturateValue = string | number | CssVars | AnyString

export type ScaleValue = WithEscapeHatch<Globals | "auto" | CssVars>

export type ScaleXValue = string | number | CssVars | AnyString

export type ScaleYValue = string | number | CssVars | AnyString

export type ScrollBehaviorValue = WithEscapeHatch<Globals | CssVars | OnlyKnown<"auto" | "smooth">>

export type ScrollSnapAlignValue = string | number | CssVars | AnyString

export type ScrollSnapCoordinateValue = string | number | CssVars | AnyString

export type ScrollSnapDestinationValue = string | number | CssVars | AnyString

export type ScrollSnapPointsXValue = string | number | CssVars | AnyString

export type ScrollSnapPointsYValue = string | number | CssVars | AnyString

export type ScrollSnapStopValue = string | number | CssVars | AnyString

export type ScrollSnapStrictnessValue = WithEscapeHatch<Globals | "mandatory" | "proximity" | CssVars>

export type ScrollSnapTypeValue = WithEscapeHatch<Globals | "both" | "none" | "x" | "y" | CssVars>

export type ScrollSnapTypeXValue = string | number | CssVars | AnyString

export type ScrollSnapTypeYValue = string | number | CssVars | AnyString

export type ScrollTimelineAxisValue = string | number | CssVars | AnyString

export type ScrollTimelineNameValue = string | number | CssVars | AnyString

export type ScrollTimelineValue = string | number | CssVars | AnyString

export type ScrollbarGutterValue = string | number | CssVars | AnyString

export type ScrollbarValue = WithEscapeHatch<Globals | "hidden" | "visible" | CssVars>

export type SepiaValue = string | number | CssVars | AnyString

export type ShadowsValue = WithEscapeHatch<Globals | TokenValue<"shadows"> | CssVars>

export type SizesValue = WithEscapeHatch<DimensionGlobals | TokenValue<"sizes"> | CssVars>

export type SpacingValue = WithEscapeHatch<AutoGlobals | TokenValue<"spacing"> | CssVars>

export type SrOnlyValue = WithEscapeHatch<Globals | boolean | CssVars>

export type StrokeDasharrayValue = string | number | CssVars | AnyString

export type StrokeDashoffsetValue = string | number | CssVars | AnyString

export type StrokeLinecapValue = string | number | CssVars | AnyString

export type StrokeLinejoinValue = string | number | CssVars | AnyString

export type StrokeMiterlimitValue = string | number | CssVars | AnyString

export type StrokeOpacityValue = string | number | CssVars | AnyString

export type TableLayoutValue = string | number | CssVars | AnyString

export type TextAlignValue = string | number | CssVars | AnyString

export type TextDecorationStyleValue = string | number | CssVars | AnyString

export type TextDecorationThicknessValue = string | number | CssVars | AnyString

export type TextDecorationValue = string | number | CssVars | AnyString

export type TextGradientValue = WithEscapeHatch<Globals | "to-b" | "to-bl" | "to-br" | "to-l" | "to-r" | "to-t" | "to-tl" | "to-tr" | CssVars>

export type TextOverflowValue = string | number | CssVars | AnyString

export type TextSizeAdjustValue = string | number | CssVars | AnyString

export type TextStyleValue = WithEscapeHatch<Globals | "2xl" | "3xl" | "4xl" | "5xl" | "6xl" | "7xl" | "lg" | "md" | "sm" | "xl" | "xs" | CssVars>

export type TextTransformValue = string | number | CssVars | AnyString

export type TextUnderlineOffsetValue = string | number | CssVars | AnyString

export type TextWrapValue = WithEscapeHatch<Globals | "balance" | "nowrap" | "wrap" | CssVars>

export type TouchActionValue = WithEscapeHatch<Globals | CssVars | OnlyKnown<"-ms-manipulation" | "-ms-none" | "-ms-pan-x" | "-ms-pan-y" | "-ms-pinch-zoom" | "auto" | "manipulation" | "none" | "pan-down" | "pan-left" | "pan-right" | "pan-up" | "pan-x" | "pan-y" | "pinch-zoom">>

export type TransformBoxValue = WithEscapeHatch<Globals | CssVars | OnlyKnown<"border-box" | "content-box" | "fill-box" | "stroke-box" | "view-box">>

export type TransformOriginValue = string | number | CssVars | AnyString

export type TransformStyleValue = WithEscapeHatch<Globals | CssVars | OnlyKnown<"flat" | "preserve-3d">>

export type TransformValue = string | number | CssVars | AnyString

export type TransitionPropertyValue = WithEscapeHatch<Globals | "background" | "colors" | "common" | "position" | "size" | CssVars>

export type TransitionValue = WithEscapeHatch<Globals | "all" | "background" | "colors" | "common" | "opacity" | "position" | "shadow" | "size" | "transform" | CssVars>

export type TranslateValue = WithEscapeHatch<Globals | "auto" | "auto-3d" | CssVars>

export type TranslateXValue = WithEscapeHatch<Globals | "-1" | "-1/2" | "-1/3" | "-1/4" | "-10" | "-11" | "-12" | "-14" | "-16" | "-2" | "-2/3" | "-2/4" | "-20" | "-24" | "-28" | "-3" | "-3/4" | "-32" | "-36" | "-4" | "-40" | "-44" | "-48" | "-5" | "-52" | "-56" | "-6" | "-60" | "-64" | "-7" | "-72" | "-8" | "-80" | "-9" | "-96" | "-full" | "-safeBottom" | "-safeTop" | "-scrollGutter" | "0" | "0.-5" | "0.5" | "1" | "1.-5" | "1.5" | "1/2" | "1/3" | "1/4" | "10" | "10.-5" | "10.5" | "11" | "12" | "14" | "16" | "2" | "2.-5" | "2.5" | "2/3" | "2/4" | "20" | "24" | "28" | "3" | "3.-5" | "3.5" | "3/4" | "32" | "36" | "4" | "4.-5" | "4.5" | "40" | "44" | "48" | "5" | "5.-5" | "5.5" | "52" | "56" | "6" | "60" | "64" | "7" | "7.-5" | "7.5" | "72" | "8" | "80" | "9" | "9.-5" | "9.5" | "96" | "full" | "safeBottom" | "safeTop" | "scrollGutter" | CssVars>

export type TranslateYValue = WithEscapeHatch<Globals | "-1" | "-1/2" | "-1/3" | "-1/4" | "-10" | "-11" | "-12" | "-14" | "-16" | "-2" | "-2/3" | "-2/4" | "-20" | "-24" | "-28" | "-3" | "-3/4" | "-32" | "-36" | "-4" | "-40" | "-44" | "-48" | "-5" | "-52" | "-56" | "-6" | "-60" | "-64" | "-7" | "-72" | "-8" | "-80" | "-9" | "-96" | "-full" | "-safeBottom" | "-safeTop" | "-scrollGutter" | "0" | "0.-5" | "0.5" | "1" | "1.-5" | "1.5" | "1/2" | "1/3" | "1/4" | "10" | "10.-5" | "10.5" | "11" | "12" | "14" | "16" | "2" | "2.-5" | "2.5" | "2/3" | "2/4" | "20" | "24" | "28" | "3" | "3.-5" | "3.5" | "3/4" | "32" | "36" | "4" | "4.-5" | "4.5" | "40" | "44" | "48" | "5" | "5.-5" | "5.5" | "52" | "56" | "6" | "60" | "64" | "7" | "7.-5" | "7.5" | "72" | "8" | "80" | "9" | "9.-5" | "9.5" | "96" | "full" | "safeBottom" | "safeTop" | "scrollGutter" | CssVars>

export type TranslateZValue = WithEscapeHatch<Globals | "-1" | "-1/2" | "-1/3" | "-1/4" | "-10" | "-11" | "-12" | "-14" | "-16" | "-2" | "-2/3" | "-2/4" | "-20" | "-24" | "-28" | "-3" | "-3/4" | "-32" | "-36" | "-4" | "-40" | "-44" | "-48" | "-5" | "-52" | "-56" | "-6" | "-60" | "-64" | "-7" | "-72" | "-8" | "-80" | "-9" | "-96" | "-full" | "-safeBottom" | "-safeTop" | "-scrollGutter" | "0" | "0.-5" | "0.5" | "1" | "1.-5" | "1.5" | "1/2" | "1/3" | "1/4" | "10" | "10.-5" | "10.5" | "11" | "12" | "14" | "16" | "2" | "2.-5" | "2.5" | "2/3" | "2/4" | "20" | "24" | "28" | "3" | "3.-5" | "3.5" | "3/4" | "32" | "36" | "4" | "4.-5" | "4.5" | "40" | "44" | "48" | "5" | "5.-5" | "5.5" | "52" | "56" | "6" | "60" | "64" | "7" | "7.-5" | "7.5" | "72" | "8" | "80" | "9" | "9.-5" | "9.5" | "96" | "full" | "safeBottom" | "safeTop" | "scrollGutter" | CssVars>

export type TruncateValue = WithEscapeHatch<Globals | boolean | CssVars>

export type UserSelectValue = WithEscapeHatch<Globals | CssVars | OnlyKnown<"-moz-none" | "all" | "auto" | "none" | "text">>

export type VerticalAlignValue = string | number | CssVars | AnyString

export type VisibilityValue = WithEscapeHatch<Globals | CssVars | OnlyKnown<"collapse" | "hidden" | "visible">>

export type WidthValue = WithEscapeHatch<Globals | "0" | "0.5" | "1" | "1.5" | "1/12" | "1/2" | "1/3" | "1/4" | "1/5" | "1/6" | "10" | "10.5" | "10/12" | "11" | "11/12" | "12" | "14" | "16" | "2" | "2.5" | "2/12" | "2/3" | "2/4" | "2/5" | "2/6" | "20" | "24" | "28" | "2xl" | "2xs" | "3" | "3.5" | "3/12" | "3/4" | "3/5" | "3/6" | "32" | "36" | "3xl" | "4" | "4.5" | "4/12" | "4/5" | "4/6" | "40" | "44" | "48" | "4xl" | "5" | "5.5" | "5/12" | "5/6" | "52" | "56" | "5xl" | "6" | "6/12" | "60" | "64" | "6xl" | "7" | "7.5" | "7/12" | "72" | "7xl" | "8" | "8/12" | "80" | "8xl" | "9" | "9.5" | "9/12" | "96" | "auto" | "breakpoint-2xl" | "breakpoint-lg" | "breakpoint-md" | "breakpoint-sm" | "breakpoint-xl" | "fit" | "full" | "lg" | "max" | "md" | "min" | "prose" | "safeBottom" | "screen" | "scrollGutter" | "sm" | "xl" | "xs" | CssVars>

export type WordBreakValue = WithEscapeHatch<Globals | CssVars | OnlyKnown<"auto-phrase" | "break-all" | "break-word" | "keep-all" | "normal">>

export type ZIndexValue = WithEscapeHatch<AutoGlobals | TokenValue<"zIndex"> | CssVars>

export type AtRuleType = "media" | "layer" | "container" | "supports" | "page" | "scope" | "starting-style"

export type Selector = `${string}&` | `&${string}` | `@${AtRuleType}${string}`

export type AnySelector = Selector | string

export type Nested<P> = P & {
  [K in Selector]?: Nested<P>
} & {
  [K in AnySelector]?: Nested<P>
} & {
  [K in Condition]?: Nested<P>
}

export type Globals = "inherit" | "initial" | "revert" | "revert-layer" | "unset"

export type ColorGlobals = Globals | "currentColor" | "transparent"

export type DimensionGlobals = Globals | "auto" | "fit-content" | "max-content" | "min-content"

export type AutoGlobals = Globals | "auto"

export type CssValue = Globals | (string & {}) | number

export interface CssProperties {
  WebkitAppearance?: ConditionalValue<CssValue>
  WebkitBorderBefore?: ConditionalValue<CssValue>
  WebkitBorderBeforeColor?: ConditionalValue<CssValue>
  WebkitBorderBeforeStyle?: ConditionalValue<CssValue>
  WebkitBorderBeforeWidth?: ConditionalValue<CssValue>
  WebkitBoxReflect?: ConditionalValue<CssValue>
  WebkitLineClamp?: ConditionalValue<CssValue>
  WebkitMask?: ConditionalValue<CssValue>
  WebkitMaskAttachment?: ConditionalValue<CssValue>
  WebkitMaskClip?: ConditionalValue<CssValue>
  WebkitMaskComposite?: ConditionalValue<CssValue>
  WebkitMaskImage?: ConditionalValue<CssValue>
  WebkitMaskOrigin?: ConditionalValue<CssValue>
  WebkitMaskPosition?: ConditionalValue<CssValue>
  WebkitMaskPositionX?: ConditionalValue<CssValue>
  WebkitMaskPositionY?: ConditionalValue<CssValue>
  WebkitMaskRepeat?: ConditionalValue<CssValue>
  WebkitMaskRepeatX?: ConditionalValue<CssValue>
  WebkitMaskRepeatY?: ConditionalValue<CssValue>
  WebkitMaskSize?: ConditionalValue<CssValue>
  WebkitOverflowScrolling?: ConditionalValue<CssValue>
  WebkitTapHighlightColor?: ConditionalValue<CssValue>
  WebkitTextFillColor?: ConditionalValue<CssValue>
  WebkitTextStroke?: ConditionalValue<CssValue>
  WebkitTextStrokeColor?: ConditionalValue<CssValue>
  WebkitTextStrokeWidth?: ConditionalValue<CssValue>
  WebkitTouchCallout?: ConditionalValue<CssValue>
  WebkitUserModify?: ConditionalValue<CssValue>
  WebkitUserSelect?: ConditionalValue<CssValue>
  accentColor?: ConditionalValue<CssValue>
  alignContent?: ConditionalValue<CssValue>
  alignItems?: ConditionalValue<CssValue>
  alignSelf?: ConditionalValue<CssValue>
  alignTracks?: ConditionalValue<CssValue>
  alignmentBaseline?: ConditionalValue<CssValue>
  all?: ConditionalValue<CssValue>
  anchorName?: ConditionalValue<CssValue>
  anchorScope?: ConditionalValue<CssValue>
  animation?: ConditionalValue<CssValue>
  animationComposition?: ConditionalValue<CssValue>
  animationDelay?: ConditionalValue<CssValue>
  animationDirection?: ConditionalValue<CssValue>
  animationDuration?: ConditionalValue<CssValue>
  animationFillMode?: ConditionalValue<CssValue>
  animationIterationCount?: ConditionalValue<CssValue>
  animationName?: ConditionalValue<CssValue>
  animationPlayState?: ConditionalValue<CssValue>
  animationRange?: ConditionalValue<CssValue>
  animationRangeEnd?: ConditionalValue<CssValue>
  animationRangeStart?: ConditionalValue<CssValue>
  animationTimeline?: ConditionalValue<CssValue>
  animationTimingFunction?: ConditionalValue<CssValue>
  appearance?: ConditionalValue<CssValue>
  aspectRatio?: ConditionalValue<CssValue>
  backdropFilter?: ConditionalValue<CssValue>
  backfaceVisibility?: ConditionalValue<CssValue>
  background?: ConditionalValue<CssValue>
  backgroundAttachment?: ConditionalValue<CssValue>
  backgroundBlendMode?: ConditionalValue<CssValue>
  backgroundClip?: ConditionalValue<CssValue>
  backgroundColor?: ConditionalValue<CssValue>
  backgroundImage?: ConditionalValue<CssValue>
  backgroundOrigin?: ConditionalValue<CssValue>
  backgroundPosition?: ConditionalValue<CssValue>
  backgroundPositionX?: ConditionalValue<CssValue>
  backgroundPositionY?: ConditionalValue<CssValue>
  backgroundRepeat?: ConditionalValue<CssValue>
  backgroundSize?: ConditionalValue<CssValue>
  baselineShift?: ConditionalValue<CssValue>
  blockSize?: ConditionalValue<CssValue>
  border?: ConditionalValue<CssValue>
  borderBlock?: ConditionalValue<CssValue>
  borderBlockColor?: ConditionalValue<CssValue>
  borderBlockEnd?: ConditionalValue<CssValue>
  borderBlockEndColor?: ConditionalValue<CssValue>
  borderBlockEndStyle?: ConditionalValue<CssValue>
  borderBlockEndWidth?: ConditionalValue<CssValue>
  borderBlockStart?: ConditionalValue<CssValue>
  borderBlockStartColor?: ConditionalValue<CssValue>
  borderBlockStartStyle?: ConditionalValue<CssValue>
  borderBlockStartWidth?: ConditionalValue<CssValue>
  borderBlockStyle?: ConditionalValue<CssValue>
  borderBlockWidth?: ConditionalValue<CssValue>
  borderBottom?: ConditionalValue<CssValue>
  borderBottomColor?: ConditionalValue<CssValue>
  borderBottomLeftRadius?: ConditionalValue<CssValue>
  borderBottomRightRadius?: ConditionalValue<CssValue>
  borderBottomStyle?: ConditionalValue<CssValue>
  borderBottomWidth?: ConditionalValue<CssValue>
  borderCollapse?: ConditionalValue<CssValue>
  borderColor?: ConditionalValue<CssValue>
  borderEndEndRadius?: ConditionalValue<CssValue>
  borderEndStartRadius?: ConditionalValue<CssValue>
  borderImage?: ConditionalValue<CssValue>
  borderImageOutset?: ConditionalValue<CssValue>
  borderImageRepeat?: ConditionalValue<CssValue>
  borderImageSlice?: ConditionalValue<CssValue>
  borderImageSource?: ConditionalValue<CssValue>
  borderImageWidth?: ConditionalValue<CssValue>
  borderInline?: ConditionalValue<CssValue>
  borderInlineColor?: ConditionalValue<CssValue>
  borderInlineEnd?: ConditionalValue<CssValue>
  borderInlineEndColor?: ConditionalValue<CssValue>
  borderInlineEndStyle?: ConditionalValue<CssValue>
  borderInlineEndWidth?: ConditionalValue<CssValue>
  borderInlineStart?: ConditionalValue<CssValue>
  borderInlineStartColor?: ConditionalValue<CssValue>
  borderInlineStartStyle?: ConditionalValue<CssValue>
  borderInlineStartWidth?: ConditionalValue<CssValue>
  borderInlineStyle?: ConditionalValue<CssValue>
  borderInlineWidth?: ConditionalValue<CssValue>
  borderLeft?: ConditionalValue<CssValue>
  borderLeftColor?: ConditionalValue<CssValue>
  borderLeftStyle?: ConditionalValue<CssValue>
  borderLeftWidth?: ConditionalValue<CssValue>
  borderRadius?: ConditionalValue<CssValue>
  borderRight?: ConditionalValue<CssValue>
  borderRightColor?: ConditionalValue<CssValue>
  borderRightStyle?: ConditionalValue<CssValue>
  borderRightWidth?: ConditionalValue<CssValue>
  borderSpacing?: ConditionalValue<CssValue>
  borderStartEndRadius?: ConditionalValue<CssValue>
  borderStartStartRadius?: ConditionalValue<CssValue>
  borderStyle?: ConditionalValue<CssValue>
  borderTop?: ConditionalValue<CssValue>
  borderTopColor?: ConditionalValue<CssValue>
  borderTopLeftRadius?: ConditionalValue<CssValue>
  borderTopRightRadius?: ConditionalValue<CssValue>
  borderTopStyle?: ConditionalValue<CssValue>
  borderTopWidth?: ConditionalValue<CssValue>
  borderWidth?: ConditionalValue<CssValue>
  bottom?: ConditionalValue<CssValue>
  boxAlign?: ConditionalValue<CssValue>
  boxDecorationBreak?: ConditionalValue<CssValue>
  boxDirection?: ConditionalValue<CssValue>
  boxFlex?: ConditionalValue<CssValue>
  boxFlexGroup?: ConditionalValue<CssValue>
  boxLines?: ConditionalValue<CssValue>
  boxOrdinalGroup?: ConditionalValue<CssValue>
  boxOrient?: ConditionalValue<CssValue>
  boxPack?: ConditionalValue<CssValue>
  boxShadow?: ConditionalValue<CssValue>
  boxSizing?: ConditionalValue<CssValue>
  breakAfter?: ConditionalValue<CssValue>
  breakBefore?: ConditionalValue<CssValue>
  breakInside?: ConditionalValue<CssValue>
  captionSide?: ConditionalValue<CssValue>
  caret?: ConditionalValue<CssValue>
  caretColor?: ConditionalValue<CssValue>
  caretShape?: ConditionalValue<CssValue>
  clear?: ConditionalValue<CssValue>
  clip?: ConditionalValue<CssValue>
  clipPath?: ConditionalValue<CssValue>
  clipRule?: ConditionalValue<CssValue>
  color?: ConditionalValue<CssValue>
  colorInterpolation?: ConditionalValue<CssValue>
  colorInterpolationFilters?: ConditionalValue<CssValue>
  colorRendering?: ConditionalValue<CssValue>
  colorScheme?: ConditionalValue<CssValue>
  columnCount?: ConditionalValue<CssValue>
  columnFill?: ConditionalValue<CssValue>
  columnGap?: ConditionalValue<CssValue>
  columnRule?: ConditionalValue<CssValue>
  columnRuleColor?: ConditionalValue<CssValue>
  columnRuleStyle?: ConditionalValue<CssValue>
  columnRuleWidth?: ConditionalValue<CssValue>
  columnSpan?: ConditionalValue<CssValue>
  columnWidth?: ConditionalValue<CssValue>
  columns?: ConditionalValue<CssValue>
  contain?: ConditionalValue<CssValue>
  containIntrinsicBlockSize?: ConditionalValue<CssValue>
  containIntrinsicHeight?: ConditionalValue<CssValue>
  containIntrinsicInlineSize?: ConditionalValue<CssValue>
  containIntrinsicSize?: ConditionalValue<CssValue>
  containIntrinsicWidth?: ConditionalValue<CssValue>
  container?: ConditionalValue<CssValue>
  containerName?: ConditionalValue<CssValue>
  containerType?: ConditionalValue<CssValue>
  content?: ConditionalValue<CssValue>
  contentVisibility?: ConditionalValue<CssValue>
  counterIncrement?: ConditionalValue<CssValue>
  counterReset?: ConditionalValue<CssValue>
  counterSet?: ConditionalValue<CssValue>
  cursor?: ConditionalValue<CssValue>
  cx?: ConditionalValue<CssValue>
  cy?: ConditionalValue<CssValue>
  d?: ConditionalValue<CssValue>
  direction?: ConditionalValue<CssValue>
  display?: ConditionalValue<CssValue>
  dominantBaseline?: ConditionalValue<CssValue>
  emptyCells?: ConditionalValue<CssValue>
  fieldSizing?: ConditionalValue<CssValue>
  fill?: ConditionalValue<CssValue>
  fillOpacity?: ConditionalValue<CssValue>
  fillRule?: ConditionalValue<CssValue>
  filter?: ConditionalValue<CssValue>
  flex?: ConditionalValue<CssValue>
  flexBasis?: ConditionalValue<CssValue>
  flexDirection?: ConditionalValue<CssValue>
  flexFlow?: ConditionalValue<CssValue>
  flexGrow?: ConditionalValue<CssValue>
  flexShrink?: ConditionalValue<CssValue>
  flexWrap?: ConditionalValue<CssValue>
  float?: ConditionalValue<CssValue>
  floodColor?: ConditionalValue<CssValue>
  floodOpacity?: ConditionalValue<CssValue>
  font?: ConditionalValue<CssValue>
  fontFamily?: ConditionalValue<CssValue>
  fontFeatureSettings?: ConditionalValue<CssValue>
  fontKerning?: ConditionalValue<CssValue>
  fontLanguageOverride?: ConditionalValue<CssValue>
  fontOpticalSizing?: ConditionalValue<CssValue>
  fontPalette?: ConditionalValue<CssValue>
  fontSize?: ConditionalValue<CssValue>
  fontSizeAdjust?: ConditionalValue<CssValue>
  fontSmooth?: ConditionalValue<CssValue>
  fontStretch?: ConditionalValue<CssValue>
  fontStyle?: ConditionalValue<CssValue>
  fontSynthesis?: ConditionalValue<CssValue>
  fontSynthesisPosition?: ConditionalValue<CssValue>
  fontSynthesisSmallCaps?: ConditionalValue<CssValue>
  fontSynthesisStyle?: ConditionalValue<CssValue>
  fontSynthesisWeight?: ConditionalValue<CssValue>
  fontVariant?: ConditionalValue<CssValue>
  fontVariantAlternates?: ConditionalValue<CssValue>
  fontVariantCaps?: ConditionalValue<CssValue>
  fontVariantEastAsian?: ConditionalValue<CssValue>
  fontVariantEmoji?: ConditionalValue<CssValue>
  fontVariantLigatures?: ConditionalValue<CssValue>
  fontVariantNumeric?: ConditionalValue<CssValue>
  fontVariantPosition?: ConditionalValue<CssValue>
  fontVariationSettings?: ConditionalValue<CssValue>
  fontWeight?: ConditionalValue<CssValue>
  fontWidth?: ConditionalValue<CssValue>
  forcedColorAdjust?: ConditionalValue<CssValue>
  gap?: ConditionalValue<CssValue>
  glyphOrientationVertical?: ConditionalValue<CssValue>
  grid?: ConditionalValue<CssValue>
  gridArea?: ConditionalValue<CssValue>
  gridAutoColumns?: ConditionalValue<CssValue>
  gridAutoFlow?: ConditionalValue<CssValue>
  gridAutoRows?: ConditionalValue<CssValue>
  gridColumn?: ConditionalValue<CssValue>
  gridColumnEnd?: ConditionalValue<CssValue>
  gridColumnGap?: ConditionalValue<CssValue>
  gridColumnStart?: ConditionalValue<CssValue>
  gridGap?: ConditionalValue<CssValue>
  gridRow?: ConditionalValue<CssValue>
  gridRowEnd?: ConditionalValue<CssValue>
  gridRowGap?: ConditionalValue<CssValue>
  gridRowStart?: ConditionalValue<CssValue>
  gridTemplate?: ConditionalValue<CssValue>
  gridTemplateAreas?: ConditionalValue<CssValue>
  gridTemplateColumns?: ConditionalValue<CssValue>
  gridTemplateRows?: ConditionalValue<CssValue>
  hangingPunctuation?: ConditionalValue<CssValue>
  height?: ConditionalValue<CssValue>
  hyphenateCharacter?: ConditionalValue<CssValue>
  hyphenateLimitChars?: ConditionalValue<CssValue>
  hyphens?: ConditionalValue<CssValue>
  imageOrientation?: ConditionalValue<CssValue>
  imageRendering?: ConditionalValue<CssValue>
  imageResolution?: ConditionalValue<CssValue>
  imeMode?: ConditionalValue<CssValue>
  initialLetter?: ConditionalValue<CssValue>
  initialLetterAlign?: ConditionalValue<CssValue>
  inlineSize?: ConditionalValue<CssValue>
  inset?: ConditionalValue<CssValue>
  insetBlock?: ConditionalValue<CssValue>
  insetBlockEnd?: ConditionalValue<CssValue>
  insetBlockStart?: ConditionalValue<CssValue>
  insetInline?: ConditionalValue<CssValue>
  insetInlineEnd?: ConditionalValue<CssValue>
  insetInlineStart?: ConditionalValue<CssValue>
  interpolateSize?: ConditionalValue<CssValue>
  isolation?: ConditionalValue<CssValue>
  justifyContent?: ConditionalValue<CssValue>
  justifyItems?: ConditionalValue<CssValue>
  justifySelf?: ConditionalValue<CssValue>
  justifyTracks?: ConditionalValue<CssValue>
  left?: ConditionalValue<CssValue>
  letterSpacing?: ConditionalValue<CssValue>
  lightingColor?: ConditionalValue<CssValue>
  lineBreak?: ConditionalValue<CssValue>
  lineClamp?: ConditionalValue<CssValue>
  lineHeight?: ConditionalValue<CssValue>
  lineHeightStep?: ConditionalValue<CssValue>
  listStyle?: ConditionalValue<CssValue>
  listStyleImage?: ConditionalValue<CssValue>
  listStylePosition?: ConditionalValue<CssValue>
  listStyleType?: ConditionalValue<CssValue>
  margin?: ConditionalValue<CssValue>
  marginBlock?: ConditionalValue<CssValue>
  marginBlockEnd?: ConditionalValue<CssValue>
  marginBlockStart?: ConditionalValue<CssValue>
  marginBottom?: ConditionalValue<CssValue>
  marginInline?: ConditionalValue<CssValue>
  marginInlineEnd?: ConditionalValue<CssValue>
  marginInlineStart?: ConditionalValue<CssValue>
  marginLeft?: ConditionalValue<CssValue>
  marginRight?: ConditionalValue<CssValue>
  marginTop?: ConditionalValue<CssValue>
  marginTrim?: ConditionalValue<CssValue>
  marker?: ConditionalValue<CssValue>
  markerEnd?: ConditionalValue<CssValue>
  markerMid?: ConditionalValue<CssValue>
  markerStart?: ConditionalValue<CssValue>
  mask?: ConditionalValue<CssValue>
  maskBorder?: ConditionalValue<CssValue>
  maskBorderMode?: ConditionalValue<CssValue>
  maskBorderOutset?: ConditionalValue<CssValue>
  maskBorderRepeat?: ConditionalValue<CssValue>
  maskBorderSlice?: ConditionalValue<CssValue>
  maskBorderSource?: ConditionalValue<CssValue>
  maskBorderWidth?: ConditionalValue<CssValue>
  maskClip?: ConditionalValue<CssValue>
  maskComposite?: ConditionalValue<CssValue>
  maskImage?: ConditionalValue<CssValue>
  maskMode?: ConditionalValue<CssValue>
  maskOrigin?: ConditionalValue<CssValue>
  maskPosition?: ConditionalValue<CssValue>
  maskRepeat?: ConditionalValue<CssValue>
  maskSize?: ConditionalValue<CssValue>
  maskType?: ConditionalValue<CssValue>
  masonryAutoFlow?: ConditionalValue<CssValue>
  mathDepth?: ConditionalValue<CssValue>
  mathShift?: ConditionalValue<CssValue>
  mathStyle?: ConditionalValue<CssValue>
  maxBlockSize?: ConditionalValue<CssValue>
  maxHeight?: ConditionalValue<CssValue>
  maxInlineSize?: ConditionalValue<CssValue>
  maxLines?: ConditionalValue<CssValue>
  maxWidth?: ConditionalValue<CssValue>
  minBlockSize?: ConditionalValue<CssValue>
  minHeight?: ConditionalValue<CssValue>
  minInlineSize?: ConditionalValue<CssValue>
  minWidth?: ConditionalValue<CssValue>
  mixBlendMode?: ConditionalValue<CssValue>
  objectFit?: ConditionalValue<CssValue>
  objectPosition?: ConditionalValue<CssValue>
  objectViewBox?: ConditionalValue<CssValue>
  offset?: ConditionalValue<CssValue>
  offsetAnchor?: ConditionalValue<CssValue>
  offsetDistance?: ConditionalValue<CssValue>
  offsetPath?: ConditionalValue<CssValue>
  offsetPosition?: ConditionalValue<CssValue>
  offsetRotate?: ConditionalValue<CssValue>
  opacity?: ConditionalValue<CssValue>
  order?: ConditionalValue<CssValue>
  orphans?: ConditionalValue<CssValue>
  outline?: ConditionalValue<CssValue>
  outlineColor?: ConditionalValue<CssValue>
  outlineOffset?: ConditionalValue<CssValue>
  outlineStyle?: ConditionalValue<CssValue>
  outlineWidth?: ConditionalValue<CssValue>
  overflow?: ConditionalValue<CssValue>
  overflowAnchor?: ConditionalValue<CssValue>
  overflowBlock?: ConditionalValue<CssValue>
  overflowClipBox?: ConditionalValue<CssValue>
  overflowClipMargin?: ConditionalValue<CssValue>
  overflowInline?: ConditionalValue<CssValue>
  overflowWrap?: ConditionalValue<CssValue>
  overflowX?: ConditionalValue<CssValue>
  overflowY?: ConditionalValue<CssValue>
  overlay?: ConditionalValue<CssValue>
  overscrollBehavior?: ConditionalValue<CssValue>
  overscrollBehaviorBlock?: ConditionalValue<CssValue>
  overscrollBehaviorInline?: ConditionalValue<CssValue>
  overscrollBehaviorX?: ConditionalValue<CssValue>
  overscrollBehaviorY?: ConditionalValue<CssValue>
  padding?: ConditionalValue<CssValue>
  paddingBlock?: ConditionalValue<CssValue>
  paddingBlockEnd?: ConditionalValue<CssValue>
  paddingBlockStart?: ConditionalValue<CssValue>
  paddingBottom?: ConditionalValue<CssValue>
  paddingInline?: ConditionalValue<CssValue>
  paddingInlineEnd?: ConditionalValue<CssValue>
  paddingInlineStart?: ConditionalValue<CssValue>
  paddingLeft?: ConditionalValue<CssValue>
  paddingRight?: ConditionalValue<CssValue>
  paddingTop?: ConditionalValue<CssValue>
  page?: ConditionalValue<CssValue>
  pageBreakAfter?: ConditionalValue<CssValue>
  pageBreakBefore?: ConditionalValue<CssValue>
  pageBreakInside?: ConditionalValue<CssValue>
  paintOrder?: ConditionalValue<CssValue>
  perspective?: ConditionalValue<CssValue>
  perspectiveOrigin?: ConditionalValue<CssValue>
  placeContent?: ConditionalValue<CssValue>
  placeItems?: ConditionalValue<CssValue>
  placeSelf?: ConditionalValue<CssValue>
  pointerEvents?: ConditionalValue<CssValue>
  position?: ConditionalValue<CssValue>
  positionAnchor?: ConditionalValue<CssValue>
  positionArea?: ConditionalValue<CssValue>
  positionTry?: ConditionalValue<CssValue>
  positionTryFallbacks?: ConditionalValue<CssValue>
  positionTryOrder?: ConditionalValue<CssValue>
  positionVisibility?: ConditionalValue<CssValue>
  printColorAdjust?: ConditionalValue<CssValue>
  quotes?: ConditionalValue<CssValue>
  r?: ConditionalValue<CssValue>
  resize?: ConditionalValue<CssValue>
  right?: ConditionalValue<CssValue>
  rotate?: ConditionalValue<CssValue>
  rowGap?: ConditionalValue<CssValue>
  rubyAlign?: ConditionalValue<CssValue>
  rubyMerge?: ConditionalValue<CssValue>
  rubyOverhang?: ConditionalValue<CssValue>
  rubyPosition?: ConditionalValue<CssValue>
  rx?: ConditionalValue<CssValue>
  ry?: ConditionalValue<CssValue>
  scale?: ConditionalValue<CssValue>
  scrollBehavior?: ConditionalValue<CssValue>
  scrollInitialTarget?: ConditionalValue<CssValue>
  scrollMargin?: ConditionalValue<CssValue>
  scrollMarginBlock?: ConditionalValue<CssValue>
  scrollMarginBlockEnd?: ConditionalValue<CssValue>
  scrollMarginBlockStart?: ConditionalValue<CssValue>
  scrollMarginBottom?: ConditionalValue<CssValue>
  scrollMarginInline?: ConditionalValue<CssValue>
  scrollMarginInlineEnd?: ConditionalValue<CssValue>
  scrollMarginInlineStart?: ConditionalValue<CssValue>
  scrollMarginLeft?: ConditionalValue<CssValue>
  scrollMarginRight?: ConditionalValue<CssValue>
  scrollMarginTop?: ConditionalValue<CssValue>
  scrollPadding?: ConditionalValue<CssValue>
  scrollPaddingBlock?: ConditionalValue<CssValue>
  scrollPaddingBlockEnd?: ConditionalValue<CssValue>
  scrollPaddingBlockStart?: ConditionalValue<CssValue>
  scrollPaddingBottom?: ConditionalValue<CssValue>
  scrollPaddingInline?: ConditionalValue<CssValue>
  scrollPaddingInlineEnd?: ConditionalValue<CssValue>
  scrollPaddingInlineStart?: ConditionalValue<CssValue>
  scrollPaddingLeft?: ConditionalValue<CssValue>
  scrollPaddingRight?: ConditionalValue<CssValue>
  scrollPaddingTop?: ConditionalValue<CssValue>
  scrollSnapAlign?: ConditionalValue<CssValue>
  scrollSnapCoordinate?: ConditionalValue<CssValue>
  scrollSnapDestination?: ConditionalValue<CssValue>
  scrollSnapPointsX?: ConditionalValue<CssValue>
  scrollSnapPointsY?: ConditionalValue<CssValue>
  scrollSnapStop?: ConditionalValue<CssValue>
  scrollSnapType?: ConditionalValue<CssValue>
  scrollSnapTypeX?: ConditionalValue<CssValue>
  scrollSnapTypeY?: ConditionalValue<CssValue>
  scrollTimeline?: ConditionalValue<CssValue>
  scrollTimelineAxis?: ConditionalValue<CssValue>
  scrollTimelineName?: ConditionalValue<CssValue>
  scrollbarColor?: ConditionalValue<CssValue>
  scrollbarGutter?: ConditionalValue<CssValue>
  scrollbarWidth?: ConditionalValue<CssValue>
  shapeImageThreshold?: ConditionalValue<CssValue>
  shapeMargin?: ConditionalValue<CssValue>
  shapeOutside?: ConditionalValue<CssValue>
  shapeRendering?: ConditionalValue<CssValue>
  speakAs?: ConditionalValue<CssValue>
  stopColor?: ConditionalValue<CssValue>
  stopOpacity?: ConditionalValue<CssValue>
  stroke?: ConditionalValue<CssValue>
  strokeColor?: ConditionalValue<CssValue>
  strokeDasharray?: ConditionalValue<CssValue>
  strokeDashoffset?: ConditionalValue<CssValue>
  strokeLinecap?: ConditionalValue<CssValue>
  strokeLinejoin?: ConditionalValue<CssValue>
  strokeMiterlimit?: ConditionalValue<CssValue>
  strokeOpacity?: ConditionalValue<CssValue>
  strokeWidth?: ConditionalValue<CssValue>
  tabSize?: ConditionalValue<CssValue>
  tableLayout?: ConditionalValue<CssValue>
  textAlign?: ConditionalValue<CssValue>
  textAlignLast?: ConditionalValue<CssValue>
  textAnchor?: ConditionalValue<CssValue>
  textAutospace?: ConditionalValue<CssValue>
  textBox?: ConditionalValue<CssValue>
  textBoxEdge?: ConditionalValue<CssValue>
  textBoxTrim?: ConditionalValue<CssValue>
  textCombineUpright?: ConditionalValue<CssValue>
  textDecoration?: ConditionalValue<CssValue>
  textDecorationColor?: ConditionalValue<CssValue>
  textDecorationLine?: ConditionalValue<CssValue>
  textDecorationSkip?: ConditionalValue<CssValue>
  textDecorationSkipInk?: ConditionalValue<CssValue>
  textDecorationStyle?: ConditionalValue<CssValue>
  textDecorationThickness?: ConditionalValue<CssValue>
  textEmphasis?: ConditionalValue<CssValue>
  textEmphasisColor?: ConditionalValue<CssValue>
  textEmphasisPosition?: ConditionalValue<CssValue>
  textEmphasisStyle?: ConditionalValue<CssValue>
  textIndent?: ConditionalValue<CssValue>
  textJustify?: ConditionalValue<CssValue>
  textOrientation?: ConditionalValue<CssValue>
  textOverflow?: ConditionalValue<CssValue>
  textRendering?: ConditionalValue<CssValue>
  textShadow?: ConditionalValue<CssValue>
  textSizeAdjust?: ConditionalValue<CssValue>
  textSpacingTrim?: ConditionalValue<CssValue>
  textTransform?: ConditionalValue<CssValue>
  textUnderlineOffset?: ConditionalValue<CssValue>
  textUnderlinePosition?: ConditionalValue<CssValue>
  textWrap?: ConditionalValue<CssValue>
  textWrapMode?: ConditionalValue<CssValue>
  textWrapStyle?: ConditionalValue<CssValue>
  timelineScope?: ConditionalValue<CssValue>
  top?: ConditionalValue<CssValue>
  touchAction?: ConditionalValue<CssValue>
  transform?: ConditionalValue<CssValue>
  transformBox?: ConditionalValue<CssValue>
  transformOrigin?: ConditionalValue<CssValue>
  transformStyle?: ConditionalValue<CssValue>
  transition?: ConditionalValue<CssValue>
  transitionBehavior?: ConditionalValue<CssValue>
  transitionDelay?: ConditionalValue<CssValue>
  transitionDuration?: ConditionalValue<CssValue>
  transitionProperty?: ConditionalValue<CssValue>
  transitionTimingFunction?: ConditionalValue<CssValue>
  translate?: ConditionalValue<CssValue>
  unicodeBidi?: ConditionalValue<CssValue>
  userSelect?: ConditionalValue<CssValue>
  vectorEffect?: ConditionalValue<CssValue>
  verticalAlign?: ConditionalValue<CssValue>
  viewTimeline?: ConditionalValue<CssValue>
  viewTimelineAxis?: ConditionalValue<CssValue>
  viewTimelineInset?: ConditionalValue<CssValue>
  viewTimelineName?: ConditionalValue<CssValue>
  viewTransitionClass?: ConditionalValue<CssValue>
  viewTransitionName?: ConditionalValue<CssValue>
  visibility?: ConditionalValue<CssValue>
  whiteSpace?: ConditionalValue<CssValue>
  whiteSpaceCollapse?: ConditionalValue<CssValue>
  widows?: ConditionalValue<CssValue>
  width?: ConditionalValue<CssValue>
  willChange?: ConditionalValue<CssValue>
  wordBreak?: ConditionalValue<CssValue>
  wordSpacing?: ConditionalValue<CssValue>
  wordWrap?: ConditionalValue<CssValue>
  writingMode?: ConditionalValue<CssValue>
  x?: ConditionalValue<CssValue>
  y?: ConditionalValue<CssValue>
  zIndex?: ConditionalValue<CssValue>
  zoom?: ConditionalValue<CssValue>
}

export interface SystemProperties extends CssProperties {
  WebkitTextFillColor?: ConditionalValue<ColorsValue>
  accentColor?: ConditionalValue<ColorsValue>
  alignContent?: ConditionalValue<AlignContentValue>
  alignItems?: ConditionalValue<AlignItemsValue>
  alignSelf?: ConditionalValue<AlignSelfValue>
  animation?: ConditionalValue<AnimationsValue>
  animationComposition?: ConditionalValue<AnimationCompositionValue>
  animationDelay?: ConditionalValue<DurationsValue>
  animationDirection?: ConditionalValue<AnimationDirectionValue>
  animationDuration?: ConditionalValue<DurationsValue>
  animationFillMode?: ConditionalValue<AnimationFillModeValue>
  animationIterationCount?: ConditionalValue<AnimationIterationCountValue>
  animationName?: ConditionalValue<KeyframesValue>
  animationPlayState?: ConditionalValue<AnimationPlayStateValue>
  animationRange?: ConditionalValue<AnimationRangeValue>
  animationRangeEnd?: ConditionalValue<AnimationRangeEndValue>
  animationRangeStart?: ConditionalValue<AnimationRangeStartValue>
  animationState?: ConditionalValue<AnimationStateValue>
  animationTimeline?: ConditionalValue<AnimationTimelineValue>
  animationTimingFunction?: ConditionalValue<EasingsValue>
  appearance?: ConditionalValue<AppearanceValue>
  aspectRatio?: ConditionalValue<AspectRatiosValue>
  backdropBlur?: ConditionalValue<BlursValue>
  backdropBrightness?: ConditionalValue<BackdropBrightnessValue>
  backdropContrast?: ConditionalValue<BackdropContrastValue>
  backdropFilter?: ConditionalValue<BackdropFilterValue>
  backdropGrayscale?: ConditionalValue<BackdropGrayscaleValue>
  backdropHueRotate?: ConditionalValue<BackdropHueRotateValue>
  backdropInvert?: ConditionalValue<BackdropInvertValue>
  backdropOpacity?: ConditionalValue<BackdropOpacityValue>
  backdropSaturate?: ConditionalValue<BackdropSaturateValue>
  backdropSepia?: ConditionalValue<BackdropSepiaValue>
  backfaceVisibility?: ConditionalValue<BackfaceVisibilityValue>
  background?: ConditionalValue<ColorsValue>
  backgroundAttachment?: ConditionalValue<BackgroundAttachmentValue>
  backgroundBlendMode?: ConditionalValue<BackgroundBlendModeValue>
  backgroundClip?: ConditionalValue<BackgroundClipValue>
  backgroundColor?: ConditionalValue<ColorsValue>
  backgroundConic?: ConditionalValue<BackgroundConicValue>
  backgroundGradient?: ConditionalValue<BackgroundGradientValue>
  backgroundImage?: ConditionalValue<AssetsValue>
  backgroundLinear?: ConditionalValue<BackgroundLinearValue>
  backgroundOrigin?: ConditionalValue<BackgroundOriginValue>
  backgroundPosition?: ConditionalValue<BackgroundPositionValue>
  backgroundPositionX?: ConditionalValue<BackgroundPositionXValue>
  backgroundPositionY?: ConditionalValue<BackgroundPositionYValue>
  backgroundRadial?: ConditionalValue<GradientsValue>
  backgroundRepeat?: ConditionalValue<BackgroundRepeatValue>
  backgroundSize?: ConditionalValue<BackgroundSizeValue>
  bg?: ConditionalValue<ColorsValue>
  bgAttachment?: ConditionalValue<BackgroundAttachmentValue>
  bgBlendMode?: ConditionalValue<BackgroundBlendModeValue>
  bgClip?: ConditionalValue<BackgroundClipValue>
  bgColor?: ConditionalValue<ColorsValue>
  bgConic?: ConditionalValue<BackgroundConicValue>
  bgGradient?: ConditionalValue<BackgroundGradientValue>
  bgImage?: ConditionalValue<AssetsValue>
  bgLinear?: ConditionalValue<BackgroundLinearValue>
  bgOrigin?: ConditionalValue<BackgroundOriginValue>
  bgPosition?: ConditionalValue<BackgroundPositionValue>
  bgPositionX?: ConditionalValue<BackgroundPositionXValue>
  bgPositionY?: ConditionalValue<BackgroundPositionYValue>
  bgRadial?: ConditionalValue<GradientsValue>
  bgRepeat?: ConditionalValue<BackgroundRepeatValue>
  bgSize?: ConditionalValue<BackgroundSizeValue>
  blockSize?: ConditionalValue<BlockSizeValue>
  blur?: ConditionalValue<BlursValue>
  border?: ConditionalValue<BordersValue>
  borderBlock?: ConditionalValue<BordersValue>
  borderBlockColor?: ConditionalValue<ColorsValue>
  borderBlockEnd?: ConditionalValue<BordersValue>
  borderBlockEndColor?: ConditionalValue<ColorsValue>
  borderBlockEndWidth?: ConditionalValue<BorderWidthsValue>
  borderBlockStart?: ConditionalValue<BordersValue>
  borderBlockStartColor?: ConditionalValue<ColorsValue>
  borderBlockStartWidth?: ConditionalValue<BorderWidthsValue>
  borderBlockWidth?: ConditionalValue<BorderWidthsValue>
  borderBottom?: ConditionalValue<BordersValue>
  borderBottomColor?: ConditionalValue<ColorsValue>
  borderBottomLeftRadius?: ConditionalValue<RadiiValue>
  borderBottomRadius?: ConditionalValue<RadiiValue>
  borderBottomRightRadius?: ConditionalValue<RadiiValue>
  borderBottomWidth?: ConditionalValue<BorderWidthsValue>
  borderCollapse?: ConditionalValue<BorderCollapseValue>
  borderColor?: ConditionalValue<ColorsValue>
  borderEnd?: ConditionalValue<BordersValue>
  borderEndColor?: ConditionalValue<ColorsValue>
  borderEndEndRadius?: ConditionalValue<RadiiValue>
  borderEndRadius?: ConditionalValue<RadiiValue>
  borderEndStartRadius?: ConditionalValue<RadiiValue>
  borderEndWidth?: ConditionalValue<BorderWidthsValue>
  borderInline?: ConditionalValue<BordersValue>
  borderInlineColor?: ConditionalValue<ColorsValue>
  borderInlineEnd?: ConditionalValue<BordersValue>
  borderInlineEndColor?: ConditionalValue<ColorsValue>
  borderInlineEndWidth?: ConditionalValue<BorderWidthsValue>
  borderInlineStart?: ConditionalValue<BordersValue>
  borderInlineStartColor?: ConditionalValue<ColorsValue>
  borderInlineStartWidth?: ConditionalValue<BorderWidthsValue>
  borderInlineWidth?: ConditionalValue<BorderWidthsValue>
  borderLeft?: ConditionalValue<BordersValue>
  borderLeftColor?: ConditionalValue<ColorsValue>
  borderLeftRadius?: ConditionalValue<RadiiValue>
  borderLeftWidth?: ConditionalValue<BorderWidthsValue>
  borderRadius?: ConditionalValue<RadiiValue>
  borderRight?: ConditionalValue<BordersValue>
  borderRightColor?: ConditionalValue<ColorsValue>
  borderRightRadius?: ConditionalValue<RadiiValue>
  borderRightWidth?: ConditionalValue<BorderWidthsValue>
  borderSpacing?: ConditionalValue<BorderSpacingValue>
  borderSpacingX?: ConditionalValue<SpacingValue>
  borderSpacingY?: ConditionalValue<SpacingValue>
  borderStart?: ConditionalValue<BordersValue>
  borderStartColor?: ConditionalValue<ColorsValue>
  borderStartEndRadius?: ConditionalValue<RadiiValue>
  borderStartRadius?: ConditionalValue<RadiiValue>
  borderStartStartRadius?: ConditionalValue<RadiiValue>
  borderStartWidth?: ConditionalValue<BorderWidthsValue>
  borderTop?: ConditionalValue<BordersValue>
  borderTopColor?: ConditionalValue<ColorsValue>
  borderTopLeftRadius?: ConditionalValue<RadiiValue>
  borderTopRadius?: ConditionalValue<RadiiValue>
  borderTopRightRadius?: ConditionalValue<RadiiValue>
  borderTopWidth?: ConditionalValue<BorderWidthsValue>
  borderWidth?: ConditionalValue<BorderWidthsValue>
  borderX?: ConditionalValue<BordersValue>
  borderXColor?: ConditionalValue<ColorsValue>
  borderXWidth?: ConditionalValue<BorderWidthsValue>
  borderY?: ConditionalValue<BordersValue>
  borderYColor?: ConditionalValue<ColorsValue>
  borderYWidth?: ConditionalValue<BorderWidthsValue>
  bottom?: ConditionalValue<SpacingValue>
  boxDecorationBreak?: ConditionalValue<BoxDecorationBreakValue>
  boxShadow?: ConditionalValue<ShadowsValue>
  boxShadowColor?: ConditionalValue<ColorsValue>
  boxSize?: ConditionalValue<BoxSizeValue>
  boxSizing?: ConditionalValue<BoxSizingValue>
  brightness?: ConditionalValue<BrightnessValue>
  caretColor?: ConditionalValue<ColorsValue>
  clipPath?: ConditionalValue<ClipPathValue>
  color?: ConditionalValue<ColorsValue>
  colorPalette?: ConditionalValue<ColorPaletteValue>
  columnGap?: ConditionalValue<SpacingValue>
  containerType?: ConditionalValue<ContainerTypeValue>
  contrast?: ConditionalValue<ContrastValue>
  cursor?: ConditionalValue<CursorValue>
  debug?: ConditionalValue<DebugValue>
  display?: ConditionalValue<DisplayValue>
  divideColor?: ConditionalValue<ColorsValue>
  divideStyle?: ConditionalValue<BorderStyleValue>
  divideX?: ConditionalValue<BorderWidthsValue>
  divideY?: ConditionalValue<BorderWidthsValue>
  dropShadow?: ConditionalValue<DropShadowsValue>
  end?: ConditionalValue<SpacingValue>
  fill?: ConditionalValue<ColorsValue>
  filter?: ConditionalValue<FilterValue>
  flex?: ConditionalValue<FlexValue>
  flexBasis?: ConditionalValue<FlexBasisValue>
  flexDir?: ConditionalValue<FlexDirectionValue>
  flexDirection?: ConditionalValue<FlexDirectionValue>
  flexGrow?: ConditionalValue<FlexGrowValue>
  flexShrink?: ConditionalValue<FlexShrinkValue>
  float?: ConditionalValue<FloatValue>
  focusRing?: ConditionalValue<FocusRingValue>
  focusRingColor?: ConditionalValue<ColorsValue>
  focusRingOffset?: ConditionalValue<SpacingValue>
  focusRingStyle?: ConditionalValue<BorderStylesValue>
  focusRingWidth?: ConditionalValue<BorderWidthsValue>
  focusVisibleRing?: ConditionalValue<FocusVisibleRingValue>
  fontFamily?: ConditionalValue<FontsValue>
  fontFeatureSettings?: ConditionalValue<FontFeatureSettingsValue>
  fontKerning?: ConditionalValue<FontKerningValue>
  fontPalette?: ConditionalValue<FontPaletteValue>
  fontSize?: ConditionalValue<FontSizesValue>
  fontSizeAdjust?: ConditionalValue<FontSizeAdjustValue>
  fontSmoothing?: ConditionalValue<FontSmoothingValue>
  fontVariant?: ConditionalValue<FontVariantValue>
  fontVariantAlternates?: ConditionalValue<FontVariantAlternatesValue>
  fontVariantCaps?: ConditionalValue<FontVariantCapsValue>
  fontVariantNumeric?: ConditionalValue<FontVariantNumericValue>
  fontVariationSettings?: ConditionalValue<FontVariationSettingsValue>
  fontWeight?: ConditionalValue<FontWeightsValue>
  gap?: ConditionalValue<SpacingValue>
  gradientFrom?: ConditionalValue<ColorsValue>
  gradientFromPosition?: ConditionalValue<GradientFromPositionValue>
  gradientTo?: ConditionalValue<ColorsValue>
  gradientToPosition?: ConditionalValue<GradientToPositionValue>
  gradientVia?: ConditionalValue<ColorsValue>
  gradientViaPosition?: ConditionalValue<GradientViaPositionValue>
  grayscale?: ConditionalValue<GrayscaleValue>
  gridAutoColumns?: ConditionalValue<GridAutoColumnsValue>
  gridAutoFlow?: ConditionalValue<GridAutoFlowValue>
  gridAutoRows?: ConditionalValue<GridAutoRowsValue>
  gridColumn?: ConditionalValue<GridColumnValue>
  gridColumnEnd?: ConditionalValue<GridColumnEndValue>
  gridColumnGap?: ConditionalValue<SpacingValue>
  gridColumnStart?: ConditionalValue<GridColumnStartValue>
  gridGap?: ConditionalValue<SpacingValue>
  gridRow?: ConditionalValue<GridRowValue>
  gridRowGap?: ConditionalValue<SpacingValue>
  gridTemplateColumns?: ConditionalValue<GridTemplateColumnsValue>
  gridTemplateRows?: ConditionalValue<GridTemplateRowsValue>
  h?: ConditionalValue<HeightValue>
  height?: ConditionalValue<HeightValue>
  hideBelow?: ConditionalValue<BreakpointsValue>
  hideFrom?: ConditionalValue<BreakpointsValue>
  hueRotate?: ConditionalValue<HueRotateValue>
  hyphens?: ConditionalValue<HyphensValue>
  inlineSize?: ConditionalValue<InlineSizeValue>
  inset?: ConditionalValue<InsetValue>
  insetBlock?: ConditionalValue<SpacingValue>
  insetBlockEnd?: ConditionalValue<SpacingValue>
  insetBlockStart?: ConditionalValue<SpacingValue>
  insetEnd?: ConditionalValue<SpacingValue>
  insetInline?: ConditionalValue<SpacingValue>
  insetInlineEnd?: ConditionalValue<SpacingValue>
  insetInlineStart?: ConditionalValue<SpacingValue>
  insetStart?: ConditionalValue<SpacingValue>
  insetX?: ConditionalValue<SpacingValue>
  insetY?: ConditionalValue<SpacingValue>
  invert?: ConditionalValue<InvertValue>
  justifyContent?: ConditionalValue<JustifyContentValue>
  left?: ConditionalValue<SpacingValue>
  letterSpacing?: ConditionalValue<LetterSpacingsValue>
  lineClamp?: ConditionalValue<LineClampValue>
  lineHeight?: ConditionalValue<LineHeightsValue>
  listStyle?: ConditionalValue<ListStyleValue>
  listStyleImage?: ConditionalValue<AssetsValue>
  listStylePosition?: ConditionalValue<ListStylePositionValue>
  listStyleType?: ConditionalValue<ListStyleTypeValue>
  m?: ConditionalValue<MarginValue>
  margin?: ConditionalValue<MarginValue>
  marginBlock?: ConditionalValue<MarginBlockValue>
  marginBlockEnd?: ConditionalValue<MarginBlockEndValue>
  marginBlockStart?: ConditionalValue<MarginBlockStartValue>
  marginBottom?: ConditionalValue<MarginBottomValue>
  marginEnd?: ConditionalValue<MarginInlineEndValue>
  marginInline?: ConditionalValue<MarginInlineValue>
  marginInlineEnd?: ConditionalValue<MarginInlineEndValue>
  marginInlineStart?: ConditionalValue<MarginInlineStartValue>
  marginLeft?: ConditionalValue<MarginLeftValue>
  marginRight?: ConditionalValue<MarginRightValue>
  marginStart?: ConditionalValue<MarginInlineStartValue>
  marginTop?: ConditionalValue<MarginTopValue>
  marginX?: ConditionalValue<MarginInlineValue>
  marginY?: ConditionalValue<MarginBlockValue>
  mask?: ConditionalValue<MaskValue>
  maskImage?: ConditionalValue<MaskImageValue>
  maskSize?: ConditionalValue<MaskSizeValue>
  maxBlockSize?: ConditionalValue<MaxBlockSizeValue>
  maxH?: ConditionalValue<MaxHeightValue>
  maxHeight?: ConditionalValue<MaxHeightValue>
  maxInlineSize?: ConditionalValue<MaxInlineSizeValue>
  maxW?: ConditionalValue<MaxWidthValue>
  maxWidth?: ConditionalValue<MaxWidthValue>
  mb?: ConditionalValue<MarginBottomValue>
  me?: ConditionalValue<MarginInlineEndValue>
  minBlockSize?: ConditionalValue<MinBlockSizeValue>
  minH?: ConditionalValue<MinHeightValue>
  minHeight?: ConditionalValue<MinHeightValue>
  minInlineSize?: ConditionalValue<MinInlineSizeValue>
  minW?: ConditionalValue<MinWidthValue>
  minWidth?: ConditionalValue<MinWidthValue>
  mixBlendMode?: ConditionalValue<MixBlendModeValue>
  ml?: ConditionalValue<MarginLeftValue>
  mr?: ConditionalValue<MarginRightValue>
  ms?: ConditionalValue<MarginInlineStartValue>
  mt?: ConditionalValue<MarginTopValue>
  mx?: ConditionalValue<MarginInlineValue>
  my?: ConditionalValue<MarginBlockValue>
  objectFit?: ConditionalValue<ObjectFitValue>
  objectPosition?: ConditionalValue<ObjectPositionValue>
  opacity?: ConditionalValue<OpacityValue>
  outline?: ConditionalValue<BordersValue>
  outlineColor?: ConditionalValue<ColorsValue>
  outlineOffset?: ConditionalValue<SpacingValue>
  outlineWidth?: ConditionalValue<BorderWidthsValue>
  overflow?: ConditionalValue<OverflowValue>
  overflowAnchor?: ConditionalValue<OverflowAnchorValue>
  overflowBlock?: ConditionalValue<OverflowBlockValue>
  overflowClipBox?: ConditionalValue<OverflowClipBoxValue>
  overflowClipMargin?: ConditionalValue<OverflowClipMarginValue>
  overflowInline?: ConditionalValue<OverflowInlineValue>
  overflowWrap?: ConditionalValue<OverflowWrapValue>
  overflowX?: ConditionalValue<OverflowXValue>
  overflowY?: ConditionalValue<OverflowYValue>
  overscrollBehavior?: ConditionalValue<OverscrollBehaviorValue>
  overscrollBehaviorBlock?: ConditionalValue<OverscrollBehaviorBlockValue>
  overscrollBehaviorInline?: ConditionalValue<OverscrollBehaviorInlineValue>
  overscrollBehaviorX?: ConditionalValue<OverscrollBehaviorXValue>
  overscrollBehaviorY?: ConditionalValue<OverscrollBehaviorYValue>
  p?: ConditionalValue<SpacingValue>
  padding?: ConditionalValue<SpacingValue>
  paddingBlock?: ConditionalValue<SpacingValue>
  paddingBlockEnd?: ConditionalValue<SpacingValue>
  paddingBlockStart?: ConditionalValue<SpacingValue>
  paddingBottom?: ConditionalValue<SpacingValue>
  paddingEnd?: ConditionalValue<SpacingValue>
  paddingInline?: ConditionalValue<SpacingValue>
  paddingInlineEnd?: ConditionalValue<SpacingValue>
  paddingInlineStart?: ConditionalValue<SpacingValue>
  paddingLeft?: ConditionalValue<SpacingValue>
  paddingRight?: ConditionalValue<SpacingValue>
  paddingStart?: ConditionalValue<SpacingValue>
  paddingTop?: ConditionalValue<SpacingValue>
  paddingX?: ConditionalValue<SpacingValue>
  paddingY?: ConditionalValue<SpacingValue>
  pb?: ConditionalValue<SpacingValue>
  pe?: ConditionalValue<SpacingValue>
  pl?: ConditionalValue<SpacingValue>
  pos?: ConditionalValue<PositionValue>
  position?: ConditionalValue<PositionValue>
  pr?: ConditionalValue<SpacingValue>
  ps?: ConditionalValue<SpacingValue>
  pt?: ConditionalValue<SpacingValue>
  px?: ConditionalValue<SpacingValue>
  py?: ConditionalValue<SpacingValue>
  right?: ConditionalValue<SpacingValue>
  ring?: ConditionalValue<BordersValue>
  ringColor?: ConditionalValue<ColorsValue>
  ringOffset?: ConditionalValue<SpacingValue>
  ringWidth?: ConditionalValue<BorderWidthsValue>
  rotate?: ConditionalValue<RotateValue>
  rotateX?: ConditionalValue<RotateValue>
  rotateY?: ConditionalValue<RotateValue>
  rotateZ?: ConditionalValue<RotateValue>
  rounded?: ConditionalValue<RadiiValue>
  roundedBottom?: ConditionalValue<RadiiValue>
  roundedBottomLeft?: ConditionalValue<RadiiValue>
  roundedBottomRight?: ConditionalValue<RadiiValue>
  roundedEnd?: ConditionalValue<RadiiValue>
  roundedEndEnd?: ConditionalValue<RadiiValue>
  roundedEndStart?: ConditionalValue<RadiiValue>
  roundedLeft?: ConditionalValue<RadiiValue>
  roundedRight?: ConditionalValue<RadiiValue>
  roundedStart?: ConditionalValue<RadiiValue>
  roundedStartEnd?: ConditionalValue<RadiiValue>
  roundedStartStart?: ConditionalValue<RadiiValue>
  roundedTop?: ConditionalValue<RadiiValue>
  roundedTopLeft?: ConditionalValue<RadiiValue>
  roundedTopRight?: ConditionalValue<RadiiValue>
  rowGap?: ConditionalValue<SpacingValue>
  saturate?: ConditionalValue<SaturateValue>
  scale?: ConditionalValue<ScaleValue>
  scaleX?: ConditionalValue<ScaleXValue>
  scaleY?: ConditionalValue<ScaleYValue>
  scrollBehavior?: ConditionalValue<ScrollBehaviorValue>
  scrollMargin?: ConditionalValue<SpacingValue>
  scrollMarginBlock?: ConditionalValue<SpacingValue>
  scrollMarginBlockEnd?: ConditionalValue<SpacingValue>
  scrollMarginBlockStart?: ConditionalValue<SpacingValue>
  scrollMarginBottom?: ConditionalValue<SpacingValue>
  scrollMarginInline?: ConditionalValue<SpacingValue>
  scrollMarginInlineEnd?: ConditionalValue<SpacingValue>
  scrollMarginInlineStart?: ConditionalValue<SpacingValue>
  scrollMarginLeft?: ConditionalValue<SpacingValue>
  scrollMarginRight?: ConditionalValue<SpacingValue>
  scrollMarginTop?: ConditionalValue<SpacingValue>
  scrollMarginX?: ConditionalValue<SpacingValue>
  scrollMarginY?: ConditionalValue<SpacingValue>
  scrollPadding?: ConditionalValue<SpacingValue>
  scrollPaddingBlock?: ConditionalValue<SpacingValue>
  scrollPaddingBlockEnd?: ConditionalValue<SpacingValue>
  scrollPaddingBlockStart?: ConditionalValue<SpacingValue>
  scrollPaddingBottom?: ConditionalValue<SpacingValue>
  scrollPaddingInline?: ConditionalValue<SpacingValue>
  scrollPaddingInlineEnd?: ConditionalValue<SpacingValue>
  scrollPaddingInlineStart?: ConditionalValue<SpacingValue>
  scrollPaddingLeft?: ConditionalValue<SpacingValue>
  scrollPaddingRight?: ConditionalValue<SpacingValue>
  scrollPaddingTop?: ConditionalValue<SpacingValue>
  scrollPaddingX?: ConditionalValue<SpacingValue>
  scrollPaddingY?: ConditionalValue<SpacingValue>
  scrollSnapAlign?: ConditionalValue<ScrollSnapAlignValue>
  scrollSnapCoordinate?: ConditionalValue<ScrollSnapCoordinateValue>
  scrollSnapDestination?: ConditionalValue<ScrollSnapDestinationValue>
  scrollSnapMargin?: ConditionalValue<SpacingValue>
  scrollSnapMarginBottom?: ConditionalValue<SpacingValue>
  scrollSnapMarginLeft?: ConditionalValue<SpacingValue>
  scrollSnapMarginRight?: ConditionalValue<SpacingValue>
  scrollSnapMarginTop?: ConditionalValue<SpacingValue>
  scrollSnapPointsX?: ConditionalValue<ScrollSnapPointsXValue>
  scrollSnapPointsY?: ConditionalValue<ScrollSnapPointsYValue>
  scrollSnapStop?: ConditionalValue<ScrollSnapStopValue>
  scrollSnapStrictness?: ConditionalValue<ScrollSnapStrictnessValue>
  scrollSnapType?: ConditionalValue<ScrollSnapTypeValue>
  scrollSnapTypeX?: ConditionalValue<ScrollSnapTypeXValue>
  scrollSnapTypeY?: ConditionalValue<ScrollSnapTypeYValue>
  scrollTimeline?: ConditionalValue<ScrollTimelineValue>
  scrollTimelineAxis?: ConditionalValue<ScrollTimelineAxisValue>
  scrollTimelineName?: ConditionalValue<ScrollTimelineNameValue>
  scrollbar?: ConditionalValue<ScrollbarValue>
  scrollbarColor?: ConditionalValue<ColorsValue>
  scrollbarGutter?: ConditionalValue<ScrollbarGutterValue>
  scrollbarWidth?: ConditionalValue<SizesValue>
  sepia?: ConditionalValue<SepiaValue>
  shadow?: ConditionalValue<ShadowsValue>
  shadowColor?: ConditionalValue<ColorsValue>
  spaceX?: ConditionalValue<MarginInlineStartValue>
  spaceY?: ConditionalValue<MarginBlockStartValue>
  srOnly?: ConditionalValue<SrOnlyValue>
  start?: ConditionalValue<SpacingValue>
  stroke?: ConditionalValue<ColorsValue>
  strokeDasharray?: ConditionalValue<StrokeDasharrayValue>
  strokeDashoffset?: ConditionalValue<StrokeDashoffsetValue>
  strokeLinecap?: ConditionalValue<StrokeLinecapValue>
  strokeLinejoin?: ConditionalValue<StrokeLinejoinValue>
  strokeMiterlimit?: ConditionalValue<StrokeMiterlimitValue>
  strokeOpacity?: ConditionalValue<StrokeOpacityValue>
  strokeWidth?: ConditionalValue<BorderWidthsValue>
  tableLayout?: ConditionalValue<TableLayoutValue>
  textAlign?: ConditionalValue<TextAlignValue>
  textDecoration?: ConditionalValue<TextDecorationValue>
  textDecorationColor?: ConditionalValue<ColorsValue>
  textDecorationStyle?: ConditionalValue<TextDecorationStyleValue>
  textDecorationThickness?: ConditionalValue<TextDecorationThicknessValue>
  textEmphasisColor?: ConditionalValue<ColorsValue>
  textGradient?: ConditionalValue<TextGradientValue>
  textIndent?: ConditionalValue<SpacingValue>
  textOverflow?: ConditionalValue<TextOverflowValue>
  textShadow?: ConditionalValue<ShadowsValue>
  textShadowColor?: ConditionalValue<ColorsValue>
  textSizeAdjust?: ConditionalValue<TextSizeAdjustValue>
  textStyle?: ConditionalValue<TextStyleValue>
  textTransform?: ConditionalValue<TextTransformValue>
  textUnderlineOffset?: ConditionalValue<TextUnderlineOffsetValue>
  textWrap?: ConditionalValue<TextWrapValue>
  top?: ConditionalValue<SpacingValue>
  touchAction?: ConditionalValue<TouchActionValue>
  transform?: ConditionalValue<TransformValue>
  transformBox?: ConditionalValue<TransformBoxValue>
  transformOrigin?: ConditionalValue<TransformOriginValue>
  transformStyle?: ConditionalValue<TransformStyleValue>
  transition?: ConditionalValue<TransitionValue>
  transitionDelay?: ConditionalValue<DurationsValue>
  transitionDuration?: ConditionalValue<DurationsValue>
  transitionProperty?: ConditionalValue<TransitionPropertyValue>
  transitionTimingFunction?: ConditionalValue<EasingsValue>
  translate?: ConditionalValue<TranslateValue>
  translateX?: ConditionalValue<TranslateXValue>
  translateY?: ConditionalValue<TranslateYValue>
  translateZ?: ConditionalValue<TranslateZValue>
  truncate?: ConditionalValue<TruncateValue>
  userSelect?: ConditionalValue<UserSelectValue>
  verticalAlign?: ConditionalValue<VerticalAlignValue>
  visibility?: ConditionalValue<VisibilityValue>
  w?: ConditionalValue<WidthValue>
  width?: ConditionalValue<WidthValue>
  wordBreak?: ConditionalValue<WordBreakValue>
  x?: ConditionalValue<TranslateXValue>
  y?: ConditionalValue<TranslateYValue>
  z?: ConditionalValue<TranslateZValue>
  zIndex?: ConditionalValue<ZIndexValue>
  container?: ConditionalValue<ContainerValue>
  containerName?: ConditionalValue<ContainerName>
  all?: ConditionalValue<WithEscapeHatch<Globals | CssVars>>
  borderBlockEndStyle?: ConditionalValue<WithEscapeHatch<Globals | CssVars | OnlyKnown<"dashed" | "dotted" | "double" | "groove" | "hidden" | "inset" | "none" | "outset" | "ridge" | "solid">>>
  borderBlockStartStyle?: ConditionalValue<WithEscapeHatch<Globals | CssVars | OnlyKnown<"dashed" | "dotted" | "double" | "groove" | "hidden" | "inset" | "none" | "outset" | "ridge" | "solid">>>
  borderBlockStyle?: ConditionalValue<WithEscapeHatch<Globals | CssVars | OnlyKnown<"dashed" | "dotted" | "double" | "groove" | "hidden" | "inset" | "none" | "outset" | "ridge" | "solid">>>
  borderBottomStyle?: ConditionalValue<WithEscapeHatch<Globals | CssVars | OnlyKnown<"dashed" | "dotted" | "double" | "groove" | "hidden" | "inset" | "none" | "outset" | "ridge" | "solid">>>
  borderInlineEndStyle?: ConditionalValue<WithEscapeHatch<Globals | CssVars | OnlyKnown<"dashed" | "dotted" | "double" | "groove" | "hidden" | "inset" | "none" | "outset" | "ridge" | "solid">>>
  borderInlineStartStyle?: ConditionalValue<WithEscapeHatch<Globals | CssVars | OnlyKnown<"dashed" | "dotted" | "double" | "groove" | "hidden" | "inset" | "none" | "outset" | "ridge" | "solid">>>
  borderInlineStyle?: ConditionalValue<WithEscapeHatch<Globals | CssVars | OnlyKnown<"dashed" | "dotted" | "double" | "groove" | "hidden" | "inset" | "none" | "outset" | "ridge" | "solid">>>
  borderLeftStyle?: ConditionalValue<WithEscapeHatch<Globals | CssVars | OnlyKnown<"dashed" | "dotted" | "double" | "groove" | "hidden" | "inset" | "none" | "outset" | "ridge" | "solid">>>
  borderRightStyle?: ConditionalValue<WithEscapeHatch<Globals | CssVars | OnlyKnown<"dashed" | "dotted" | "double" | "groove" | "hidden" | "inset" | "none" | "outset" | "ridge" | "solid">>>
  borderTopStyle?: ConditionalValue<WithEscapeHatch<Globals | CssVars | OnlyKnown<"dashed" | "dotted" | "double" | "groove" | "hidden" | "inset" | "none" | "outset" | "ridge" | "solid">>>
  breakAfter?: ConditionalValue<WithEscapeHatch<Globals | CssVars | OnlyKnown<"all" | "always" | "auto" | "avoid" | "avoid-column" | "avoid-page" | "avoid-region" | "column" | "left" | "page" | "recto" | "region" | "right" | "verso">>>
  breakBefore?: ConditionalValue<WithEscapeHatch<Globals | CssVars | OnlyKnown<"all" | "always" | "auto" | "avoid" | "avoid-column" | "avoid-page" | "avoid-region" | "column" | "left" | "page" | "recto" | "region" | "right" | "verso">>>
  breakInside?: ConditionalValue<WithEscapeHatch<Globals | CssVars | OnlyKnown<"auto" | "avoid" | "avoid-column" | "avoid-page" | "avoid-region">>>
  captionSide?: ConditionalValue<WithEscapeHatch<Globals | CssVars | OnlyKnown<"bottom" | "top">>>
  clear?: ConditionalValue<WithEscapeHatch<Globals | CssVars | OnlyKnown<"both" | "inline-end" | "inline-start" | "left" | "none" | "right">>>
  columnFill?: ConditionalValue<WithEscapeHatch<Globals | CssVars | OnlyKnown<"auto" | "balance">>>
  columnRuleStyle?: ConditionalValue<WithEscapeHatch<Globals | CssVars | OnlyKnown<"dashed" | "dotted" | "double" | "groove" | "hidden" | "inset" | "none" | "outset" | "ridge" | "solid">>>
  contentVisibility?: ConditionalValue<WithEscapeHatch<Globals | CssVars | OnlyKnown<"auto" | "hidden" | "visible">>>
  direction?: ConditionalValue<WithEscapeHatch<Globals | CssVars | OnlyKnown<"ltr" | "rtl">>>
  emptyCells?: ConditionalValue<WithEscapeHatch<Globals | CssVars | OnlyKnown<"hide" | "show">>>
  flexWrap?: ConditionalValue<WithEscapeHatch<Globals | CssVars | OnlyKnown<"nowrap" | "wrap" | "wrap-reverse">>>
  forcedColorAdjust?: ConditionalValue<WithEscapeHatch<Globals | CssVars | OnlyKnown<"auto" | "none" | "preserve-parent-color">>>
  isolation?: ConditionalValue<WithEscapeHatch<Globals | CssVars | OnlyKnown<"auto" | "isolate">>>
  lineBreak?: ConditionalValue<WithEscapeHatch<Globals | CssVars | OnlyKnown<"anywhere" | "auto" | "loose" | "normal" | "strict">>>
  pointerEvents?: ConditionalValue<WithEscapeHatch<Globals | CssVars | OnlyKnown<"all" | "auto" | "fill" | "none" | "painted" | "stroke" | "visible" | "visibleFill" | "visiblePainted" | "visibleStroke">>>
  resize?: ConditionalValue<WithEscapeHatch<Globals | CssVars | OnlyKnown<"block" | "both" | "horizontal" | "inline" | "none" | "vertical">>>
  writingMode?: ConditionalValue<WithEscapeHatch<Globals | CssVars | OnlyKnown<"horizontal-tb" | "sideways-lr" | "sideways-rl" | "vertical-lr" | "vertical-rl">>>
}

export type CssVarValue = ConditionalValue<CssVars | AnyString | AnyNumber>

export type CssVarProperties = {
  [K in `--${string}`]?: CssVarValue
}

export type NestedStyles = {
  [K in Selector | Condition]?: SystemStyleObject
}

export interface SystemStyleObject extends SystemProperties, CssVarProperties, NestedStyles {}

export interface GlobalStyleObject {
  [selector: string]: SystemStyleObject
}

export interface CssKeyframes {
  [name: string]: {
    [time: string]: CssProperties
  }
}

export interface GlobalFontfaceRule {
  fontFamily: string
  src: string
  fontStyle?: string
  fontWeight?: string | number
  fontDisplay?: "auto" | "block" | "swap" | "fallback" | "optional"
  unicodeRange?: string
  fontFeatureSettings?: string
  fontVariationSettings?: string
  fontStretch?: string
  ascentOverride?: string
  descentOverride?: string
  lineGapOverride?: string
  sizeAdjust?: string
}

export type FontfaceRule = Omit<GlobalFontfaceRule, "fontFamily">

export interface GlobalFontface {
  [name: string]: FontfaceRule | FontfaceRule[]
}

interface WithCss {
  css?: SystemStyleObject | SystemStyleObject[]
}

export type JsxStyleProps = SystemStyleObject & WithCss

export interface PatchedHTMLProps {
  htmlWidth?: string | number
  htmlHeight?: string | number
  htmlTranslate?: "yes" | "no" | undefined
  htmlContent?: string
}

export type OmittedHTMLProps = "color" | "translate" | "transition" | "width" | "height" | "content"

type WithHTMLProps<T> = DistributiveOmit<T, OmittedHTMLProps> & PatchedHTMLProps

export type JsxHTMLProps<T extends Record<string, any>, P extends Record<string, any> = {}> = Assign<WithHTMLProps<T>, P>