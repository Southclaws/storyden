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
  "@/2xl": string
  "@/2xlDown": string
  "@/2xlOnly": string
  "@/2xlTo3xl": string
  "@/2xlTo4xl": string
  "@/2xlTo5xl": string
  "@/2xlTo6xl": string
  "@/2xlTo7xl": string
  "@/2xlTo8xl": string
  "@/2xs": string
  "@/2xsDown": string
  "@/2xsOnly": string
  "@/2xsTo2xl": string
  "@/2xsTo3xl": string
  "@/2xsTo4xl": string
  "@/2xsTo5xl": string
  "@/2xsTo6xl": string
  "@/2xsTo7xl": string
  "@/2xsTo8xl": string
  "@/2xsToLg": string
  "@/2xsToMd": string
  "@/2xsToSm": string
  "@/2xsToXl": string
  "@/2xsToXs": string
  "@/3xl": string
  "@/3xlDown": string
  "@/3xlOnly": string
  "@/3xlTo4xl": string
  "@/3xlTo5xl": string
  "@/3xlTo6xl": string
  "@/3xlTo7xl": string
  "@/3xlTo8xl": string
  "@/3xs": string
  "@/3xsDown": string
  "@/3xsOnly": string
  "@/3xsTo2xl": string
  "@/3xsTo2xs": string
  "@/3xsTo3xl": string
  "@/3xsTo4xl": string
  "@/3xsTo5xl": string
  "@/3xsTo6xl": string
  "@/3xsTo7xl": string
  "@/3xsTo8xl": string
  "@/3xsToLg": string
  "@/3xsToMd": string
  "@/3xsToSm": string
  "@/3xsToXl": string
  "@/3xsToXs": string
  "@/4xl": string
  "@/4xlDown": string
  "@/4xlOnly": string
  "@/4xlTo5xl": string
  "@/4xlTo6xl": string
  "@/4xlTo7xl": string
  "@/4xlTo8xl": string
  "@/5xl": string
  "@/5xlDown": string
  "@/5xlOnly": string
  "@/5xlTo6xl": string
  "@/5xlTo7xl": string
  "@/5xlTo8xl": string
  "@/6xl": string
  "@/6xlDown": string
  "@/6xlOnly": string
  "@/6xlTo7xl": string
  "@/6xlTo8xl": string
  "@/7xl": string
  "@/7xlDown": string
  "@/7xlOnly": string
  "@/7xlTo8xl": string
  "@/8xl": string
  "@/8xlDown": string
  "@/8xlOnly": string
  "@/lg": string
  "@/lgDown": string
  "@/lgOnly": string
  "@/lgTo2xl": string
  "@/lgTo3xl": string
  "@/lgTo4xl": string
  "@/lgTo5xl": string
  "@/lgTo6xl": string
  "@/lgTo7xl": string
  "@/lgTo8xl": string
  "@/lgToXl": string
  "@/md": string
  "@/mdDown": string
  "@/mdOnly": string
  "@/mdTo2xl": string
  "@/mdTo3xl": string
  "@/mdTo4xl": string
  "@/mdTo5xl": string
  "@/mdTo6xl": string
  "@/mdTo7xl": string
  "@/mdTo8xl": string
  "@/mdToLg": string
  "@/mdToXl": string
  "@/sm": string
  "@/smDown": string
  "@/smOnly": string
  "@/smTo2xl": string
  "@/smTo3xl": string
  "@/smTo4xl": string
  "@/smTo5xl": string
  "@/smTo6xl": string
  "@/smTo7xl": string
  "@/smTo8xl": string
  "@/smToLg": string
  "@/smToMd": string
  "@/smToXl": string
  "@/xl": string
  "@/xlDown": string
  "@/xlOnly": string
  "@/xlTo2xl": string
  "@/xlTo3xl": string
  "@/xlTo4xl": string
  "@/xlTo5xl": string
  "@/xlTo6xl": string
  "@/xlTo7xl": string
  "@/xlTo8xl": string
  "@/xs": string
  "@/xsDown": string
  "@/xsOnly": string
  "@/xsTo2xl": string
  "@/xsTo3xl": string
  "@/xsTo4xl": string
  "@/xsTo5xl": string
  "@/xsTo6xl": string
  "@/xsTo7xl": string
  "@/xsTo8xl": string
  "@/xsToLg": string
  "@/xsToMd": string
  "@/xsToSm": string
  "@/xsToXl": string
  "_active": string
  "_after": string
  "_anyPointerCoarse": string
  "_anyPointerFine": string
  "_anyPointerNone": string
  "_atValue": string
  "_autofill": string
  "_backdrop": string
  "_before": string
  "_checked": string
  "_closed": string
  "_complete": string
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
  "_inert": string
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
  "_pointerCoarse": string
  "_pointerFine": string
  "_pointerNone": string
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
  "_userInvalid": string
  "_userValid": string
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
export type CssContainer = ContainerName | `${ContainerName} / inline-size` | `${ContainerName} / size` | AnyString

export type Condition = keyof Conditions

export type ConditionalValue<T> =
  | T
  | Array<T | null | undefined>
  | { [K in Condition]?: ConditionalValue<T> }
  | undefined

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

export type AlignContentValue = WithEscapeHatch<CssGlobals | CssVars | OnlyKnown<"baseline" | "center" | "end" | "first baseline" | "flex-end" | "flex-start" | "last baseline" | "normal" | "safe center" | "safe end" | "safe start" | "space-around" | "space-between" | "space-evenly" | "start" | "stretch">>

export type AlignItemsValue = WithEscapeHatch<CssGlobals | CssVars | OnlyKnown<"anchor-center" | "baseline" | "center" | "end" | "first baseline" | "flex-end" | "flex-start" | "last baseline" | "normal" | "safe center" | "safe end" | "safe start" | "self-end" | "self-start" | "start" | "stretch">>

export type AlignSelfValue = WithEscapeHatch<CssGlobals | CssVars | OnlyKnown<"anchor-center" | "auto" | "baseline" | "center" | "end" | "first baseline" | "flex-end" | "flex-start" | "last baseline" | "normal" | "safe center" | "safe end" | "safe start" | "self-end" | "self-start" | "start" | "stretch">>

export type AnimationCompositionValue = string | number | CssVars | AnyString

export type AnimationDirectionValue = string | number | CssVars | AnyString

export type AnimationFillModeValue = string | number | CssVars | AnyString

export type AnimationIterationCountValue = string | number | CssVars | AnyString

export type AnimationPlayStateValue = string | number | CssVars | AnyString

export type AnimationRangeEndValue = string | number | CssVars | AnyString

export type AnimationRangeStartValue = string | number | CssVars | AnyString

export type AnimationRangeValue = string | number | CssVars | AnyString

export type AnimationStateValue = string | number | CssVars | AnyString

export type AnimationTimelineValue = string | number | CssVars | AnyString

export type AnimationsValue = CssGlobals | TokenValue<"animations"> | CssVars | AnyString | AnyNumber

export type AppearanceValue = WithEscapeHatch<PropertyValueMap["appearance"]>

export type AspectRatiosValue = CssAutoGlobals | TokenValue<"aspectRatios"> | CssVars | AnyString | AnyNumber

export type AssetsValue = PropertyValueMap["listStyleImage"]

export type BackdropBrightnessValue = string | number | CssVars | AnyString

export type BackdropContrastValue = string | number | CssVars | AnyString

export type BackdropFilterValue = "auto" | CssVars | AnyString | AnyNumber

export type BackdropGrayscaleValue = string | number | CssVars | AnyString

export type BackdropHueRotateValue = string | number | CssVars | AnyString

export type BackdropInvertValue = string | number | CssVars | AnyString

export type BackdropOpacityValue = string | number | CssVars | AnyString

export type BackdropSaturateValue = string | number | CssVars | AnyString

export type BackdropSepiaValue = string | number | CssVars | AnyString

export type BackfaceVisibilityValue = string | number | CssVars | AnyString

export type BackgroundAttachmentValue = string | number | CssVars | AnyString

export type BackgroundBlendModeValue = string | number | CssVars | AnyString

export type BackgroundClipValue = string | number | CssVars | AnyString

export type BackgroundConicValue = string | number | CssVars | AnyString

export type BackgroundGradientValue = "to-b" | "to-bl" | "to-br" | "to-l" | "to-r" | "to-t" | "to-tl" | "to-tr" | CssVars | AnyString | AnyNumber

export type BackgroundLinearValue = "to-b" | "to-bl" | "to-br" | "to-l" | "to-r" | "to-t" | "to-tl" | "to-tr" | CssVars | AnyString | AnyNumber

export type BackgroundOriginValue = string | number | CssVars | AnyString

export type BackgroundPositionValue = PropertyValueMap["backgroundPosition"]

export type BackgroundPositionXValue = string | number | CssVars | AnyString

export type BackgroundPositionYValue = string | number | CssVars | AnyString

export type BackgroundRepeatValue = PropertyValueMap["backgroundRepeat"]

export type BackgroundSizeValue = PropertyValueMap["backgroundSize"]

export type BlockSizeValue = CssDimensionGlobals | TokenValue<"sizes"> | "1/2" | "1/3" | "1/4" | "1/5" | "1/6" | "2/3" | "2/4" | "2/5" | "2/6" | "3/4" | "3/5" | "3/6" | "4/5" | "4/6" | "5/6" | "auto" | "dvh" | "lvh" | "screen" | "svh" | CssVars | AnyString | AnyNumber

export type BlursValue = CssGlobals | TokenValue<"blurs"> | CssVars | AnyString | AnyNumber

export type BorderCollapseValue = string | number | CssVars | AnyString

export type BorderStyleValue = PropertyValueMap["borderStyle"] | string | number | CssVars | AnyString

export type BorderStylesValue = CssGlobals | TokenValue<"borderStyles"> | PropertyValueMap["outlineStyle"] | CssVars | AnyString | AnyNumber

export type BorderWidthsValue = PropertyValueMap["strokeWidth"]

export type BordersValue = CssGlobals | TokenValue<"borders"> | CssVars | AnyString | AnyNumber

export type BoxDecorationBreakValue = string | number | CssVars | AnyString

export type BoxSizeValue = CssDimensionGlobals | TokenValue<"sizes"> | "1/12" | "1/2" | "1/3" | "1/4" | "1/5" | "1/6" | "10/12" | "11/12" | "2/12" | "2/3" | "2/4" | "2/5" | "2/6" | "3/12" | "3/4" | "3/5" | "3/6" | "4/12" | "4/5" | "4/6" | "5/12" | "5/6" | "6/12" | "7/12" | "8/12" | "9/12" | "auto" | "screen" | CssVars | AnyString | AnyNumber

export type BoxSizingValue = string | number | CssVars | AnyString

export type BreakpointsValue = CssGlobals | TokenValue<"breakpoints"> | CssVars | AnyString | AnyNumber

export type BrightnessValue = string | number | CssVars | AnyString

export type ClipPathValue = PropertyValueMap["clipPath"]

export type ColorPaletteValue = "Mono" | "Primary" | "Shades" | "amber" | "black" | "blue" | "current" | "cyan" | "emerald" | "fuchsia" | "gray" | "green" | "indigo" | "lime" | "mauve" | "mist" | "neutral" | "olive" | "orange" | "pink" | "purple" | "red" | "rose" | "sky" | "slate" | "stone" | "taupe" | "teal" | "transparent" | "violet" | "white" | "yellow" | "zinc" | CssVars | AnyString | AnyNumber

export type ColorsValue = CssColorGlobals | TokenValue<"colors"> | CssVars | AnyString | AnyNumber

export type ContainerNamesValue = CssGlobals | TokenValue<"containerNames"> | CssAny | CssVars | AnyString | AnyNumber

export type ContainerTypeValue = PropertyValueMap["containerType"]

export type ContainerValue = string | number | CssVars | AnyString

export type ContrastValue = string | number | CssVars | AnyString

export type CursorValue = PropertyValueMap["cursor"]

export type DebugValue = boolean | CssVars | AnyString | AnyNumber

export type DisplayValue = WithEscapeHatch<CssGlobals | CssVars | OnlyKnown<"-ms-flexbox" | "-ms-grid" | "-ms-inline-flexbox" | "-ms-inline-grid" | "-webkit-flex" | "-webkit-inline-flex" | "block" | "contents" | "flex" | "flow" | "flow-root" | "grid" | "inline" | "inline-block" | "inline-flex" | "inline-grid" | "inline-list-item" | "inline-table" | "list-item" | "none" | "ruby" | "ruby-base" | "ruby-base-container" | "ruby-text" | "ruby-text-container" | "run-in" | "table" | "table-caption" | "table-cell" | "table-column" | "table-column-group" | "table-footer-group" | "table-header-group" | "table-row" | "table-row-group">>

export type DropShadowsValue = CssGlobals | TokenValue<"dropShadows"> | CssVars | AnyString | AnyNumber

export type DurationsValue = CssGlobals | TokenValue<"durations"> | CssVars | AnyString | AnyNumber

export type EasingsValue = CssGlobals | TokenValue<"easings"> | CssVars | AnyString | AnyNumber

export type FilterValue = "auto" | CssVars | AnyString | AnyNumber

export type FlexBasisValue = "0" | "0.5" | "1" | "1.5" | "1/12" | "1/2" | "1/3" | "1/4" | "1/5" | "1/6" | "10" | "10/12" | "11" | "11/12" | "12" | "14" | "16" | "2" | "2.5" | "2/12" | "2/3" | "2/4" | "2/5" | "2/6" | "20" | "24" | "28" | "2xl" | "3" | "3.5" | "3/12" | "3/4" | "3/5" | "3/6" | "32" | "36" | "3xl" | "4" | "4.5" | "4/12" | "4/5" | "4/6" | "40" | "44" | "48" | "4xl" | "5" | "5.5" | "5/12" | "5/6" | "52" | "56" | "5xl" | "6" | "6/12" | "60" | "64" | "6xl" | "7" | "7/12" | "72" | "7xl" | "8" | "8/12" | "80" | "8xl" | "9" | "9/12" | "96" | "breakpoint-2xl" | "breakpoint-lg" | "breakpoint-md" | "breakpoint-sm" | "breakpoint-xl" | "fit" | "full" | "lg" | "max" | "md" | "min" | "prose" | "sm" | "xl" | "xs" | CssVars | AnyString | AnyNumber

export type FlexDirectionValue = string | number | CssVars | AnyString

export type FlexGrowValue = PropertyValueMap["flexGrow"]

export type FlexShrinkValue = PropertyValueMap["flexShrink"]

export type FlexValue = "1" | "auto" | "initial" | "none" | CssVars | AnyString | AnyNumber

export type FloatValue = WithEscapeHatch<CssGlobals | CssVars | OnlyKnown<"end" | "start" | PropertyValueMap["float"]>>

export type FocusRingValue = "inside" | "mixed" | "none" | "outside" | CssVars | AnyString | AnyNumber

export type FocusVisibleRingValue = "inside" | "mixed" | "none" | "outside" | CssVars | AnyString | AnyNumber

