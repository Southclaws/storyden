import type { Conditions, Breakpoints, ContainerName, SystemProperties } from './system';
import type { Tokens, Token, ColorOpacityModifier, ColorOpacityToken, ColorPalette, AnimationToken, BlurToken, BorderWidthToken, BorderToken, BreakpointToken, ColorToken, DurationToken, EasingToken, FontSizeToken, FontWeightToken, FontToken, LetterSpacingToken, LineHeightToken, OpacityToken, RadiusToken, ShadowToken, SizeToken, SpacingToken, ZIndexToken } from './tokens';

export interface Register {}

export type Registered<K extends PropertyKey, Local> = Register extends { system: infer System }
  ? K extends keyof System
    ? System[K]
    : Local
  : Local

export interface PandaSystem {
  Conditions: Conditions
  Breakpoints: Breakpoints
  ContainerName: ContainerName
  SystemProperties: SystemProperties
  Tokens: Tokens
  Token: Token
  ColorOpacityModifier: ColorOpacityModifier
  ColorOpacityToken: ColorOpacityToken
  ColorPalette: ColorPalette
  AnimationToken: AnimationToken
  BlurToken: BlurToken
  BorderWidthToken: BorderWidthToken
  BorderToken: BorderToken
  BreakpointToken: BreakpointToken
  ColorToken: ColorToken
  DurationToken: DurationToken
  EasingToken: EasingToken
  FontSizeToken: FontSizeToken
  FontWeightToken: FontWeightToken
  FontToken: FontToken
  LetterSpacingToken: LetterSpacingToken
  LineHeightToken: LineHeightToken
  OpacityToken: OpacityToken
  RadiusToken: RadiusToken
  ShadowToken: ShadowToken
  SizeToken: SizeToken
  SpacingToken: SpacingToken
  ZIndexToken: ZIndexToken
}