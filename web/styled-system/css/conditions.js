import { withoutSpace } from '../helpers';

const conditions = new Set("2xl,2xlDown,2xlOnly,_active,_after,_atValue,_autofill,_backdrop,_before,_checked,_closed,_collapsed,_complete,_containerLarge,_containerMedium,_containerSmall,_current,_currentPage,_currentStep,_dark,_default,_disabled,_dragging,_empty,_enabled,_even,_expanded,_file,_first,_firstLetter,_firstLine,_firstOfType,_focus,_focusVisible,_focusWithin,_fullscreen,_grabbed,_groupActive,_groupChecked,_groupDisabled,_groupExpanded,_groupFocus,_groupFocusVisible,_groupFocusWithin,_groupHover,_groupInvalid,_hidden,_highContrast,_highlighted,_horizontal,_hover,_icon,_inRange,_incomplete,_indeterminate,_invalid,_invertedColors,_landscape,_last,_lastOfType,_lessContrast,_light,_loading,_ltr,_marker,_moreContrast,_motionReduce,_motionSafe,_noscript,_now,_odd,_off,_on,_only,_onlyOfType,_open,_optional,_osDark,_osLight,_outOfRange,_overValue,_peerActive,_peerChecked,_peerDisabled,_peerExpanded,_peerFocus,_peerFocusVisible,_peerFocusWithin,_peerHover,_peerInvalid,_peerPlaceholderShown,_placeholder,_placeholderShown,_portrait,_pressed,_print,_rangeEnd,_rangeStart,_readOnly,_readWrite,_required,_rtl,_scrollbar,_scrollbarThumb,_scrollbarTrack,_selected,_selection,_starting,_target,_today,_topmost,_unavailable,_underValue,_valid,_vertical,_visited,base,lg,lgDown,lgOnly,lgTo2xl,lgToXl,md,mdDown,mdOnly,mdTo2xl,mdToLg,mdToXl,sm,smDown,smOnly,smTo2xl,smToLg,smToMd,smToXl,xl,xlDown,xlOnly,xlTo2xl".split(','))
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