export type FontFeatureSettingsValue = string | number | CssVars | AnyString

export type FontKerningValue = string | number | CssVars | AnyString

export type FontPaletteValue = string | number | CssVars | AnyString

export type FontSizeAdjustValue = PropertyValueMap["fontSizeAdjust"]

export type FontSizesValue = CssGlobals | TokenValue<"fontSizes"> | CssVars | AnyString | AnyNumber

export type FontSmoothingValue = "antialiased" | "subpixel-antialiased" | CssVars | AnyString | AnyNumber

export type FontVariantAlternatesValue = string | number | CssVars | AnyString

export type FontVariantCapsValue = PropertyValueMap["fontVariantCaps"] | AnyString

export type FontVariantNumericValue = PropertyValueMap["fontVariantNumeric"]

export type FontVariantValue = string | number | CssVars | AnyString

export type FontVariationSettingsValue = string | number | CssVars | AnyString

export type FontWeightsValue = CssGlobals | TokenValue<"fontWeights"> | CssVars | AnyString | AnyNumber

export type FontsValue = CssGlobals | TokenValue<"fonts"> | CssVars | AnyString | AnyNumber

export type GradientFromPositionValue = string | number | CssVars | AnyString

export type GradientToPositionValue = string | number | CssVars | AnyString

export type GradientViaPositionValue = string | number | CssVars | AnyString

export type GradientsValue = string | number | CssVars | AnyString

export type GrayscaleValue = string | number | CssVars | AnyString

export type GridAutoColumnsValue = "fr" | "max" | "min" | CssVars | AnyString | AnyNumber

export type GridAutoFlowValue = PropertyValueMap["gridAutoFlow"]

export type GridAutoRowsValue = "fr" | "max" | "min" | CssVars | AnyString | AnyNumber

export type GridColumnEndValue = string | number | CssVars | AnyString

export type GridColumnStartValue = string | number | CssVars | AnyString

export type GridColumnValue = string | number | CssVars | AnyString

export type GridRowValue = string | number | CssVars | AnyString

export type GridTemplateColumnsValue = string | number | CssVars | AnyString

export type GridTemplateRowsValue = string | number | CssVars | AnyString

export type HeightValue = CssDimensionGlobals | TokenValue<"sizes"> | "1/2" | "1/3" | "1/4" | "1/5" | "1/6" | "2/3" | "2/4" | "2/5" | "2/6" | "3/4" | "3/5" | "3/6" | "4/5" | "4/6" | "5/6" | "auto" | "dvh" | "lvh" | "screen" | "svh" | CssVars | AnyString | AnyNumber

export type HueRotateValue = string | number | CssVars | AnyString

export type HyphensValue = PropertyValueMap["hyphens"] | AnyString

export type InlineSizeValue = CssDimensionGlobals | TokenValue<"sizes"> | "1/12" | "1/2" | "1/3" | "1/4" | "1/5" | "1/6" | "10/12" | "11/12" | "2/12" | "2/3" | "2/4" | "2/5" | "2/6" | "3/12" | "3/4" | "3/5" | "3/6" | "4/12" | "4/5" | "4/6" | "5/12" | "5/6" | "6/12" | "7/12" | "8/12" | "9/12" | "auto" | "screen" | CssVars | AnyString | AnyNumber

export type InvertValue = string | number | CssVars | AnyString

export type JustifyContentValue = PropertyValueMap["justifyContent"]

export type KeyframesValue = "bounce" | "ping" | "pulse" | "spin" | CssVars | AnyString | AnyNumber

export type LetterSpacingsValue = CssGlobals | TokenValue<"letterSpacings"> | CssVars | AnyString | AnyNumber

export type LineClampValue = string | number | CssVars | AnyString

export type LineHeightsValue = CssGlobals | TokenValue<"lineHeights"> | CssVars | AnyString | AnyNumber

export type ListStylePositionValue = PropertyValueMap["listStylePosition"] | AnyString

export type ListStyleTypeValue = PropertyValueMap["listStyleType"]

export type ListStyleValue = string | number | CssVars | AnyString

export type MaskClipValue = string | number | CssVars | AnyString

export type MaskCompositeValue = string | number | CssVars | AnyString

export type MaskConicValue = string | number | CssVars | AnyString

export type MaskImageValue = string | number | CssVars | AnyString

export type MaskLinearValue = "to-b" | "to-bl" | "to-br" | "to-l" | "to-r" | "to-t" | "to-tl" | "to-tr" | CssVars | AnyString | AnyNumber

export type MaskModeValue = string | number | CssVars | AnyString

export type MaskOriginValue = string | number | CssVars | AnyString

export type MaskPositionValue = PropertyValueMap["maskPosition"]

export type MaskRadialAtValue = "bottom" | "bottom left" | "bottom right" | "center" | "left" | "right" | "top" | "top left" | "top right" | CssVars | AnyString | AnyNumber

export type MaskRadialShapeValue = "circle" | "ellipse" | CssVars | AnyString | AnyNumber

export type MaskRadialSizeValue = "closest-corner" | "closest-side" | "farthest-corner" | "farthest-side" | CssVars | AnyString | AnyNumber

export type MaskRadialValue = string | number | CssVars | AnyString

export type MaskRepeatValue = PropertyValueMap["maskRepeat"]

export type MaskSizeValue = PropertyValueMap["maskSize"]

export type MaskTypeValue = string | number | CssVars | AnyString

export type MaskValue = string | number | CssVars | AnyString

export type MaxBlockSizeValue = CssDimensionGlobals | TokenValue<"sizes"> | "1/2" | "1/3" | "1/4" | "1/5" | "1/6" | "2/3" | "2/4" | "2/5" | "2/6" | "3/4" | "3/5" | "3/6" | "4/5" | "4/6" | "5/6" | "auto" | "dvh" | "lvh" | "screen" | "svh" | CssVars | AnyString | AnyNumber

export type MaxHeightValue = CssDimensionGlobals | TokenValue<"sizes"> | "1/2" | "1/3" | "1/4" | "1/5" | "1/6" | "2/3" | "2/4" | "2/5" | "2/6" | "3/4" | "3/5" | "3/6" | "4/5" | "4/6" | "5/6" | "auto" | "dvh" | "lvh" | "screen" | "svh" | CssVars | AnyString | AnyNumber

export type MaxInlineSizeValue = CssDimensionGlobals | TokenValue<"sizes"> | "1/12" | "1/2" | "1/3" | "1/4" | "1/5" | "1/6" | "10/12" | "11/12" | "2/12" | "2/3" | "2/4" | "2/5" | "2/6" | "3/12" | "3/4" | "3/5" | "3/6" | "4/12" | "4/5" | "4/6" | "5/12" | "5/6" | "6/12" | "7/12" | "8/12" | "9/12" | "auto" | "screen" | CssVars | AnyString | AnyNumber

export type MaxWidthValue = CssDimensionGlobals | TokenValue<"sizes"> | "1/12" | "1/2" | "1/3" | "1/4" | "1/5" | "1/6" | "10/12" | "11/12" | "2/12" | "2/3" | "2/4" | "2/5" | "2/6" | "3/12" | "3/4" | "3/5" | "3/6" | "4/12" | "4/5" | "4/6" | "5/12" | "5/6" | "6/12" | "7/12" | "8/12" | "9/12" | "auto" | "screen" | CssVars | AnyString | AnyNumber

export type MinBlockSizeValue = CssDimensionGlobals | TokenValue<"sizes"> | "1/2" | "1/3" | "1/4" | "1/5" | "1/6" | "2/3" | "2/4" | "2/5" | "2/6" | "3/4" | "3/5" | "3/6" | "4/5" | "4/6" | "5/6" | "auto" | "dvh" | "lvh" | "screen" | "svh" | CssVars | AnyString | AnyNumber

export type MinHeightValue = CssDimensionGlobals | TokenValue<"sizes"> | "1/2" | "1/3" | "1/4" | "1/5" | "1/6" | "2/3" | "2/4" | "2/5" | "2/6" | "3/4" | "3/5" | "3/6" | "4/5" | "4/6" | "5/6" | "auto" | "dvh" | "lvh" | "screen" | "svh" | CssVars | AnyString | AnyNumber

export type MinInlineSizeValue = CssDimensionGlobals | TokenValue<"sizes"> | "1/12" | "1/2" | "1/3" | "1/4" | "1/5" | "1/6" | "10/12" | "11/12" | "2/12" | "2/3" | "2/4" | "2/5" | "2/6" | "3/12" | "3/4" | "3/5" | "3/6" | "4/12" | "4/5" | "4/6" | "5/12" | "5/6" | "6/12" | "7/12" | "8/12" | "9/12" | "auto" | "screen" | CssVars | AnyString | AnyNumber

export type MinWidthValue = CssDimensionGlobals | TokenValue<"sizes"> | "1/12" | "1/2" | "1/3" | "1/4" | "1/5" | "1/6" | "10/12" | "11/12" | "2/12" | "2/3" | "2/4" | "2/5" | "2/6" | "3/12" | "3/4" | "3/5" | "3/6" | "4/12" | "4/5" | "4/6" | "5/12" | "5/6" | "6/12" | "7/12" | "8/12" | "9/12" | "auto" | "screen" | CssVars | AnyString | AnyNumber

export type MixBlendModeValue = WithEscapeHatch<PropertyValueMap["mixBlendMode"]>

export type ObjectFitValue = WithEscapeHatch<PropertyValueMap["objectFit"]>

export type ObjectPositionValue = PropertyValueMap["objectPosition"]

export type OpacityValue = PropertyValueMap["opacity"]

export type OverflowAnchorValue = string | number | CssVars | AnyString

export type OverflowBlockValue = PropertyValueMap["overflowBlock"] | AnyString

export type OverflowClipBoxValue = string | number | CssVars | AnyString

export type OverflowClipMarginValue = PropertyValueMap["overflowClipMargin"]

export type OverflowInlineValue = PropertyValueMap["overflowInline"] | AnyString

export type OverflowValue = WithEscapeHatch<PropertyValueMap["overflow"]>

export type OverflowWrapValue = string | number | CssVars | AnyString

export type OverflowXValue = WithEscapeHatch<PropertyValueMap["overflowX"]>

export type OverflowYValue = PropertyValueMap["overflowY"] | AnyString

export type OverscrollBehaviorBlockValue = string | number | CssVars | AnyString

export type OverscrollBehaviorInlineValue = string | number | CssVars | AnyString

export type OverscrollBehaviorValue = PropertyValueMap["overscrollBehavior"]

export type OverscrollBehaviorXValue = PropertyValueMap["overscrollBehaviorX"] | AnyString

export type OverscrollBehaviorYValue = PropertyValueMap["overscrollBehaviorY"] | AnyString

export type PositionValue = WithEscapeHatch<PropertyValueMap["position"]>

export type RadiiValue = CssGlobals | TokenValue<"radii"> | CssVars | AnyString | AnyNumber

export type RotateValue = CssAny | "auto" | "auto-3d" | CssVars | AnyString | AnyNumber

export type SaturateValue = string | number | CssVars | AnyString

export type ScaleValue = CssAny | "auto" | CssVars | AnyString | AnyNumber

export type ScaleXValue = string | number | CssVars | AnyString

export type ScaleYValue = string | number | CssVars | AnyString

export type ScrollBehaviorValue = string | number | CssVars | AnyString

export type ScrollSnapAlignValue = PropertyValueMap["scrollSnapAlign"]

export type ScrollSnapCoordinateValue = string | number | CssVars | AnyString

export type ScrollSnapDestinationValue = PropertyValueMap["scrollSnapDestination"]

export type ScrollSnapPointsXValue = string | number | CssVars | AnyString

export type ScrollSnapPointsYValue = string | number | CssVars | AnyString

export type ScrollSnapStopValue = PropertyValueMap["scrollSnapStop"] | AnyString

export type ScrollSnapStrictnessValue = "mandatory" | "proximity" | CssVars | AnyString | AnyNumber

export type ScrollSnapTypeValue = "both" | "none" | "x" | "y" | CssVars | AnyString | AnyNumber

export type ScrollSnapTypeXValue = string | number | CssVars | AnyString

export type ScrollSnapTypeYValue = string | number | CssVars | AnyString

export type ScrollTimelineAxisValue = string | number | CssVars | AnyString

export type ScrollTimelineNameValue = string | number | CssVars | AnyString

export type ScrollTimelineValue = string | number | CssVars | AnyString

export type ScrollbarColorValue = string | number | CssVars | AnyString

export type ScrollbarGutterValue = PropertyValueMap["scrollbarGutter"]

export type ScrollbarValue = "hidden" | "visible" | CssVars | AnyString | AnyNumber

export type ScrollbarWidthValue = WithEscapeHatch<PropertyValueMap["scrollbarWidth"]>

export type SepiaValue = string | number | CssVars | AnyString

export type ShadowsValue = CssGlobals | TokenValue<"shadows"> | CssVars | AnyString | AnyNumber

export type SpacingValue = CssAutoGlobals | TokenValue<"spacing"> | "auto" | CssVars | AnyString | AnyNumber

export type SrOnlyValue = boolean | CssVars | AnyString | AnyNumber

export type StrokeDasharrayValue = string | number | CssVars | AnyString

export type StrokeDashoffsetValue = string | number | CssVars | AnyString

export type StrokeLinecapValue = string | number | CssVars | AnyString

export type StrokeLinejoinValue = PropertyValueMap["strokeLinejoin"] | AnyString

export type StrokeMiterlimitValue = string | number | CssVars | AnyString

export type StrokeOpacityValue = string | number | CssVars | AnyString

export type TableLayoutValue = PropertyValueMap["tableLayout"] | AnyString

export type TextAlignValue = PropertyValueMap["textAlign"] | AnyString

export type TextDecorationStyleValue = PropertyValueMap["textDecorationStyle"] | AnyString

export type TextDecorationThicknessValue = PropertyValueMap["textDecorationThickness"]

export type TextDecorationValue = string | number | CssVars | AnyString

export type TextGradientValue = "to-b" | "to-bl" | "to-br" | "to-l" | "to-r" | "to-t" | "to-tl" | "to-tr" | CssVars | AnyString | AnyNumber

export type TextOverflowValue = PropertyValueMap["textOverflow"]

export type TextSizeAdjustValue = PropertyValueMap["textSizeAdjust"]

export type TextStyleValue = "2xl" | "3xl" | "4xl" | "5xl" | "6xl" | "7xl" | "8xl" | "9xl" | "lg" | "md" | "sm" | "xl" | "xs" | CssVars | AnyString | AnyNumber

export type TextTransformValue = PropertyValueMap["textTransform"]

export type TextUnderlineOffsetValue = PropertyValueMap["textUnderlineOffset"]

export type TextWrapValue = PropertyValueMap["textWrap"]

export type TouchActionValue = WithEscapeHatch<CssGlobals | CssVars | OnlyKnown<"-ms-manipulation" | "-ms-none" | "-ms-pan-x" | "-ms-pan-y" | "-ms-pinch-zoom" | "auto" | "manipulation" | "none" | "pan-down" | "pan-left" | "pan-right" | "pan-up" | "pan-x" | "pan-y" | "pinch-zoom">>

