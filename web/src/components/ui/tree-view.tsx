import type { Assign } from "@ark-ui/react";
import {
  TreeView as ArkTreeView,
  type TreeViewRootProps,
} from "@ark-ui/react/tree-view";
import { CircleDashed, CircleIcon, DotIcon } from "lucide-react";
import Link from "next/link";
import { forwardRef } from "react";

import { css, cx } from "@/styled-system/css";
import { splitCssProps } from "@/styled-system/jsx";
import { type TreeViewVariantProps, treeView } from "@/styled-system/recipes";
import { token } from "@/styled-system/tokens";
import type { JsxStyleProps } from "@/styled-system/types";

export interface Child {
  value: string;
  name: string;
  url: string;
  children?: Child[];
}

export interface TreeViewData {
  label: string;
  children: Child[];
}

export interface TreeViewProps
  extends Assign<JsxStyleProps, TreeViewRootProps>,
    TreeViewVariantProps {
  data: TreeViewData;
}

export const TreeView = forwardRef<HTMLDivElement, TreeViewProps>(
  (props, ref) => {
    const [cssProps, localProps] = splitCssProps(props);
    const { data, className, ...rootProps } = localProps;
    const styles = treeView();

    const renderChild = (child: Child) => (
      <ArkTreeView.Branch
        key={child.value}
        value={child.value}
        className={styles.branch}
      >
        <ArkTreeView.BranchControl className={styles.branchControl}>
          <ArkTreeView.BranchIndicator className={styles.branchIndicator}>
            {child.children?.length ? <ChevronRightIcon /> : <BulletIcon />}
          </ArkTreeView.BranchIndicator>

          <ArkTreeView.BranchText className={styles.branchText}>
            <Link href={child.url}>{child.name}</Link>
          </ArkTreeView.BranchText>
        </ArkTreeView.BranchControl>

        <ArkTreeView.BranchContent className={styles.branchContent}>
          {child.children?.map((child) =>
            child.children ? (
              renderChild(child)
            ) : (
              <ArkTreeView.Item
                key={child.value}
                value={child.value}
                className={styles.item}
              >
                <ArkTreeView.ItemText className={styles.itemText}>
                  <Link href={child.url}>{child.name}</Link>
                </ArkTreeView.ItemText>
              </ArkTreeView.Item>
            ),
          )}
        </ArkTreeView.BranchContent>
      </ArkTreeView.Branch>
    );

    return (
      <ArkTreeView.Root
        ref={ref}
        aria-label={data.label}
        className={cx(styles.root, css(cssProps), className)}
        {...rootProps}
      >
        <ArkTreeView.Tree className={styles.tree}>
          {data.children.map(renderChild)}
        </ArkTreeView.Tree>
      </ArkTreeView.Root>
    );
  },
);

TreeView.displayName = "TreeView";

const ChevronRightIcon = () => (
  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
    <title>Chevron Right Icon</title>
    <path
      fill="none"
      stroke="currentColor"
      strokeLinecap="round"
      strokeLinejoin="round"
      strokeWidth="2"
      d="m9 18l6-6l-6-6"
    />
  </svg>
);

const BulletIcon = () => (
  <svg
    xmlns="http://www.w3.org/2000/svg"
    width="16"
    height="16"
    viewBox="0 0 24 24"
    fill={token("colors.fg.muted")}
  >
    <circle cx="12.1" cy="12.1" r="2.5" />
  </svg>
);
