import { withoutSpace } from '../helpers';

const conditions = /* @__PURE__ */ new Set(/* @__PURE__ */ "2xl,2xlDown,2xlOnly,@/2xl,@/2xlDown,@/2xlOnly,@/2xlTo3xl,@/2xlTo4xl,@/2xlTo5xl,@/2xlTo6xl,@/2xlTo7xl,@/2xlTo8xl,@/2xs,@/2xsDown,@/2xsOnly,@/2xsTo2xl,@/2xsTo3xl,@/2xsTo4xl,@/2xsTo5xl,@/2xsTo6xl,@/2xsTo7xl,@/2xsTo8xl,@/2xsToLg,@/2xsToMd,@/2xsToSm,@/2xsToXl,@/2xsToXs,@/3xl,@/3xlDown,@/3xlOnly,@/3xlTo4xl,@/3xlTo5xl,@/3xlTo6xl,@/3xlTo7xl,@/3xlTo8xl,@/3xs,@/3xsDown,@/3xsOnly,@/3xsTo2xl,@/3xsTo2xs,@/3xsTo3xl,@/3xsTo4xl,@/3xsTo5xl,@/3xsTo6xl,@/3xsTo7xl,@/3xsTo8xl,@/3xsToLg,@/3xsToMd,@/3xsToSm,@/3xsToXl,@/3xsToXs,@/4xl,@/4xlDown,@/4xlOnly,@/4xlTo5xl,@/4xlTo6xl,@/4xlTo7xl,@/4xlTo8xl,@/5xl,@/5xlDown,@/5xlOnly,@/5xlTo6xl,@/5xlTo7xl,@/5xlTo8xl,@/6xl,@/6xlDown,@/6xlOnly,@/6xlTo7xl,@/6xlTo8xl,@/7xl,@/7xlDown,@/7xlOnly,@/7xlTo8xl,@/8xl,@/8xlDown,@/8xlOnly,@/lg,@/lgDown,@/lgOnly,@/lgTo2xl,@/lgTo3xl,@/lgTo4xl,@/lgTo5xl,@/lgTo6xl,@/lgTo7xl,@/lgTo8xl,@/lgToXl,@/md,@/mdDown,@/mdOnly,@/mdTo2xl,@/mdTo3xl,@/mdTo4xl,@/mdTo5xl,@/mdTo6xl,@/mdTo7xl,@/mdTo8xl,@/mdToLg,@/mdToXl,@/sm,@/smDown,@/smOnly,@/smTo2xl,@/smTo3xl,@/smTo4xl,@/smTo5xl,@/smTo6xl,@/smTo7xl,@/smTo8xl,@/smToLg,@/smToMd,@/smToXl,@/xl,@/xlDown,@/xlOnly,@/xlTo2xl,@/xlTo3xl,@/xlTo4xl,@/xlTo5xl,@/xlTo6xl,@/xlTo7xl,@/xlTo8xl,@/xs,@/xsDown,@/xsOnly,@/xsTo2xl,@/xsTo3xl,@/xsTo4xl,@/xsTo5xl,@/xsTo6xl,@/xsTo7xl,@/xsTo8xl,@/xsToLg,@/xsToMd,@/xsToSm,@/xsToXl,_active,_after,_anyPointerCoarse,_anyPointerFine,_anyPointerNone,_atValue,_autofill,_backdrop,_before,_checked,_closed,_complete,_current,_currentPage,_currentStep,_dark,_default,_disabled,_dragging,_empty,_enabled,_even,_expanded,_file,_first,_firstLetter,_firstLine,_firstOfType,_focus,_focusVisible,_focusWithin,_fullscreen,_grabbed,_groupActive,_groupChecked,_groupDisabled,_groupExpanded,_groupFocus,_groupFocusVisible,_groupFocusWithin,_groupHover,_groupInvalid,_hidden,_highContrast,_highlighted,_horizontal,_hover,_icon,_inRange,_incomplete,_indeterminate,_inert,_invalid,_invertedColors,_landscape,_last,_lastOfType,_lessContrast,_light,_loading,_ltr,_marker,_moreContrast,_motionReduce,_motionSafe,_noscript,_now,_odd,_only,_onlyOfType,_open,_optional,_osDark,_osLight,_outOfRange,_overValue,_peerActive,_peerChecked,_peerDisabled,_peerExpanded,_peerFocus,_peerFocusVisible,_peerFocusWithin,_peerHover,_peerInvalid,_peerPlaceholderShown,_placeholder,_placeholderShown,_pointerCoarse,_pointerFine,_pointerNone,_portrait,_pressed,_print,_rangeEnd,_rangeStart,_readOnly,_readWrite,_required,_rtl,_scrollbar,_scrollbarThumb,_scrollbarTrack,_selected,_selection,_starting,_target,_today,_topmost,_unavailable,_underValue,_userInvalid,_userValid,_valid,_vertical,_visited,base,lg,lgDown,lgOnly,lgTo2xl,lgToXl,md,mdDown,mdOnly,mdTo2xl,mdToLg,mdToXl,sm,smDown,smOnly,smTo2xl,smToLg,smToMd,smToXl,xl,xlDown,xlOnly,xlTo2xl".split(','))
const conditionRe = /^@|&/
const underscoreRe = /^_/
const selectorRe = /&|@/

export const breakpointKeys = ["base","sm","md","lg","xl","2xl"]

export function isCondition(v) {
  return conditions.has(v) || conditionRe.test(v)
}

export function finalizeConditions(paths) {
  return paths.map((p) => {
    if (conditions.has(p)) {
      return p.replace(underscoreRe, '')
    }
    if (selectorRe.test(p)) {
      return `[${withoutSpace(p.trim())}]`
    }
    return p
  })
}

export function sortConditions(paths) {
  return [...paths].sort((a, b) => {
    const aa = isCondition(a)
    const bb = isCondition(b)
    return aa && !bb ? 1 : !aa && bb ? -1 : 0
  })
}