export type TransformBoxValue = WithEscapeHatch<PropertyValueMap["transformBox"]>

export type TransformOriginValue = string | number | CssVars | AnyString

export type TransformStyleValue = string | number | CssVars | AnyString

export type TransformValue = string | number | CssVars | AnyString

export type TransitionPropertyValue = "background" | "colors" | "common" | "position" | "size" | CssVars | AnyString | AnyNumber

export type TransitionValue = "all" | "background" | "colors" | "common" | "opacity" | "position" | "shadow" | "size" | "transform" | CssVars | AnyString | AnyNumber

export type TranslateValue = CssAny | "auto" | "auto-3d" | CssVars | AnyString | AnyNumber

export type TranslateXValue = CssAutoGlobals | TokenValue<"spacing"> | "-1/2" | "-1/3" | "-1/4" | "-2/3" | "-2/4" | "-3/4" | "-full" | "1/2" | "1/3" | "1/4" | "2/3" | "2/4" | "3/4" | "full" | CssVars | AnyString | AnyNumber

export type TranslateYValue = CssAutoGlobals | TokenValue<"spacing"> | "-1/2" | "-1/3" | "-1/4" | "-2/3" | "-2/4" | "-3/4" | "-full" | "1/2" | "1/3" | "1/4" | "2/3" | "2/4" | "3/4" | "full" | CssVars | AnyString | AnyNumber

export type TruncateValue = boolean | CssVars | AnyString | AnyNumber

export type UserSelectValue = WithEscapeHatch<PropertyValueMap["userSelect"]>

export type VerticalAlignValue = PropertyValueMap["verticalAlign"]

export type VisibilityValue = string | number | CssVars | AnyString

export type WidthValue = CssDimensionGlobals | TokenValue<"sizes"> | "1/12" | "1/2" | "1/3" | "1/4" | "1/5" | "1/6" | "10/12" | "11/12" | "2/12" | "2/3" | "2/4" | "2/5" | "2/6" | "3/12" | "3/4" | "3/5" | "3/6" | "4/12" | "4/5" | "4/6" | "5/12" | "5/6" | "6/12" | "7/12" | "8/12" | "9/12" | "auto" | "screen" | CssVars | AnyString | AnyNumber

export type WordBreakValue = WithEscapeHatch<PropertyValueMap["wordBreak"]>

export type ZIndexValue = PropertyValueMap["zIndex"]

export type AtRuleType = "media" | "layer" | "container" | "supports" | "page" | "scope" | "starting-style"

export type Selector = `${string}&` | `&${string}` | `@${AtRuleType}${string}`

export type AnySelector = Selector | string

export type CssGlobals = "inherit" | "initial" | "revert" | "revert-layer" | "unset"

export type CssColorGlobals = CssGlobals | "currentColor" | "transparent"

export type CssDimensionGlobals = CssGlobals | "auto" | "fit-content" | "max-content" | "min-content"

export type CssAutoGlobals = CssGlobals | "auto"

export type CssLength = CssDimensionGlobals | (string & {}) | number

export type CssNamedColor = "aliceblue" | "antiquewhite" | "aqua" | "aquamarine" | "azure" | "beige" | "bisque" | "black" | "blanchedalmond" | "blue" | "blueviolet" | "brown" | "burlywood" | "cadetblue" | "chartreuse" | "chocolate" | "coral" | "cornflowerblue" | "cornsilk" | "crimson" | "cyan" | "darkblue" | "darkcyan" | "darkgoldenrod" | "darkgray" | "darkgreen" | "darkgrey" | "darkkhaki" | "darkmagenta" | "darkolivegreen" | "darkorange" | "darkorchid" | "darkred" | "darksalmon" | "darkseagreen" | "darkslateblue" | "darkslategray" | "darkslategrey" | "darkturquoise" | "darkviolet" | "deeppink" | "deepskyblue" | "dimgray" | "dimgrey" | "dodgerblue" | "firebrick" | "floralwhite" | "forestgreen" | "fuchsia" | "gainsboro" | "ghostwhite" | "gold" | "goldenrod" | "gray" | "green" | "greenyellow" | "grey" | "honeydew" | "hotpink" | "indianred" | "indigo" | "ivory" | "khaki" | "lavender" | "lavenderblush" | "lawngreen" | "lemonchiffon" | "lightblue" | "lightcoral" | "lightcyan" | "lightgoldenrodyellow" | "lightgray" | "lightgreen" | "lightgrey" | "lightpink" | "lightsalmon" | "lightseagreen" | "lightskyblue" | "lightslategray" | "lightslategrey" | "lightsteelblue" | "lightyellow" | "lime" | "limegreen" | "linen" | "magenta" | "maroon" | "mediumaquamarine" | "mediumblue" | "mediumorchid" | "mediumpurple" | "mediumseagreen" | "mediumslateblue" | "mediumspringgreen" | "mediumturquoise" | "mediumvioletred" | "midnightblue" | "mintcream" | "mistyrose" | "moccasin" | "navajowhite" | "navy" | "oldlace" | "olive" | "olivedrab" | "orange" | "orangered" | "orchid" | "palegoldenrod" | "palegreen" | "paleturquoise" | "palevioletred" | "papayawhip" | "peachpuff" | "peru" | "pink" | "plum" | "powderblue" | "purple" | "rebeccapurple" | "red" | "rosybrown" | "royalblue" | "saddlebrown" | "salmon" | "sandybrown" | "seagreen" | "seashell" | "sienna" | "silver" | "skyblue" | "slateblue" | "slategray" | "slategrey" | "snow" | "springgreen" | "steelblue" | "tan" | "teal" | "thistle" | "tomato" | "turquoise" | "violet" | "wheat" | "white" | "whitesmoke" | "yellow" | "yellowgreen"

export type CssSystemColor = "AccentColor" | "AccentColorText" | "ActiveText" | "ButtonBorder" | "ButtonFace" | "ButtonText" | "Canvas" | "CanvasText" | "Field" | "FieldText" | "GrayText" | "Highlight" | "HighlightText" | "LinkText" | "Mark" | "MarkText" | "SelectedItem" | "SelectedItemText" | "VisitedText"

export type CssColor = CssColorGlobals | CssNamedColor | CssSystemColor | (string & {})

export type CssAbsoluteSize = "large" | "medium" | "small" | "x-large" | "x-small" | "xx-large" | "xx-small" | "xxx-large"

export type CssBgSize = CssLength | "auto" | "contain" | "cover"

export type CssFontSize = CssLength | CssAbsoluteSize | "larger" | "smaller" | "math"

export type CssFontWeight = CssGlobals | "bold" | "normal" | (string & {}) | number

export type CssLineWidth = CssLength | "medium" | "thick" | "thin"

export type CssLineStyle = CssGlobals | "dashed" | "dotted" | "double" | "groove" | "hidden" | "inset" | "none" | "outset" | "ridge" | "solid"

export type CssLineStyleOpen = CssLineStyle | (string & {})

export type CssRepeatStyle = CssGlobals | "no-repeat" | "repeat" | "repeat-x" | "repeat-y" | "round" | "space" | (string & {})

export type CssFontStretch = CssGlobals | "condensed" | "expanded" | "extra-condensed" | "extra-expanded" | "normal" | "semi-condensed" | "semi-expanded" | "ultra-condensed" | "ultra-expanded" | (string & {})

export type CssOverflow = CssGlobals | "-moz-hidden-unscrollable" | "auto" | "clip" | "hidden" | "overlay" | "scroll" | "visible"

export type CssOverflowOpen = CssOverflow | (string & {})

export type CssOverflowShort = CssGlobals | "auto" | "clip" | "hidden" | "scroll" | "visible"

export type CssOverscrollBehavior = CssGlobals | "auto" | "contain" | "none"

export type CssOverscrollBehaviorOpen = CssOverscrollBehavior | (string & {})

export type CssPosition = CssLength | "bottom" | "center" | "left" | "right" | "top"

export type CssGenericFamily = "-apple-system" | "cursive" | "emoji" | "fangsong" | "fantasy" | "math" | "monospace" | "sans-serif" | "serif" | "system-ui" | "ui-monospace" | "ui-rounded" | "ui-sans-serif" | "ui-serif"

export type CssFontFamily = CssGlobals | CssGenericFamily | (string & {})

export type CssNumeric = CssGlobals | (string & {}) | number

export type CssZIndex = CssAutoGlobals | (string & {}) | number

export type CssAny = CssGlobals | (string & {}) | number

