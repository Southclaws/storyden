"use client";

import type { Assign } from "@ark-ui/react";
import {
  Progress as ArkProgress,
  type ProgressRootProps,
} from "@ark-ui/react/progress";

import { css, cx } from "@/styled-system/css";
import { splitCssProps } from "@/styled-system/jsx";
import { type ProgressVariantProps, progress } from "@/styled-system/recipes";
import type { JsxStyleProps } from "@/styled-system/types";

interface BaseProgressProps
  extends Assign<JsxStyleProps, ProgressRootProps>,
    ProgressVariantProps {}

export interface ProgressCircleProps extends BaseProgressProps {
  showValue?: boolean;
}

export const ProgressCircle = (props: ProgressCircleProps) => {
  const [variantProps, progressProps] = progress.splitVariantProps(props);
  const [cssProps, localProps] = splitCssProps(progressProps);
  const { showValue = true, className, ...rootProps } = localProps;
  const styles = progress(variantProps);

  return (
    <ArkProgress.Root
      className={cx(styles.root, css(cssProps), className)}
      {...rootProps}
    >
      <ArkProgress.Circle className={styles.circle}>
        <ArkProgress.CircleTrack className={styles.circleTrack} />
        <ArkProgress.CircleRange className={styles.circleRange} />
      </ArkProgress.Circle>
      {showValue && <ArkProgress.ValueText className={styles.valueText} />}
    </ArkProgress.Root>
  );
};

export interface ProgressHorizontalProps extends BaseProgressProps {
  showValue?: boolean;
}

export const ProgressHorizontal = (props: ProgressHorizontalProps) => {
  const [variantProps, progressProps] = progress.splitVariantProps(props);
  const [cssProps, localProps] = splitCssProps(progressProps);
  const { showValue = true, className, ...rootProps } = localProps;
  const styles = progress(variantProps);

  return (
    <ArkProgress.Root
      className={cx(styles.root, css(cssProps), className)}
      {...rootProps}
    >
      <ArkProgress.Track className={styles.track}>
        <ArkProgress.Range className={styles.range} />
      </ArkProgress.Track>
      {showValue && <ArkProgress.ValueText className={styles.valueText} />}
    </ArkProgress.Root>
  );
};