export type PropertyValueMap = {
  accentColor: CssColor
  alignContent: CssGlobals | "baseline" | "center" | "end" | "first baseline" | "flex-end" | "flex-start" | "last baseline" | "normal" | "safe center" | "safe end" | "safe start" | "space-around" | "space-between" | "space-evenly" | "start" | "stretch" | (string & {})
  alignItems: CssGlobals | "anchor-center" | "baseline" | "center" | "end" | "first baseline" | "flex-end" | "flex-start" | "last baseline" | "normal" | "safe center" | "safe end" | "safe start" | "self-end" | "self-start" | "start" | "stretch" | (string & {})
  alignSelf: CssGlobals | "anchor-center" | "auto" | "baseline" | "center" | "end" | "first baseline" | "flex-end" | "flex-start" | "last baseline" | "normal" | "safe center" | "safe end" | "safe start" | "self-end" | "self-start" | "start" | "stretch" | (string & {})
  alignmentBaseline: CssGlobals | "alphabetic" | "baseline" | "central" | "ideographic" | "mathematical" | "middle" | "text-after-edge" | "text-before-edge"
  appearance: CssGlobals | "auto" | "button" | "checkbox" | "listbox" | "menulist" | "menulist-button" | "meter" | "none" | "progress-bar" | "radio" | "searchfield" | "textarea" | "textfield"
  backgroundColor: CssColor
  backgroundPosition: CssPosition
  backgroundRepeat: CssRepeatStyle
  backgroundSize: CssBgSize
  baselineShift: CssLength
  blockSize: CssLength
  borderBlockColor: CssColor
  borderBlockEndColor: CssColor
  borderBlockEndStyle: CssLineStyle
  borderBlockEndWidth: CssLineWidth
  borderBlockStartColor: CssColor
  borderBlockStartStyle: CssLineStyle
  borderBlockStartWidth: CssLineWidth
  borderBlockStyle: CssLineStyleOpen
  borderBlockWidth: CssLineWidth
  borderBottomColor: CssColor
  borderBottomLeftRadius: CssLength
  borderBottomRightRadius: CssLength
  borderBottomStyle: CssLineStyle
  borderBottomWidth: CssLineWidth
  borderColor: CssColor
  borderEndEndRadius: CssLength
  borderEndStartRadius: CssLength
  borderImage: CssGlobals | "none" | "repeat" | "round" | "space" | "stretch" | (string & {})
  borderImageWidth: CssLength
  borderInlineColor: CssColor
  borderInlineEndColor: CssColor
  borderInlineEndStyle: CssLineStyle
  borderInlineEndWidth: CssLineWidth
  borderInlineStartColor: CssColor
  borderInlineStartStyle: CssLineStyle
  borderInlineStartWidth: CssLineWidth
  borderInlineStyle: CssLineStyleOpen
  borderInlineWidth: CssLineWidth
  borderLeftColor: CssColor
  borderLeftStyle: CssLineStyle
  borderLeftWidth: CssLineWidth
  borderRadius: CssLength
  borderRightColor: CssColor
  borderRightStyle: CssLineStyle
  borderRightWidth: CssLineWidth
  borderSpacing: CssLength
  borderStartEndRadius: CssLength
  borderStartStartRadius: CssLength
  borderStyle: CssLineStyleOpen
  borderTopColor: CssColor
  borderTopLeftRadius: CssLength
  borderTopRightRadius: CssLength
  borderTopStyle: CssLineStyle
  borderTopWidth: CssLineWidth
  borderWidth: CssLineWidth
  bottom: CssLength
  boxAlign: CssGlobals | "baseline" | "center" | "end" | "start" | "stretch"
  breakAfter: CssGlobals | "all" | "always" | "auto" | "avoid" | "avoid-column" | "avoid-page" | "avoid-region" | "column" | "left" | "page" | "recto" | "region" | "right" | "verso"
  breakBefore: PropertyValueMap["breakAfter"]
  breakInside: CssGlobals | "auto" | "avoid" | "avoid-column" | "avoid-page" | "avoid-region"
  caret: CssGlobals | "auto" | "bar" | "block" | "currentColor" | "underscore" | (string & {})
  caretColor: CssColor
  clear: CssGlobals | "both" | "inline-end" | "inline-start" | "left" | "none" | "right"
  clipPath: CssGlobals | "fill-box" | "margin-box" | "none" | "stroke-box" | "view-box" | (string & {})
  color: CssColor
  colorScheme: CssGlobals | "dark" | "light" | "normal" | (string & {})
  columnGap: CssLength
  columnRuleColor: CssColor
  columnRuleStyle: CssLineStyleOpen
  columnRuleWidth: CssLineWidth
  columnWidth: CssLength
  contain: CssGlobals | "content" | "inline-size" | "layout" | "none" | "paint" | "size" | "strict" | "style" | (string & {})
  containIntrinsicBlockSize: CssLength
  containIntrinsicHeight: CssLength
  containIntrinsicInlineSize: CssLength
  containIntrinsicSize: CssLength
  containIntrinsicWidth: CssLength
  containerType: CssGlobals | "inline-size" | "normal" | "scroll-state" | "size" | (string & {})
  content: CssGlobals | "close-quote" | "no-close-quote" | "no-open-quote" | "none" | "normal" | "open-quote" | (string & {})
  cursor: CssGlobals | "-moz-grab" | "-moz-zoom-in" | "-moz-zoom-out" | "-webkit-grab" | "-webkit-grabbing" | "-webkit-zoom-in" | "-webkit-zoom-out" | "alias" | "all-scroll" | "auto" | "cell" | "col-resize" | "context-menu" | "copy" | "crosshair" | "default" | "e-resize" | "ew-resize" | "grab" | "grabbing" | "help" | "move" | "n-resize" | "ne-resize" | "nesw-resize" | "no-drop" | "none" | "not-allowed" | "ns-resize" | "nw-resize" | "nwse-resize" | "pointer" | "progress" | "row-resize" | "s-resize" | "se-resize" | "sw-resize" | "text" | "vertical-text" | "w-resize" | "wait" | "zoom-in" | "zoom-out" | (string & {})
  display: CssGlobals | "-ms-flexbox" | "-ms-grid" | "-ms-inline-flexbox" | "-ms-inline-grid" | "-webkit-flex" | "-webkit-inline-flex" | "block" | "contents" | "flex" | "flow" | "flow-root" | "grid" | "inline" | "inline-block" | "inline-flex" | "inline-grid" | "inline-list-item" | "inline-table" | "list-item" | "none" | "ruby" | "ruby-base" | "ruby-base-container" | "ruby-text" | "ruby-text-container" | "run-in" | "table" | "table-caption" | "table-cell" | "table-column" | "table-column-group" | "table-footer-group" | "table-header-group" | "table-row" | "table-row-group" | (string & {})
  dominantBaseline: CssGlobals | "alphabetic" | "auto" | "central" | "hanging" | "ideographic" | "mathematical" | "middle" | "text-bottom" | "text-top"
  fieldSizing: CssGlobals | "content" | "fixed"
  flexBasis: CssLength
  flexFlow: CssGlobals | "column" | "column-reverse" | "nowrap" | "row" | "row-reverse" | "wrap" | "wrap-reverse" | (string & {})
  flexGrow: CssNumeric
  flexShrink: CssNumeric
  float: CssGlobals | "inline-end" | "inline-start" | "left" | "none" | "right"
  floodColor: CssColor
  font: CssGlobals | "caption" | "icon" | "menu" | "message-box" | "small-caption" | "status-bar" | (string & {})
  fontFamily: CssFontFamily
  fontSize: CssFontSize
  fontSizeAdjust: CssLength
  fontStretch: CssFontStretch
  fontStyle: CssGlobals | "italic" | "normal" | "oblique" | (string & {})
  fontSynthesis: CssGlobals | "none" | "position" | "small-caps" | "style" | "weight" | (string & {})
  fontVariantCaps: CssGlobals | "all-petite-caps" | "all-small-caps" | "normal" | "petite-caps" | "small-caps" | "titling-caps" | "unicase"
  fontVariantEastAsian: CssGlobals | "full-width" | "jis04" | "jis78" | "jis83" | "jis90" | "normal" | "proportional-width" | "ruby" | "simplified" | "traditional" | (string & {})
  fontVariantLigatures: CssGlobals | "common-ligatures" | "contextual" | "discretionary-ligatures" | "historical-ligatures" | "no-common-ligatures" | "no-contextual" | "no-discretionary-ligatures" | "no-historical-ligatures" | "none" | "normal" | (string & {})
  fontVariantNumeric: CssGlobals | "diagonal-fractions" | "lining-nums" | "normal" | "oldstyle-nums" | "ordinal" | "proportional-nums" | "slashed-zero" | "stacked-fractions" | "tabular-nums" | (string & {})
  fontWeight: CssFontWeight
  fontWidth: CssFontStretch
  gap: CssLength
  gridAutoFlow: CssGlobals | "column" | "dense" | "row" | (string & {})
  gridColumnGap: CssLength
  gridGap: CssLength
  gridRowGap: CssLength
  hangingPunctuation: CssGlobals | "allow-end" | "first" | "force-end" | "last" | "none" | (string & {})
  height: CssLength
  hyphens: CssGlobals | "auto" | "manual" | "none"
  imageRendering: CssGlobals | "-moz-crisp-edges" | "-webkit-optimize-contrast" | "auto" | "crisp-edges" | "pixelated" | "smooth"
  imeMode: CssGlobals | "active" | "auto" | "disabled" | "inactive" | "normal"
  inlineSize: CssLength
  inset: CssLength
  justifyContent: CssGlobals | "center" | "end" | "flex-end" | "flex-start" | "left" | "normal" | "right" | "safe center" | "safe end" | "safe start" | "space-around" | "space-between" | "space-evenly" | "start" | "stretch" | (string & {})
  justifyItems: CssGlobals | "anchor-center" | "baseline" | "center" | "end" | "first baseline" | "flex-end" | "flex-start" | "last baseline" | "left" | "legacy" | "normal" | "right" | "safe center" | "safe end" | "safe start" | "self-end" | "self-start" | "start" | "stretch" | (string & {})
  justifySelf: CssGlobals | "anchor-center" | "auto" | "baseline" | "center" | "end" | "first baseline" | "flex-end" | "flex-start" | "last baseline" | "left" | "normal" | "right" | "safe center" | "safe end" | "safe start" | "self-end" | "self-start" | "start" | "stretch" | (string & {})
  left: CssLength
  letterSpacing: CssLength
  lightingColor: CssColor
  lineBreak: CssGlobals | "anywhere" | "auto" | "loose" | "normal" | "strict"
  lineHeight: CssNumeric
  lineHeightStep: CssLength
  listStyleImage: CssGlobals | "none" | (string & {})
  listStylePosition: CssGlobals | "inside" | "outside"
  listStyleType: PropertyValueMap["listStyleImage"]
  margin: CssLength
  marginBottom: CssLength
  marginLeft: CssLength
  marginRight: CssLength
  marginTop: CssLength
  maskBorder: CssGlobals | "alpha" | "luminance" | "none" | "repeat" | "round" | "space" | "stretch" | (string & {})
  maskBorderWidth: CssLength
  maskPosition: CssPosition
  maskRepeat: CssRepeatStyle
  maskSize: CssBgSize
  maxBlockSize: CssLength
  maxHeight: CssLength
  maxInlineSize: CssLength
  maxWidth: CssLength
  minBlockSize: CssLength
  minHeight: CssLength
  minInlineSize: CssLength
  minWidth: CssLength
  mixBlendMode: CssGlobals | "color" | "color-burn" | "color-dodge" | "darken" | "difference" | "exclusion" | "hard-light" | "hue" | "lighten" | "luminosity" | "multiply" | "normal" | "overlay" | "plus-darker" | "plus-lighter" | "saturation" | "screen" | "soft-light"
  objectFit: CssGlobals | "contain" | "cover" | "fill" | "none" | "scale-down"
  objectPosition: CssPosition
  offsetAnchor: CssPosition
  offsetPosition: CssPosition
  opacity: CssNumeric
  order: CssNumeric
  orphans: CssNumeric
  outlineColor: CssColor
  outlineOffset: CssLength
  outlineStyle: CssGlobals | "auto" | "dashed" | "dotted" | "double" | "groove" | "inset" | "none" | "outset" | "ridge" | "solid"
  outlineWidth: CssLineWidth
  overflow: CssOverflowOpen
  overflowBlock: CssOverflowShort
  overflowClipMargin: CssLength
  overflowInline: CssOverflowShort
  overflowX: CssOverflow
  overflowY: CssOverflow
  overscrollBehavior: CssOverscrollBehaviorOpen
  overscrollBehaviorX: CssOverscrollBehavior
  overscrollBehaviorY: CssOverscrollBehavior
  padding: CssLength
  paddingBottom: CssLength
  paddingLeft: CssLength
  paddingRight: CssLength
  paddingTop: CssLength
  pageBreakAfter: CssGlobals | "always" | "auto" | "avoid" | "left" | "recto" | "right" | "verso"
  pageBreakBefore: PropertyValueMap["pageBreakAfter"]
  perspective: CssLength
  perspectiveOrigin: CssPosition
  placeContent: PropertyValueMap["alignContent"]
  placeItems: PropertyValueMap["alignItems"]
  placeSelf: PropertyValueMap["alignSelf"]
  pointerEvents: CssGlobals | "all" | "auto" | "fill" | "none" | "painted" | "stroke" | "visible" | "visibleFill" | "visiblePainted" | "visibleStroke"
  position: CssGlobals | "-webkit-sticky" | "absolute" | "fixed" | "relative" | "static" | "sticky"
  positionTryOrder: CssGlobals | "most-block-size" | "most-height" | "most-inline-size" | "most-width" | "normal"
  resize: CssGlobals | "block" | "both" | "horizontal" | "inline" | "none" | "vertical"
  right: CssLength
  rowGap: CssLength
  scrollMarginBlock: CssLength
  scrollMarginBlockEnd: CssLength
  scrollMarginBlockStart: CssLength
  scrollMarginBottom: CssLength
  scrollMarginInline: CssLength
  scrollMarginInlineEnd: CssLength
  scrollMarginInlineStart: CssLength
  scrollMarginLeft: CssLength
  scrollMarginRight: CssLength
  scrollMarginTop: CssLength
  scrollPaddingBlock: CssLength
  scrollPaddingBlockEnd: CssLength
  scrollPaddingBlockStart: CssLength
  scrollPaddingBottom: CssLength
  scrollPaddingInline: CssLength
  scrollPaddingInlineEnd: CssLength
  scrollPaddingInlineStart: CssLength
  scrollPaddingLeft: CssLength
  scrollPaddingRight: CssLength
  scrollPaddingTop: CssLength
  scrollSnapAlign: CssGlobals | "center" | "end" | "none" | "start" | (string & {})
  scrollSnapDestination: CssPosition
  scrollSnapStop: CssGlobals | "always" | "normal"
  scrollSnapType: CssGlobals | "block" | "both" | "inline" | "none" | "x" | "y" | (string & {})
  scrollbarGutter: CssGlobals | "auto" | "stable" | "stable both-edges" | (string & {})
  scrollbarWidth: CssGlobals | "auto" | "none" | "thin"
  shapeMargin: CssLength
  shapeOutside: CssGlobals | "border-box" | "content-box" | "margin-box" | "none" | "padding-box" | (string & {})
  speakAs: CssGlobals | "digits" | "literal-punctuation" | "no-punctuation" | "normal" | "spell-out" | (string & {})
  stopColor: CssColor
  strokeColor: CssColor
  strokeLinejoin: CssGlobals | "arcs" | "bevel" | "miter" | "miter-clip" | "round"
  strokeWidth: CssLength
  tabSize: CssLength
  tableLayout: CssGlobals | "auto" | "fixed"
  textAlign: CssGlobals | "-khtml-center" | "-khtml-left" | "-khtml-right" | "-moz-center" | "-moz-left" | "-moz-right" | "-webkit-center" | "-webkit-left" | "-webkit-match-parent" | "-webkit-right" | "center" | "end" | "justify" | "left" | "match-parent" | "right" | "start"
  textAlignLast: CssGlobals | "auto" | "center" | "end" | "justify" | "left" | "right" | "start"
  textAutospace: CssGlobals | "auto" | "ideograph-alpha" | "ideograph-numeric" | "insert" | "no-autospace" | "normal" | "punctuation" | "replace" | (string & {})
  textBox: CssGlobals | "auto" | "cap" | "ex" | "ideographic" | "ideographic-ink" | "none" | "normal" | "text" | "trim-both" | "trim-end" | "trim-start" | (string & {})
  textBoxEdge: CssGlobals | "auto" | "cap" | "ex" | "ideographic" | "ideographic-ink" | "text" | (string & {})
  textDecorationColor: CssColor
  textDecorationLine: CssGlobals | "blink" | "grammar-error" | "line-through" | "none" | "overline" | "spelling-error" | "underline" | (string & {})
  textDecorationSkip: CssGlobals | "box-decoration" | "edges" | "leading-spaces" | "none" | "objects" | "spaces" | "trailing-spaces" | (string & {})
  textDecorationSkipInk: CssGlobals | "all" | "auto" | "none"
  textDecorationStyle: CssGlobals | "dashed" | "dotted" | "double" | "solid" | "wavy"
  textDecorationThickness: CssLength
  textEmphasis: CssGlobals | "circle" | "currentColor" | "dot" | "double-circle" | "filled" | "none" | "open" | "sesame" | "triangle" | (string & {})
  textEmphasisColor: CssColor
  textEmphasisPosition: CssGlobals | "auto" | "over" | "under" | (string & {})
  textEmphasisStyle: CssGlobals | "circle" | "dot" | "double-circle" | "filled" | "none" | "open" | "sesame" | "triangle" | (string & {})
  textIndent: CssLength
  textJustify: CssGlobals | "auto" | "distribute" | "inter-character" | "inter-word" | "none"
  textOrientation: CssGlobals | "mixed" | "sideways" | "sideways-right" | "upright"
  textOverflow: CssGlobals | "clip" | "ellipsis" | (string & {})
  textRendering: CssGlobals | "auto" | "geometricPrecision" | "optimizeLegibility" | "optimizeSpeed"
  textSizeAdjust: CssLength
  textTransform: CssGlobals | "capitalize" | "full-size-kana" | "full-width" | "lowercase" | "math-auto" | "none" | "uppercase" | (string & {})
  textUnderlineOffset: CssLength
  textUnderlinePosition: CssGlobals | "auto" | "from-font" | "left" | "right" | "under" | (string & {})
  textWrap: CssGlobals | "auto" | "balance" | "nowrap" | "pretty" | "stable" | "wrap" | (string & {})
  textWrapMode: CssGlobals | "nowrap" | "wrap"
  textWrapStyle: CssGlobals | "auto" | "balance" | "pretty" | "stable"
  top: CssLength
  touchAction: CssGlobals | "-ms-manipulation" | "-ms-none" | "-ms-pan-x" | "-ms-pan-y" | "-ms-pinch-zoom" | "auto" | "manipulation" | "none" | "pan-down" | "pan-left" | "pan-right" | "pan-up" | "pan-x" | "pan-y" | "pinch-zoom" | (string & {})
  transformBox: CssGlobals | "border-box" | "content-box" | "fill-box" | "stroke-box" | "view-box"
  transitionBehavior: CssGlobals | "allow-discrete" | "normal" | (string & {})
  transitionProperty: CssGlobals | "all" | "none" | (string & {})
  unicodeBidi: CssGlobals | "-moz-isolate" | "-moz-isolate-override" | "-moz-plaintext" | "-webkit-isolate" | "-webkit-isolate-override" | "-webkit-plaintext" | "bidi-override" | "embed" | "isolate" | "isolate-override" | "normal" | "plaintext"
  userSelect: CssGlobals | "-moz-none" | "all" | "auto" | "none" | "text"
  vectorEffect: CssGlobals | "fixed-position" | "non-rotation" | "non-scaling-size" | "non-scaling-stroke" | "none"
  verticalAlign: CssGlobals | "baseline" | "bottom" | "middle" | "sub" | "super" | "text-bottom" | "text-top" | "top" | (string & {})
  viewTimelineInset: CssLength
  whiteSpace: CssGlobals | "-moz-pre-wrap" | "break-spaces" | "collapse" | "normal" | "nowrap" | "pre" | "pre-line" | "pre-wrap" | "preserve" | "preserve-breaks" | "preserve-spaces" | "wrap" | (string & {})
  whiteSpaceCollapse: CssGlobals | "break-spaces" | "collapse" | "preserve" | "preserve-breaks" | "preserve-spaces"
  widows: CssNumeric
  width: CssLength
  wordBreak: CssGlobals | "auto-phrase" | "break-all" | "break-word" | "keep-all" | "normal"
  wordSpacing: CssLength
  wordWrap: CssGlobals | "break-word" | "normal"
  writingMode: CssGlobals | "horizontal-tb" | "sideways-lr" | "sideways-rl" | "vertical-lr" | "vertical-rl"
  zIndex: CssZIndex
  zoom: CssNumeric
}

export interface SystemProperties {
  WebkitAppearance?: ConditionalValue<CssAny>
  WebkitBorderBefore?: ConditionalValue<CssAny>
  WebkitBorderBeforeColor?: ConditionalValue<CssAny>
  WebkitBorderBeforeStyle?: ConditionalValue<CssAny>
  WebkitBorderBeforeWidth?: ConditionalValue<CssAny>
  WebkitBoxReflect?: ConditionalValue<CssAny>
  WebkitLineClamp?: ConditionalValue<CssAny>
  WebkitMask?: ConditionalValue<CssAny>
  WebkitMaskAttachment?: ConditionalValue<CssAny>
  WebkitMaskClip?: ConditionalValue<CssAny>
  WebkitMaskComposite?: ConditionalValue<CssAny>
  WebkitMaskImage?: ConditionalValue<CssAny>
  WebkitMaskOrigin?: ConditionalValue<CssAny>
  WebkitMaskPosition?: ConditionalValue<CssAny>
  WebkitMaskPositionX?: ConditionalValue<CssAny>
  WebkitMaskPositionY?: ConditionalValue<CssAny>
  WebkitMaskRepeat?: ConditionalValue<CssAny>
  WebkitMaskRepeatX?: ConditionalValue<CssAny>
  WebkitMaskRepeatY?: ConditionalValue<CssAny>
  WebkitMaskSize?: ConditionalValue<CssAny>
  WebkitOverflowScrolling?: ConditionalValue<CssAny>
  WebkitTapHighlightColor?: ConditionalValue<CssAny>
  WebkitTextFillColor?: ConditionalValue<ColorsValue | CssAny>
  WebkitTextStroke?: ConditionalValue<CssAny>
  WebkitTextStrokeColor?: ConditionalValue<CssAny>
  WebkitTextStrokeWidth?: ConditionalValue<CssAny>
  WebkitTouchCallout?: ConditionalValue<CssAny>
  WebkitUserModify?: ConditionalValue<CssAny>
  WebkitUserSelect?: ConditionalValue<CssAny>
  accentColor?: ConditionalValue<ColorsValue | PropertyValueMap["accentColor"]>
  alignContent?: ConditionalValue<AlignContentValue>
  alignItems?: ConditionalValue<AlignItemsValue>
  alignSelf?: ConditionalValue<AlignSelfValue>
  alignTracks?: ConditionalValue<CssAny>
  alignmentBaseline?: ConditionalValue<PropertyValueMap["alignmentBaseline"] | AnyString>
  all?: ConditionalValue<CssAny>
  anchorName?: ConditionalValue<CssAny>
  anchorScope?: ConditionalValue<CssAny>
  animation?: ConditionalValue<AnimationsValue | CssAny>
  animationComposition?: ConditionalValue<AnimationCompositionValue | CssAny>
  animationDelay?: ConditionalValue<DurationsValue | CssAny>
  animationDirection?: ConditionalValue<AnimationDirectionValue | CssAny>
  animationDuration?: ConditionalValue<DurationsValue | CssAny>
  animationFillMode?: ConditionalValue<AnimationFillModeValue | CssAny>
  animationIterationCount?: ConditionalValue<AnimationIterationCountValue | CssAny>
  animationName?: ConditionalValue<KeyframesValue | CssAny>
  animationPlayState?: ConditionalValue<AnimationPlayStateValue | CssAny>
  animationRange?: ConditionalValue<AnimationRangeValue | CssAny>
  animationRangeEnd?: ConditionalValue<AnimationRangeEndValue | CssAny>
  animationRangeStart?: ConditionalValue<AnimationRangeStartValue | CssAny>
  animationTimeline?: ConditionalValue<AnimationTimelineValue | CssAny>
  animationTimingFunction?: ConditionalValue<EasingsValue | CssAny>
  appearance?: ConditionalValue<AppearanceValue>
  aspectRatio?: ConditionalValue<AspectRatiosValue | CssAny>
  backdropFilter?: ConditionalValue<BackdropFilterValue | CssAny>
  backfaceVisibility?: ConditionalValue<BackfaceVisibilityValue | CssAny>
  background?: ConditionalValue<ColorsValue | CssAny>
  backgroundAttachment?: ConditionalValue<BackgroundAttachmentValue | CssAny>
  backgroundBlendMode?: ConditionalValue<BackgroundBlendModeValue | CssAny>
  backgroundClip?: ConditionalValue<BackgroundClipValue | CssAny>
  backgroundColor?: ConditionalValue<ColorsValue | PropertyValueMap["backgroundColor"]>
  backgroundImage?: ConditionalValue<AssetsValue | CssAny>
  backgroundOrigin?: ConditionalValue<BackgroundOriginValue | CssAny>
  backgroundPosition?: ConditionalValue<BackgroundPositionValue>
  backgroundPositionX?: ConditionalValue<BackgroundPositionXValue | CssAny>
  backgroundPositionY?: ConditionalValue<BackgroundPositionYValue | CssAny>
  backgroundRepeat?: ConditionalValue<BackgroundRepeatValue>
  backgroundSize?: ConditionalValue<BackgroundSizeValue>
  baselineShift?: ConditionalValue<PropertyValueMap["baselineShift"]>
  blockSize?: ConditionalValue<BlockSizeValue | PropertyValueMap["blockSize"]>
  border?: ConditionalValue<BordersValue | CssAny>
  borderBlock?: ConditionalValue<BordersValue | CssAny>
  borderBlockColor?: ConditionalValue<ColorsValue | PropertyValueMap["borderBlockColor"]>
  borderBlockEnd?: ConditionalValue<BordersValue | CssAny>
  borderBlockEndColor?: ConditionalValue<ColorsValue | PropertyValueMap["borderBlockEndColor"]>
  borderBlockEndStyle?: ConditionalValue<WithEscapeHatch<PropertyValueMap["borderBlockEndStyle"]>>
  borderBlockEndWidth?: ConditionalValue<BorderWidthsValue>
  borderBlockStart?: ConditionalValue<BordersValue | CssAny>
  borderBlockStartColor?: ConditionalValue<ColorsValue | PropertyValueMap["borderBlockStartColor"]>
  borderBlockStartStyle?: ConditionalValue<WithEscapeHatch<PropertyValueMap["borderBlockStartStyle"]>>
  borderBlockStartWidth?: ConditionalValue<BorderWidthsValue>
  borderBlockStyle?: ConditionalValue<WithEscapeHatch<PropertyValueMap["borderBlockStyle"]>>
  borderBlockWidth?: ConditionalValue<BorderWidthsValue>
  borderBottom?: ConditionalValue<BordersValue | CssAny>
  borderBottomColor?: ConditionalValue<ColorsValue | PropertyValueMap["borderBottomColor"]>
  borderBottomLeftRadius?: ConditionalValue<RadiiValue | PropertyValueMap["borderBottomLeftRadius"]>
  borderBottomRightRadius?: ConditionalValue<RadiiValue | PropertyValueMap["borderBottomRightRadius"]>
  borderBottomStyle?: ConditionalValue<WithEscapeHatch<PropertyValueMap["borderBottomStyle"]>>
  borderBottomWidth?: ConditionalValue<BorderWidthsValue>
  borderCollapse?: ConditionalValue<BorderCollapseValue | CssAny>
  borderColor?: ConditionalValue<ColorsValue | PropertyValueMap["borderColor"]>
  borderEndEndRadius?: ConditionalValue<RadiiValue | PropertyValueMap["borderEndEndRadius"]>
  borderEndStartRadius?: ConditionalValue<RadiiValue | PropertyValueMap["borderEndStartRadius"]>
  borderImage?: ConditionalValue<PropertyValueMap["borderImage"]>
  borderImageOutset?: ConditionalValue<CssAny>
  borderImageRepeat?: ConditionalValue<CssAny>
  borderImageSlice?: ConditionalValue<CssAny>
  borderImageSource?: ConditionalValue<CssAny>
  borderImageWidth?: ConditionalValue<PropertyValueMap["borderImageWidth"]>
  borderInline?: ConditionalValue<BordersValue | CssAny>
  borderInlineColor?: ConditionalValue<ColorsValue | PropertyValueMap["borderInlineColor"]>
  borderInlineEnd?: ConditionalValue<BordersValue | CssAny>
  borderInlineEndColor?: ConditionalValue<ColorsValue | PropertyValueMap["borderInlineEndColor"]>
  borderInlineEndStyle?: ConditionalValue<WithEscapeHatch<PropertyValueMap["borderInlineEndStyle"]>>
  borderInlineEndWidth?: ConditionalValue<BorderWidthsValue>
  borderInlineStart?: ConditionalValue<BordersValue | CssAny>
  borderInlineStartColor?: ConditionalValue<ColorsValue | PropertyValueMap["borderInlineStartColor"]>
  borderInlineStartStyle?: ConditionalValue<WithEscapeHatch<PropertyValueMap["borderInlineStartStyle"]>>
  borderInlineStartWidth?: ConditionalValue<BorderWidthsValue>
  borderInlineStyle?: ConditionalValue<WithEscapeHatch<PropertyValueMap["borderInlineStyle"]>>
  borderInlineWidth?: ConditionalValue<BorderWidthsValue>
  borderLeft?: ConditionalValue<BordersValue | CssAny>
  borderLeftColor?: ConditionalValue<ColorsValue | PropertyValueMap["borderLeftColor"]>
  borderLeftStyle?: ConditionalValue<WithEscapeHatch<PropertyValueMap["borderLeftStyle"]>>
  borderLeftWidth?: ConditionalValue<BorderWidthsValue>
  borderRadius?: ConditionalValue<RadiiValue | PropertyValueMap["borderRadius"]>
  borderRight?: ConditionalValue<BordersValue | CssAny>
  borderRightColor?: ConditionalValue<ColorsValue | PropertyValueMap["borderRightColor"]>
  borderRightStyle?: ConditionalValue<WithEscapeHatch<PropertyValueMap["borderRightStyle"]>>
  borderRightWidth?: ConditionalValue<BorderWidthsValue>
  borderSpacing?: ConditionalValue<SpacingValue | PropertyValueMap["borderSpacing"]>
  borderStartEndRadius?: ConditionalValue<RadiiValue | PropertyValueMap["borderStartEndRadius"]>
  borderStartStartRadius?: ConditionalValue<RadiiValue | PropertyValueMap["borderStartStartRadius"]>
  borderStyle?: ConditionalValue<PropertyValueMap["borderStyle"]>
  borderTop?: ConditionalValue<BordersValue | CssAny>
  borderTopColor?: ConditionalValue<ColorsValue | PropertyValueMap["borderTopColor"]>
  borderTopLeftRadius?: ConditionalValue<RadiiValue | PropertyValueMap["borderTopLeftRadius"]>
  borderTopRightRadius?: ConditionalValue<RadiiValue | PropertyValueMap["borderTopRightRadius"]>
  borderTopStyle?: ConditionalValue<WithEscapeHatch<PropertyValueMap["borderTopStyle"]>>
  borderTopWidth?: ConditionalValue<BorderWidthsValue>
  borderWidth?: ConditionalValue<BorderWidthsValue>
  bottom?: ConditionalValue<SpacingValue | PropertyValueMap["bottom"]>
  boxAlign?: ConditionalValue<PropertyValueMap["boxAlign"] | AnyString>
  boxDecorationBreak?: ConditionalValue<BoxDecorationBreakValue | CssAny>
  boxDirection?: ConditionalValue<CssAny>
  boxFlex?: ConditionalValue<CssAny>
  boxFlexGroup?: ConditionalValue<CssAny>
  boxLines?: ConditionalValue<CssAny>
  boxOrdinalGroup?: ConditionalValue<CssAny>
  boxOrient?: ConditionalValue<CssAny>
  boxPack?: ConditionalValue<CssAny>
  boxShadow?: ConditionalValue<ShadowsValue | CssAny>
  boxSizing?: ConditionalValue<BoxSizingValue | CssAny>
  breakAfter?: ConditionalValue<WithEscapeHatch<PropertyValueMap["breakAfter"]>>
  breakBefore?: ConditionalValue<WithEscapeHatch<PropertyValueMap["breakBefore"]>>
  breakInside?: ConditionalValue<WithEscapeHatch<PropertyValueMap["breakInside"]>>
  captionSide?: ConditionalValue<CssAny>
  caret?: ConditionalValue<PropertyValueMap["caret"]>
  caretColor?: ConditionalValue<ColorsValue | PropertyValueMap["caretColor"]>
  caretShape?: ConditionalValue<CssAny>
  clear?: ConditionalValue<WithEscapeHatch<PropertyValueMap["clear"]>>
  clip?: ConditionalValue<CssAny>
  clipPath?: ConditionalValue<ClipPathValue>
  clipRule?: ConditionalValue<CssAny>
  color?: ConditionalValue<ColorsValue | PropertyValueMap["color"]>
  colorInterpolation?: ConditionalValue<CssAny>
  colorInterpolationFilters?: ConditionalValue<CssAny>
  colorRendering?: ConditionalValue<CssAny>
  colorScheme?: ConditionalValue<PropertyValueMap["colorScheme"]>
  columnCount?: ConditionalValue<CssAny>
  columnFill?: ConditionalValue<CssAny>
  columnGap?: ConditionalValue<SpacingValue | PropertyValueMap["columnGap"]>
  columnRule?: ConditionalValue<CssAny>
  columnRuleColor?: ConditionalValue<PropertyValueMap["columnRuleColor"]>
  columnRuleStyle?: ConditionalValue<WithEscapeHatch<PropertyValueMap["columnRuleStyle"]>>
  columnRuleWidth?: ConditionalValue<PropertyValueMap["columnRuleWidth"]>
  columnSpan?: ConditionalValue<CssAny>
  columnWidth?: ConditionalValue<PropertyValueMap["columnWidth"]>
  columns?: ConditionalValue<CssAny>
  contain?: ConditionalValue<PropertyValueMap["contain"]>
  containIntrinsicBlockSize?: ConditionalValue<PropertyValueMap["containIntrinsicBlockSize"]>
  containIntrinsicHeight?: ConditionalValue<PropertyValueMap["containIntrinsicHeight"]>
  containIntrinsicInlineSize?: ConditionalValue<PropertyValueMap["containIntrinsicInlineSize"]>
  containIntrinsicSize?: ConditionalValue<PropertyValueMap["containIntrinsicSize"]>
  containIntrinsicWidth?: ConditionalValue<PropertyValueMap["containIntrinsicWidth"]>
  container?: ConditionalValue<CssContainer>
  containerName?: ConditionalValue<ContainerName>
  containerType?: ConditionalValue<ContainerTypeValue>
  content?: ConditionalValue<PropertyValueMap["content"]>
  contentVisibility?: ConditionalValue<CssAny>
  counterIncrement?: ConditionalValue<CssAny>
  counterReset?: ConditionalValue<CssAny>
  counterSet?: ConditionalValue<CssAny>
  cursor?: ConditionalValue<CursorValue>
  cx?: ConditionalValue<CssAny>
  cy?: ConditionalValue<CssAny>
  d?: ConditionalValue<CssAny>
  direction?: ConditionalValue<CssAny>
  display?: ConditionalValue<DisplayValue>
  dominantBaseline?: ConditionalValue<PropertyValueMap["dominantBaseline"] | AnyString>
  emptyCells?: ConditionalValue<CssAny>
  fieldSizing?: ConditionalValue<PropertyValueMap["fieldSizing"] | AnyString>
  fill?: ConditionalValue<ColorsValue | CssAny>
  fillOpacity?: ConditionalValue<CssAny>
  fillRule?: ConditionalValue<CssAny>
  filter?: ConditionalValue<FilterValue | CssAny>
  flex?: ConditionalValue<FlexValue | CssAny>
  flexBasis?: ConditionalValue<FlexBasisValue>
  flexDirection?: ConditionalValue<FlexDirectionValue | CssAny>
  flexFlow?: ConditionalValue<PropertyValueMap["flexFlow"]>
  flexGrow?: ConditionalValue<FlexGrowValue>
  flexShrink?: ConditionalValue<FlexShrinkValue>
  flexWrap?: ConditionalValue<CssAny>
  float?: ConditionalValue<FloatValue>
  floodColor?: ConditionalValue<PropertyValueMap["floodColor"]>
  floodOpacity?: ConditionalValue<CssAny>
  font?: ConditionalValue<PropertyValueMap["font"]>
  fontFamily?: ConditionalValue<FontsValue | PropertyValueMap["fontFamily"]>
  fontFeatureSettings?: ConditionalValue<FontFeatureSettingsValue | CssAny>
  fontKerning?: ConditionalValue<FontKerningValue | CssAny>
  fontLanguageOverride?: ConditionalValue<CssAny>
  fontOpticalSizing?: ConditionalValue<CssAny>
  fontPalette?: ConditionalValue<FontPaletteValue | CssAny>
  fontSize?: ConditionalValue<FontSizesValue | PropertyValueMap["fontSize"]>
  fontSizeAdjust?: ConditionalValue<FontSizeAdjustValue>
  fontSmooth?: ConditionalValue<CssAny>
  fontStretch?: ConditionalValue<PropertyValueMap["fontStretch"]>
  fontStyle?: ConditionalValue<PropertyValueMap["fontStyle"]>
  fontSynthesis?: ConditionalValue<PropertyValueMap["fontSynthesis"]>
  fontSynthesisPosition?: ConditionalValue<CssAny>
  fontSynthesisSmallCaps?: ConditionalValue<CssAny>
  fontSynthesisStyle?: ConditionalValue<CssAny>
  fontSynthesisWeight?: ConditionalValue<CssAny>
  fontVariant?: ConditionalValue<FontVariantValue | CssAny>
  fontVariantAlternates?: ConditionalValue<FontVariantAlternatesValue | CssAny>
  fontVariantCaps?: ConditionalValue<FontVariantCapsValue>
  fontVariantEastAsian?: ConditionalValue<PropertyValueMap["fontVariantEastAsian"]>
  fontVariantEmoji?: ConditionalValue<CssAny>
  fontVariantLigatures?: ConditionalValue<PropertyValueMap["fontVariantLigatures"]>
  fontVariantNumeric?: ConditionalValue<FontVariantNumericValue>
  fontVariantPosition?: ConditionalValue<CssAny>
  fontVariationSettings?: ConditionalValue<FontVariationSettingsValue | CssAny>
  fontWeight?: ConditionalValue<FontWeightsValue | PropertyValueMap["fontWeight"]>
  fontWidth?: ConditionalValue<PropertyValueMap["fontWidth"]>
  forcedColorAdjust?: ConditionalValue<CssAny>
  gap?: ConditionalValue<SpacingValue | PropertyValueMap["gap"]>
  glyphOrientationVertical?: ConditionalValue<CssAny>
  grid?: ConditionalValue<CssAny>
  gridArea?: ConditionalValue<CssAny>
  gridAutoColumns?: ConditionalValue<GridAutoColumnsValue | CssAny>
  gridAutoFlow?: ConditionalValue<GridAutoFlowValue>
  gridAutoRows?: ConditionalValue<GridAutoRowsValue | CssAny>
  gridColumn?: ConditionalValue<GridColumnValue | CssAny>
  gridColumnEnd?: ConditionalValue<GridColumnEndValue | CssAny>
  gridColumnGap?: ConditionalValue<SpacingValue | PropertyValueMap["gridColumnGap"]>
  gridColumnStart?: ConditionalValue<GridColumnStartValue | CssAny>
  gridGap?: ConditionalValue<SpacingValue | PropertyValueMap["gridGap"]>
  gridRow?: ConditionalValue<GridRowValue | CssAny>
  gridRowEnd?: ConditionalValue<CssAny>
  gridRowGap?: ConditionalValue<SpacingValue | PropertyValueMap["gridRowGap"]>
  gridRowStart?: ConditionalValue<CssAny>
  gridTemplate?: ConditionalValue<CssAny>
  gridTemplateAreas?: ConditionalValue<CssAny>
  gridTemplateColumns?: ConditionalValue<GridTemplateColumnsValue | CssAny>
  gridTemplateRows?: ConditionalValue<GridTemplateRowsValue | CssAny>
  hangingPunctuation?: ConditionalValue<PropertyValueMap["hangingPunctuation"]>
  height?: ConditionalValue<HeightValue | PropertyValueMap["height"]>
  hyphenateCharacter?: ConditionalValue<CssAny>
  hyphenateLimitChars?: ConditionalValue<CssAny>
  hyphens?: ConditionalValue<HyphensValue>
  imageOrientation?: ConditionalValue<CssAny>
  imageRendering?: ConditionalValue<PropertyValueMap["imageRendering"] | AnyString>
  imageResolution?: ConditionalValue<CssAny>
  imeMode?: ConditionalValue<PropertyValueMap["imeMode"] | AnyString>
  initialLetter?: ConditionalValue<CssAny>
  initialLetterAlign?: ConditionalValue<CssAny>
  inlineSize?: ConditionalValue<InlineSizeValue | PropertyValueMap["inlineSize"]>
  inset?: ConditionalValue<SpacingValue | PropertyValueMap["inset"]>
  insetBlock?: ConditionalValue<SpacingValue | CssAny>
  insetBlockEnd?: ConditionalValue<SpacingValue | CssAny>
  insetBlockStart?: ConditionalValue<SpacingValue | CssAny>
  insetInline?: ConditionalValue<SpacingValue | CssAny>
  insetInlineEnd?: ConditionalValue<SpacingValue | CssAny>
  insetInlineStart?: ConditionalValue<SpacingValue | CssAny>
  interpolateSize?: ConditionalValue<CssAny>
  isolation?: ConditionalValue<CssAny>
  justifyContent?: ConditionalValue<JustifyContentValue>
  justifyItems?: ConditionalValue<PropertyValueMap["justifyItems"]>
  justifySelf?: ConditionalValue<PropertyValueMap["justifySelf"]>
  justifyTracks?: ConditionalValue<CssAny>
  left?: ConditionalValue<SpacingValue | PropertyValueMap["left"]>
  letterSpacing?: ConditionalValue<LetterSpacingsValue | PropertyValueMap["letterSpacing"]>
  lightingColor?: ConditionalValue<PropertyValueMap["lightingColor"]>
  lineBreak?: ConditionalValue<WithEscapeHatch<PropertyValueMap["lineBreak"]>>
  lineClamp?: ConditionalValue<LineClampValue | CssAny>
  lineHeight?: ConditionalValue<LineHeightsValue | PropertyValueMap["lineHeight"]>
  lineHeightStep?: ConditionalValue<PropertyValueMap["lineHeightStep"]>
  listStyle?: ConditionalValue<ListStyleValue | CssAny>
  listStyleImage?: ConditionalValue<AssetsValue>
  listStylePosition?: ConditionalValue<ListStylePositionValue>
  listStyleType?: ConditionalValue<ListStyleTypeValue>
  margin?: ConditionalValue<SpacingValue | PropertyValueMap["margin"]>
  marginBlock?: ConditionalValue<SpacingValue | CssAny>
  marginBlockEnd?: ConditionalValue<SpacingValue | CssAny>
  marginBlockStart?: ConditionalValue<SpacingValue | CssAny>
  marginBottom?: ConditionalValue<SpacingValue | PropertyValueMap["marginBottom"]>
  marginInline?: ConditionalValue<SpacingValue | CssAny>
  marginInlineEnd?: ConditionalValue<SpacingValue | CssAny>
  marginInlineStart?: ConditionalValue<SpacingValue | CssAny>
  marginLeft?: ConditionalValue<SpacingValue | PropertyValueMap["marginLeft"]>
  marginRight?: ConditionalValue<SpacingValue | PropertyValueMap["marginRight"]>
  marginTop?: ConditionalValue<SpacingValue | PropertyValueMap["marginTop"]>
  marginTrim?: ConditionalValue<CssAny>
  marker?: ConditionalValue<CssAny>
  markerEnd?: ConditionalValue<CssAny>
  markerMid?: ConditionalValue<CssAny>
  markerStart?: ConditionalValue<CssAny>
  mask?: ConditionalValue<MaskValue | CssAny>
  maskBorder?: ConditionalValue<PropertyValueMap["maskBorder"]>
  maskBorderMode?: ConditionalValue<CssAny>
  maskBorderOutset?: ConditionalValue<CssAny>
  maskBorderRepeat?: ConditionalValue<CssAny>
  maskBorderSlice?: ConditionalValue<CssAny>
  maskBorderSource?: ConditionalValue<CssAny>
  maskBorderWidth?: ConditionalValue<PropertyValueMap["maskBorderWidth"]>
  maskClip?: ConditionalValue<MaskClipValue | CssAny>
  maskComposite?: ConditionalValue<MaskCompositeValue | CssAny>
  maskImage?: ConditionalValue<MaskImageValue | CssAny>
  maskMode?: ConditionalValue<MaskModeValue | CssAny>
  maskOrigin?: ConditionalValue<MaskOriginValue | CssAny>
  maskPosition?: ConditionalValue<MaskPositionValue>
  maskRepeat?: ConditionalValue<MaskRepeatValue>
  maskSize?: ConditionalValue<MaskSizeValue>
  maskType?: ConditionalValue<MaskTypeValue | CssAny>
  masonryAutoFlow?: ConditionalValue<CssAny>
  mathDepth?: ConditionalValue<CssAny>
  mathShift?: ConditionalValue<CssAny>
  mathStyle?: ConditionalValue<CssAny>
  maxBlockSize?: ConditionalValue<MaxBlockSizeValue | PropertyValueMap["maxBlockSize"]>
  maxHeight?: ConditionalValue<MaxHeightValue | PropertyValueMap["maxHeight"]>
  maxInlineSize?: ConditionalValue<MaxInlineSizeValue | PropertyValueMap["maxInlineSize"]>
  maxLines?: ConditionalValue<CssAny>
  maxWidth?: ConditionalValue<MaxWidthValue | PropertyValueMap["maxWidth"]>
  minBlockSize?: ConditionalValue<MinBlockSizeValue | PropertyValueMap["minBlockSize"]>
  minHeight?: ConditionalValue<MinHeightValue | PropertyValueMap["minHeight"]>
  minInlineSize?: ConditionalValue<MinInlineSizeValue | PropertyValueMap["minInlineSize"]>
  minWidth?: ConditionalValue<MinWidthValue | PropertyValueMap["minWidth"]>
  mixBlendMode?: ConditionalValue<MixBlendModeValue>
  objectFit?: ConditionalValue<ObjectFitValue>
  objectPosition?: ConditionalValue<ObjectPositionValue>
  objectViewBox?: ConditionalValue<CssAny>
  offset?: ConditionalValue<CssAny>
  offsetAnchor?: ConditionalValue<PropertyValueMap["offsetAnchor"]>
  offsetDistance?: ConditionalValue<CssAny>
  offsetPath?: ConditionalValue<CssAny>
  offsetPosition?: ConditionalValue<PropertyValueMap["offsetPosition"]>
  offsetRotate?: ConditionalValue<CssAny>
  opacity?: ConditionalValue<OpacityValue>
  order?: ConditionalValue<PropertyValueMap["order"]>
  orphans?: ConditionalValue<PropertyValueMap["orphans"]>
  outline?: ConditionalValue<BordersValue | CssAny>
  outlineColor?: ConditionalValue<ColorsValue | PropertyValueMap["outlineColor"]>
  outlineOffset?: ConditionalValue<SpacingValue | PropertyValueMap["outlineOffset"]>
  outlineStyle?: ConditionalValue<WithEscapeHatch<PropertyValueMap["outlineStyle"]>>
  outlineWidth?: ConditionalValue<BorderWidthsValue>
  overflow?: ConditionalValue<OverflowValue>
  overflowAnchor?: ConditionalValue<OverflowAnchorValue | CssAny>
  overflowBlock?: ConditionalValue<OverflowBlockValue>
  overflowClipBox?: ConditionalValue<OverflowClipBoxValue | CssAny>
  overflowClipMargin?: ConditionalValue<OverflowClipMarginValue>
  overflowInline?: ConditionalValue<OverflowInlineValue>
  overflowWrap?: ConditionalValue<OverflowWrapValue | CssAny>
  overflowX?: ConditionalValue<OverflowXValue>
  overflowY?: ConditionalValue<OverflowYValue>
  overlay?: ConditionalValue<CssAny>
  overscrollBehavior?: ConditionalValue<OverscrollBehaviorValue>
  overscrollBehaviorBlock?: ConditionalValue<OverscrollBehaviorBlockValue | CssAny>
  overscrollBehaviorInline?: ConditionalValue<OverscrollBehaviorInlineValue | CssAny>
  overscrollBehaviorX?: ConditionalValue<OverscrollBehaviorXValue>
  overscrollBehaviorY?: ConditionalValue<OverscrollBehaviorYValue>
  padding?: ConditionalValue<SpacingValue | PropertyValueMap["padding"]>
  paddingBlock?: ConditionalValue<SpacingValue | CssAny>
  paddingBlockEnd?: ConditionalValue<SpacingValue | CssAny>
  paddingBlockStart?: ConditionalValue<SpacingValue | CssAny>
  paddingBottom?: ConditionalValue<SpacingValue | PropertyValueMap["paddingBottom"]>
  paddingInline?: ConditionalValue<SpacingValue | CssAny>
  paddingInlineEnd?: ConditionalValue<SpacingValue | CssAny>
  paddingInlineStart?: ConditionalValue<SpacingValue | CssAny>
  paddingLeft?: ConditionalValue<SpacingValue | PropertyValueMap["paddingLeft"]>
  paddingRight?: ConditionalValue<SpacingValue | PropertyValueMap["paddingRight"]>
  paddingTop?: ConditionalValue<SpacingValue | PropertyValueMap["paddingTop"]>
  page?: ConditionalValue<CssAny>
  pageBreakAfter?: ConditionalValue<PropertyValueMap["pageBreakAfter"] | AnyString>
  pageBreakBefore?: ConditionalValue<PropertyValueMap["pageBreakBefore"] | AnyString>
  pageBreakInside?: ConditionalValue<CssAny>
  paintOrder?: ConditionalValue<CssAny>
  perspective?: ConditionalValue<PropertyValueMap["perspective"]>
  perspectiveOrigin?: ConditionalValue<PropertyValueMap["perspectiveOrigin"]>
  placeContent?: ConditionalValue<PropertyValueMap["placeContent"]>
  placeItems?: ConditionalValue<PropertyValueMap["placeItems"]>
  placeSelf?: ConditionalValue<PropertyValueMap["placeSelf"]>
  pointerEvents?: ConditionalValue<WithEscapeHatch<PropertyValueMap["pointerEvents"]>>
  position?: ConditionalValue<PositionValue>
  positionAnchor?: ConditionalValue<CssAny>
  positionArea?: ConditionalValue<CssAny>
  positionTry?: ConditionalValue<CssAny>
  positionTryFallbacks?: ConditionalValue<CssAny>
  positionTryOrder?: ConditionalValue<PropertyValueMap["positionTryOrder"] | AnyString>
  positionVisibility?: ConditionalValue<CssAny>
  printColorAdjust?: ConditionalValue<CssAny>
  quotes?: ConditionalValue<CssAny>
  r?: ConditionalValue<CssAny>
  resize?: ConditionalValue<WithEscapeHatch<PropertyValueMap["resize"]>>
  right?: ConditionalValue<SpacingValue | PropertyValueMap["right"]>
  rotate?: ConditionalValue<RotateValue | CssAny>
  rowGap?: ConditionalValue<SpacingValue | PropertyValueMap["rowGap"]>
  rubyAlign?: ConditionalValue<CssAny>
  rubyMerge?: ConditionalValue<CssAny>
  rubyOverhang?: ConditionalValue<CssAny>
  rubyPosition?: ConditionalValue<CssAny>
  rx?: ConditionalValue<CssAny>
  ry?: ConditionalValue<CssAny>
  scale?: ConditionalValue<ScaleValue | CssAny>
  scrollBehavior?: ConditionalValue<ScrollBehaviorValue | CssAny>
  scrollInitialTarget?: ConditionalValue<CssAny>
  scrollMargin?: ConditionalValue<SpacingValue | CssAny>
  scrollMarginBlock?: ConditionalValue<SpacingValue | PropertyValueMap["scrollMarginBlock"]>
  scrollMarginBlockEnd?: ConditionalValue<SpacingValue | PropertyValueMap["scrollMarginBlockEnd"]>
  scrollMarginBlockStart?: ConditionalValue<SpacingValue | PropertyValueMap["scrollMarginBlockStart"]>
  scrollMarginBottom?: ConditionalValue<SpacingValue | PropertyValueMap["scrollMarginBottom"]>
  scrollMarginInline?: ConditionalValue<SpacingValue | PropertyValueMap["scrollMarginInline"]>
  scrollMarginInlineEnd?: ConditionalValue<SpacingValue | PropertyValueMap["scrollMarginInlineEnd"]>
  scrollMarginInlineStart?: ConditionalValue<SpacingValue | PropertyValueMap["scrollMarginInlineStart"]>
  scrollMarginLeft?: ConditionalValue<SpacingValue | PropertyValueMap["scrollMarginLeft"]>
  scrollMarginRight?: ConditionalValue<SpacingValue | PropertyValueMap["scrollMarginRight"]>
  scrollMarginTop?: ConditionalValue<SpacingValue | PropertyValueMap["scrollMarginTop"]>
  scrollPadding?: ConditionalValue<SpacingValue | CssAny>
  scrollPaddingBlock?: ConditionalValue<SpacingValue | PropertyValueMap["scrollPaddingBlock"]>
  scrollPaddingBlockEnd?: ConditionalValue<SpacingValue | PropertyValueMap["scrollPaddingBlockEnd"]>
  scrollPaddingBlockStart?: ConditionalValue<SpacingValue | PropertyValueMap["scrollPaddingBlockStart"]>
  scrollPaddingBottom?: ConditionalValue<SpacingValue | PropertyValueMap["scrollPaddingBottom"]>
  scrollPaddingInline?: ConditionalValue<SpacingValue | PropertyValueMap["scrollPaddingInline"]>
  scrollPaddingInlineEnd?: ConditionalValue<SpacingValue | PropertyValueMap["scrollPaddingInlineEnd"]>
  scrollPaddingInlineStart?: ConditionalValue<SpacingValue | PropertyValueMap["scrollPaddingInlineStart"]>
  scrollPaddingLeft?: ConditionalValue<SpacingValue | PropertyValueMap["scrollPaddingLeft"]>
  scrollPaddingRight?: ConditionalValue<SpacingValue | PropertyValueMap["scrollPaddingRight"]>
  scrollPaddingTop?: ConditionalValue<SpacingValue | PropertyValueMap["scrollPaddingTop"]>
  scrollSnapAlign?: ConditionalValue<ScrollSnapAlignValue>
  scrollSnapCoordinate?: ConditionalValue<ScrollSnapCoordinateValue | CssAny>
  scrollSnapDestination?: ConditionalValue<ScrollSnapDestinationValue>
  scrollSnapPointsX?: ConditionalValue<ScrollSnapPointsXValue | CssAny>
  scrollSnapPointsY?: ConditionalValue<ScrollSnapPointsYValue | CssAny>
  scrollSnapStop?: ConditionalValue<ScrollSnapStopValue>
  scrollSnapType?: ConditionalValue<ScrollSnapTypeValue>
  scrollSnapTypeX?: ConditionalValue<ScrollSnapTypeXValue | CssAny>
  scrollSnapTypeY?: ConditionalValue<ScrollSnapTypeYValue | CssAny>
  scrollTimeline?: ConditionalValue<ScrollTimelineValue | CssAny>
  scrollTimelineAxis?: ConditionalValue<ScrollTimelineAxisValue | CssAny>
  scrollTimelineName?: ConditionalValue<ScrollTimelineNameValue | CssAny>
  scrollbarColor?: ConditionalValue<ScrollbarColorValue | CssAny>
  scrollbarGutter?: ConditionalValue<ScrollbarGutterValue>
  scrollbarWidth?: ConditionalValue<ScrollbarWidthValue>
  shapeImageThreshold?: ConditionalValue<CssAny>
  shapeMargin?: ConditionalValue<PropertyValueMap["shapeMargin"]>
  shapeOutside?: ConditionalValue<PropertyValueMap["shapeOutside"]>
  shapeRendering?: ConditionalValue<CssAny>
  speakAs?: ConditionalValue<PropertyValueMap["speakAs"]>
  stopColor?: ConditionalValue<PropertyValueMap["stopColor"]>
  stopOpacity?: ConditionalValue<CssAny>
  stroke?: ConditionalValue<ColorsValue | CssAny>
  strokeColor?: ConditionalValue<PropertyValueMap["strokeColor"]>
  strokeDasharray?: ConditionalValue<StrokeDasharrayValue | CssAny>
  strokeDashoffset?: ConditionalValue<StrokeDashoffsetValue | CssAny>
  strokeLinecap?: ConditionalValue<StrokeLinecapValue | CssAny>
  strokeLinejoin?: ConditionalValue<StrokeLinejoinValue>
  strokeMiterlimit?: ConditionalValue<StrokeMiterlimitValue | CssAny>
  strokeOpacity?: ConditionalValue<StrokeOpacityValue | CssAny>
  strokeWidth?: ConditionalValue<BorderWidthsValue>
  tabSize?: ConditionalValue<PropertyValueMap["tabSize"]>
  tableLayout?: ConditionalValue<TableLayoutValue>
  textAlign?: ConditionalValue<TextAlignValue>
  textAlignLast?: ConditionalValue<PropertyValueMap["textAlignLast"] | AnyString>
  textAnchor?: ConditionalValue<CssAny>
  textAutospace?: ConditionalValue<PropertyValueMap["textAutospace"]>
  textBox?: ConditionalValue<PropertyValueMap["textBox"]>
  textBoxEdge?: ConditionalValue<PropertyValueMap["textBoxEdge"]>
  textBoxTrim?: ConditionalValue<CssAny>
  textCombineUpright?: ConditionalValue<CssAny>
  textDecoration?: ConditionalValue<TextDecorationValue | CssAny>
  textDecorationColor?: ConditionalValue<ColorsValue | PropertyValueMap["textDecorationColor"]>
  textDecorationLine?: ConditionalValue<PropertyValueMap["textDecorationLine"]>
  textDecorationSkip?: ConditionalValue<PropertyValueMap["textDecorationSkip"]>
  textDecorationSkipInk?: ConditionalValue<PropertyValueMap["textDecorationSkipInk"] | AnyString>
  textDecorationStyle?: ConditionalValue<TextDecorationStyleValue>
  textDecorationThickness?: ConditionalValue<TextDecorationThicknessValue>
  textEmphasis?: ConditionalValue<PropertyValueMap["textEmphasis"]>
  textEmphasisColor?: ConditionalValue<ColorsValue | PropertyValueMap["textEmphasisColor"]>
  textEmphasisPosition?: ConditionalValue<PropertyValueMap["textEmphasisPosition"]>
  textEmphasisStyle?: ConditionalValue<PropertyValueMap["textEmphasisStyle"]>
  textIndent?: ConditionalValue<SpacingValue | PropertyValueMap["textIndent"]>
  textJustify?: ConditionalValue<PropertyValueMap["textJustify"] | AnyString>
  textOrientation?: ConditionalValue<PropertyValueMap["textOrientation"] | AnyString>
  textOverflow?: ConditionalValue<TextOverflowValue>
  textRendering?: ConditionalValue<PropertyValueMap["textRendering"] | AnyString>
  textShadow?: ConditionalValue<ShadowsValue | CssAny>
  textSizeAdjust?: ConditionalValue<TextSizeAdjustValue>
  textSpacingTrim?: ConditionalValue<CssAny>
  textTransform?: ConditionalValue<TextTransformValue>
  textUnderlineOffset?: ConditionalValue<TextUnderlineOffsetValue>
  textUnderlinePosition?: ConditionalValue<PropertyValueMap["textUnderlinePosition"]>
  textWrap?: ConditionalValue<TextWrapValue>
  textWrapMode?: ConditionalValue<PropertyValueMap["textWrapMode"] | AnyString>
  textWrapStyle?: ConditionalValue<PropertyValueMap["textWrapStyle"] | AnyString>
  timelineScope?: ConditionalValue<CssAny>
  top?: ConditionalValue<SpacingValue | PropertyValueMap["top"]>
  touchAction?: ConditionalValue<TouchActionValue>
  transform?: ConditionalValue<TransformValue | CssAny>
  transformBox?: ConditionalValue<TransformBoxValue>
  transformOrigin?: ConditionalValue<TransformOriginValue | CssAny>
  transformStyle?: ConditionalValue<TransformStyleValue | CssAny>
  transition?: ConditionalValue<TransitionValue | CssAny>
  transitionBehavior?: ConditionalValue<PropertyValueMap["transitionBehavior"]>
  transitionDelay?: ConditionalValue<DurationsValue | CssAny>
  transitionDuration?: ConditionalValue<DurationsValue | CssAny>
  transitionProperty?: ConditionalValue<TransitionPropertyValue>
  transitionTimingFunction?: ConditionalValue<EasingsValue | CssAny>
  translate?: ConditionalValue<TranslateValue | CssAny>
  unicodeBidi?: ConditionalValue<PropertyValueMap["unicodeBidi"] | AnyString>
  userSelect?: ConditionalValue<UserSelectValue>
  vectorEffect?: ConditionalValue<PropertyValueMap["vectorEffect"] | AnyString>
  verticalAlign?: ConditionalValue<VerticalAlignValue>
  viewTimeline?: ConditionalValue<CssAny>
  viewTimelineAxis?: ConditionalValue<CssAny>
  viewTimelineInset?: ConditionalValue<PropertyValueMap["viewTimelineInset"]>
  viewTimelineName?: ConditionalValue<CssAny>
  viewTransitionClass?: ConditionalValue<CssAny>
  viewTransitionName?: ConditionalValue<CssAny>
  visibility?: ConditionalValue<VisibilityValue | CssAny>
  whiteSpace?: ConditionalValue<PropertyValueMap["whiteSpace"]>
  whiteSpaceCollapse?: ConditionalValue<PropertyValueMap["whiteSpaceCollapse"] | AnyString>
  widows?: ConditionalValue<PropertyValueMap["widows"]>
  width?: ConditionalValue<WidthValue | PropertyValueMap["width"]>
  willChange?: ConditionalValue<CssAny>
  wordBreak?: ConditionalValue<WordBreakValue>
  wordSpacing?: ConditionalValue<PropertyValueMap["wordSpacing"]>
  wordWrap?: ConditionalValue<PropertyValueMap["wordWrap"] | AnyString>
  writingMode?: ConditionalValue<WithEscapeHatch<PropertyValueMap["writingMode"]>>
  x?: ConditionalValue<TranslateXValue>
  y?: ConditionalValue<TranslateYValue>
  zIndex?: ConditionalValue<ZIndexValue>
  zoom?: ConditionalValue<PropertyValueMap["zoom"]>
  animationState?: ConditionalValue<AnimationStateValue | CssAny>
  backdropBlur?: ConditionalValue<BlursValue>
  backdropBrightness?: ConditionalValue<BackdropBrightnessValue | CssAny>
  backdropContrast?: ConditionalValue<BackdropContrastValue | CssAny>
  backdropGrayscale?: ConditionalValue<BackdropGrayscaleValue | CssAny>
  backdropHueRotate?: ConditionalValue<BackdropHueRotateValue | CssAny>
  backdropInvert?: ConditionalValue<BackdropInvertValue | CssAny>
  backdropOpacity?: ConditionalValue<BackdropOpacityValue | CssAny>
  backdropSaturate?: ConditionalValue<BackdropSaturateValue | CssAny>
  backdropSepia?: ConditionalValue<BackdropSepiaValue | CssAny>
  backgroundConic?: ConditionalValue<BackgroundConicValue | CssAny>
  backgroundGradient?: ConditionalValue<BackgroundGradientValue | CssAny>
  backgroundLinear?: ConditionalValue<BackgroundLinearValue | CssAny>
  backgroundRadial?: ConditionalValue<GradientsValue | CssAny>
  bg?: ConditionalValue<ColorsValue | CssAny>
  bgAttachment?: ConditionalValue<BackgroundAttachmentValue | CssAny>
  bgBlendMode?: ConditionalValue<BackgroundBlendModeValue | CssAny>
  bgClip?: ConditionalValue<BackgroundClipValue | CssAny>
  bgColor?: ConditionalValue<ColorsValue | PropertyValueMap["backgroundColor"]>
  bgConic?: ConditionalValue<BackgroundConicValue | CssAny>
  bgGradient?: ConditionalValue<BackgroundGradientValue | CssAny>
  bgImage?: ConditionalValue<AssetsValue | CssAny>
  bgLinear?: ConditionalValue<BackgroundLinearValue | CssAny>
  bgOrigin?: ConditionalValue<BackgroundOriginValue | CssAny>
  bgPosition?: ConditionalValue<BackgroundPositionValue>
  bgPositionX?: ConditionalValue<BackgroundPositionXValue | CssAny>
  bgPositionY?: ConditionalValue<BackgroundPositionYValue | CssAny>
  bgRadial?: ConditionalValue<GradientsValue | CssAny>
  bgRepeat?: ConditionalValue<BackgroundRepeatValue>
  bgSize?: ConditionalValue<BackgroundSizeValue>
  blur?: ConditionalValue<BlursValue>
  borderBottomRadius?: ConditionalValue<RadiiValue | PropertyValueMap["borderRadius"]>
  borderEnd?: ConditionalValue<BordersValue | CssAny>
  borderEndColor?: ConditionalValue<ColorsValue | PropertyValueMap["borderInlineEndColor"]>
  borderEndRadius?: ConditionalValue<RadiiValue | PropertyValueMap["borderRadius"]>
  borderEndWidth?: ConditionalValue<BorderWidthsValue>
  borderLeftRadius?: ConditionalValue<RadiiValue | PropertyValueMap["borderRadius"]>
  borderRightRadius?: ConditionalValue<RadiiValue | PropertyValueMap["borderRadius"]>
  borderSpacingX?: ConditionalValue<SpacingValue>
  borderSpacingY?: ConditionalValue<SpacingValue>
  borderStart?: ConditionalValue<BordersValue | CssAny>
  borderStartColor?: ConditionalValue<ColorsValue | PropertyValueMap["borderInlineStartColor"]>
  borderStartRadius?: ConditionalValue<RadiiValue | PropertyValueMap["borderRadius"]>
  borderStartWidth?: ConditionalValue<BorderWidthsValue>
  borderTopRadius?: ConditionalValue<RadiiValue | PropertyValueMap["borderRadius"]>
  borderX?: ConditionalValue<BordersValue | CssAny>
  borderXColor?: ConditionalValue<ColorsValue | PropertyValueMap["borderInlineColor"]>
  borderXWidth?: ConditionalValue<BorderWidthsValue>
  borderY?: ConditionalValue<BordersValue | CssAny>
  borderYColor?: ConditionalValue<ColorsValue | PropertyValueMap["borderBlockColor"]>
  borderYWidth?: ConditionalValue<BorderWidthsValue>
  boxShadowColor?: ConditionalValue<ColorsValue>
  boxSize?: ConditionalValue<BoxSizeValue>
  brightness?: ConditionalValue<BrightnessValue | CssAny>
  colorPalette?: ConditionalValue<ColorPaletteValue | CssAny>
  contrast?: ConditionalValue<ContrastValue | CssAny>
  debug?: ConditionalValue<DebugValue | CssAny>
  divideColor?: ConditionalValue<ColorsValue>
  divideStyle?: ConditionalValue<BorderStyleValue>
  divideX?: ConditionalValue<BorderWidthsValue | CssAny>
  divideY?: ConditionalValue<BorderWidthsValue | CssAny>
  dropShadow?: ConditionalValue<DropShadowsValue>
  end?: ConditionalValue<SpacingValue | CssAny>
  flexDir?: ConditionalValue<FlexDirectionValue | CssAny>
  focusRing?: ConditionalValue<FocusRingValue | CssAny>
  focusRingColor?: ConditionalValue<ColorsValue>
  focusRingOffset?: ConditionalValue<SpacingValue>
  focusRingStyle?: ConditionalValue<BorderStylesValue | PropertyValueMap["outlineStyle"]>
  focusRingWidth?: ConditionalValue<BorderWidthsValue>
  focusVisibleRing?: ConditionalValue<FocusVisibleRingValue | CssAny>
  fontSmoothing?: ConditionalValue<FontSmoothingValue | CssAny>
  gradientFrom?: ConditionalValue<ColorsValue>
  gradientFromPosition?: ConditionalValue<GradientFromPositionValue | CssAny>
  gradientTo?: ConditionalValue<ColorsValue>
  gradientToPosition?: ConditionalValue<GradientToPositionValue | CssAny>
  gradientVia?: ConditionalValue<ColorsValue>
  gradientViaPosition?: ConditionalValue<GradientViaPositionValue | CssAny>
  grayscale?: ConditionalValue<GrayscaleValue | CssAny>
  h?: ConditionalValue<HeightValue | PropertyValueMap["height"]>
  hideBelow?: ConditionalValue<BreakpointsValue>
  hideFrom?: ConditionalValue<BreakpointsValue>
  hueRotate?: ConditionalValue<HueRotateValue | CssAny>
  insetEnd?: ConditionalValue<SpacingValue | CssAny>
  insetStart?: ConditionalValue<SpacingValue | CssAny>
  insetX?: ConditionalValue<SpacingValue | CssAny>
  insetY?: ConditionalValue<SpacingValue | CssAny>
  invert?: ConditionalValue<InvertValue | CssAny>
  m?: ConditionalValue<SpacingValue | PropertyValueMap["margin"]>
  marginEnd?: ConditionalValue<SpacingValue | CssAny>
  marginStart?: ConditionalValue<SpacingValue | CssAny>
  marginX?: ConditionalValue<SpacingValue | CssAny>
  marginY?: ConditionalValue<SpacingValue | CssAny>
  maskBottomFrom?: ConditionalValue<SpacingValue>
  maskBottomFromColor?: ConditionalValue<ColorsValue>
  maskBottomTo?: ConditionalValue<SpacingValue>
  maskBottomToColor?: ConditionalValue<ColorsValue>
  maskConic?: ConditionalValue<MaskConicValue | CssAny>
  maskConicFrom?: ConditionalValue<SpacingValue>
  maskConicFromColor?: ConditionalValue<ColorsValue>
  maskConicTo?: ConditionalValue<SpacingValue>
  maskConicToColor?: ConditionalValue<ColorsValue>
  maskLeftFrom?: ConditionalValue<SpacingValue>
  maskLeftFromColor?: ConditionalValue<ColorsValue>
  maskLeftTo?: ConditionalValue<SpacingValue>
  maskLeftToColor?: ConditionalValue<ColorsValue>
  maskLinear?: ConditionalValue<MaskLinearValue | CssAny>
  maskLinearFrom?: ConditionalValue<SpacingValue>
  maskLinearFromColor?: ConditionalValue<ColorsValue>
  maskLinearTo?: ConditionalValue<SpacingValue>
  maskLinearToColor?: ConditionalValue<ColorsValue>
  maskRadial?: ConditionalValue<MaskRadialValue | CssAny>
  maskRadialAt?: ConditionalValue<MaskRadialAtValue | CssAny>
  maskRadialFrom?: ConditionalValue<SpacingValue>
  maskRadialFromColor?: ConditionalValue<ColorsValue>
  maskRadialShape?: ConditionalValue<MaskRadialShapeValue | CssAny>
  maskRadialSize?: ConditionalValue<MaskRadialSizeValue | CssAny>
  maskRadialTo?: ConditionalValue<SpacingValue>
  maskRadialToColor?: ConditionalValue<ColorsValue>
  maskRightFrom?: ConditionalValue<SpacingValue>
  maskRightFromColor?: ConditionalValue<ColorsValue>
  maskRightTo?: ConditionalValue<SpacingValue>
  maskRightToColor?: ConditionalValue<ColorsValue>
  maskTopFrom?: ConditionalValue<SpacingValue>
  maskTopFromColor?: ConditionalValue<ColorsValue>
  maskTopTo?: ConditionalValue<SpacingValue>
  maskTopToColor?: ConditionalValue<ColorsValue>
  maskXFrom?: ConditionalValue<SpacingValue>
  maskXFromColor?: ConditionalValue<ColorsValue>
  maskXTo?: ConditionalValue<SpacingValue>
  maskXToColor?: ConditionalValue<ColorsValue>
  maskYFrom?: ConditionalValue<SpacingValue>
  maskYFromColor?: ConditionalValue<ColorsValue>
  maskYTo?: ConditionalValue<SpacingValue>
  maskYToColor?: ConditionalValue<ColorsValue>
  maxH?: ConditionalValue<MaxHeightValue | PropertyValueMap["maxHeight"]>
  maxW?: ConditionalValue<MaxWidthValue | PropertyValueMap["maxWidth"]>
  mb?: ConditionalValue<SpacingValue | PropertyValueMap["marginBottom"]>
  me?: ConditionalValue<SpacingValue | CssAny>
  minH?: ConditionalValue<MinHeightValue | PropertyValueMap["minHeight"]>
  minW?: ConditionalValue<MinWidthValue | PropertyValueMap["minWidth"]>
  ml?: ConditionalValue<SpacingValue | PropertyValueMap["marginLeft"]>
  mr?: ConditionalValue<SpacingValue | PropertyValueMap["marginRight"]>
  ms?: ConditionalValue<SpacingValue | CssAny>
  mt?: ConditionalValue<SpacingValue | PropertyValueMap["marginTop"]>
  mx?: ConditionalValue<SpacingValue | CssAny>
  my?: ConditionalValue<SpacingValue | CssAny>
  p?: ConditionalValue<SpacingValue | PropertyValueMap["padding"]>
  paddingEnd?: ConditionalValue<SpacingValue | CssAny>
  paddingStart?: ConditionalValue<SpacingValue | CssAny>
  paddingX?: ConditionalValue<SpacingValue | CssAny>
  paddingY?: ConditionalValue<SpacingValue | CssAny>
  pb?: ConditionalValue<SpacingValue | PropertyValueMap["paddingBottom"]>
  pe?: ConditionalValue<SpacingValue | CssAny>
  pl?: ConditionalValue<SpacingValue | PropertyValueMap["paddingLeft"]>
  pos?: ConditionalValue<PositionValue>
  pr?: ConditionalValue<SpacingValue | PropertyValueMap["paddingRight"]>
  ps?: ConditionalValue<SpacingValue | CssAny>
  pt?: ConditionalValue<SpacingValue | PropertyValueMap["paddingTop"]>
  px?: ConditionalValue<SpacingValue | CssAny>
  py?: ConditionalValue<SpacingValue | CssAny>
  ring?: ConditionalValue<BordersValue | CssAny>
  ringColor?: ConditionalValue<ColorsValue | PropertyValueMap["outlineColor"]>
  ringOffset?: ConditionalValue<SpacingValue | PropertyValueMap["outlineOffset"]>
  ringWidth?: ConditionalValue<BorderWidthsValue>
  rotateX?: ConditionalValue<RotateValue | CssAny>
  rotateY?: ConditionalValue<RotateValue | CssAny>
  rotateZ?: ConditionalValue<RotateValue | CssAny>
  rounded?: ConditionalValue<RadiiValue | PropertyValueMap["borderRadius"]>
  roundedBottom?: ConditionalValue<RadiiValue | PropertyValueMap["borderRadius"]>
  roundedBottomLeft?: ConditionalValue<RadiiValue | PropertyValueMap["borderBottomLeftRadius"]>
  roundedBottomRight?: ConditionalValue<RadiiValue | PropertyValueMap["borderBottomRightRadius"]>
  roundedEnd?: ConditionalValue<RadiiValue | PropertyValueMap["borderRadius"]>
  roundedEndEnd?: ConditionalValue<RadiiValue | PropertyValueMap["borderEndEndRadius"]>
  roundedEndStart?: ConditionalValue<RadiiValue | PropertyValueMap["borderEndStartRadius"]>
  roundedLeft?: ConditionalValue<RadiiValue | PropertyValueMap["borderRadius"]>
  roundedRight?: ConditionalValue<RadiiValue | PropertyValueMap["borderRadius"]>
  roundedStart?: ConditionalValue<RadiiValue | PropertyValueMap["borderRadius"]>
  roundedStartEnd?: ConditionalValue<RadiiValue | PropertyValueMap["borderStartEndRadius"]>
  roundedStartStart?: ConditionalValue<RadiiValue | PropertyValueMap["borderStartStartRadius"]>
  roundedTop?: ConditionalValue<RadiiValue | PropertyValueMap["borderRadius"]>
  roundedTopLeft?: ConditionalValue<RadiiValue | PropertyValueMap["borderTopLeftRadius"]>
  roundedTopRight?: ConditionalValue<RadiiValue | PropertyValueMap["borderTopRightRadius"]>
  saturate?: ConditionalValue<SaturateValue | CssAny>
  scaleX?: ConditionalValue<ScaleXValue | CssAny>
  scaleY?: ConditionalValue<ScaleYValue | CssAny>
  scrollMarginX?: ConditionalValue<SpacingValue | PropertyValueMap["scrollMarginInline"]>
  scrollMarginY?: ConditionalValue<SpacingValue | PropertyValueMap["scrollMarginBlock"]>
  scrollPaddingX?: ConditionalValue<SpacingValue | PropertyValueMap["scrollPaddingInline"]>
  scrollPaddingY?: ConditionalValue<SpacingValue | PropertyValueMap["scrollPaddingBlock"]>
  scrollSnapMargin?: ConditionalValue<SpacingValue>
  scrollSnapMarginBottom?: ConditionalValue<SpacingValue>
  scrollSnapMarginLeft?: ConditionalValue<SpacingValue>
  scrollSnapMarginRight?: ConditionalValue<SpacingValue>
  scrollSnapMarginTop?: ConditionalValue<SpacingValue>
  scrollSnapStrictness?: ConditionalValue<ScrollSnapStrictnessValue | CssAny>
  scrollbar?: ConditionalValue<ScrollbarValue | CssAny>
  scrollbarThumb?: ConditionalValue<ColorsValue>
  scrollbarTrack?: ConditionalValue<ColorsValue>
  sepia?: ConditionalValue<SepiaValue | CssAny>
  shadow?: ConditionalValue<ShadowsValue | CssAny>
  shadowColor?: ConditionalValue<ColorsValue>
  spaceX?: ConditionalValue<SpacingValue | CssAny>
  spaceY?: ConditionalValue<SpacingValue | CssAny>
  srOnly?: ConditionalValue<SrOnlyValue | CssAny>
  start?: ConditionalValue<SpacingValue | CssAny>
  textGradient?: ConditionalValue<TextGradientValue | CssAny>
  textShadowColor?: ConditionalValue<ColorsValue>
  textStyle?: ConditionalValue<TextStyleValue | CssAny>
  translateX?: ConditionalValue<TranslateXValue>
  translateY?: ConditionalValue<TranslateYValue>
  translateZ?: ConditionalValue<SpacingValue>
  truncate?: ConditionalValue<TruncateValue | CssAny>
  w?: ConditionalValue<WidthValue | PropertyValueMap["width"]>
  z?: ConditionalValue<SpacingValue>
}

export type CssVarValue = ConditionalValue<CssVars | AnyString | AnyNumber>

export type CssVarProperties = {
  [K in `--${string}`]?: CssVarValue
}

export type NestedStyles = {
  [K in Selector | Condition]?: SystemStyleObject
}

export interface SystemStyleObject extends SystemProperties, CssVarProperties, NestedStyles {}

export type SystemStyleObjectWith<K extends keyof SystemProperties> = Pick<SystemProperties, K> & {
  [S in Selector | Condition]?: SystemStyleObjectWith<K>
}

export type SystemStyleObjectWithout<K extends keyof SystemProperties> = SystemStyleObjectWith<Exclude<keyof SystemProperties, K>>

export interface GlobalStyleObject {
  [selector: string]: SystemStyleObject
}

export interface CssKeyframes {
  [name: string]: {
    [time: string]: SystemStyleObject
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