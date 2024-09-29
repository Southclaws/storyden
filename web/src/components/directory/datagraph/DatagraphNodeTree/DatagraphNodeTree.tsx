import type { Assign } from "@ark-ui/react";
import {
  TreeView as ArkTreeView,
  type TreeViewRootProps,
} from "@ark-ui/react/tree-view";
import Link from "next/link";
import { forwardRef } from "react";

import { NodeWithChildren } from "@/api/openapi-schema";
import { css, cx } from "@/styled-system/css";
import { splitCssProps } from "@/styled-system/jsx";
import { type TreeViewVariantProps, treeView } from "@/styled-system/recipes";
import { token } from "@/styled-system/tokens";
import type { JsxStyleProps } from "@/styled-system/types";

import { DatagraphNodeMenu } from "../DatagraphNodeMenu/DatagraphNodeMenu";

import { useDatagraphNodeTree } from "./useDatagraphNodeTree";

export interface TreeViewData {
  label: string;
  children: NodeWithChildren[];
}

export interface TreeViewProps
  extends Assign<JsxStyleProps, TreeViewRootProps>,
    TreeViewVariantProps {
  data: TreeViewData;
  currentNode: string | undefined;
}

export const DatagraphNodeTree = forwardRef<HTMLDivElement, TreeViewProps>(
  (props, ref) => {
    const [cssProps, localProps] = splitCssProps(props);
    const { data, currentNode, className, ...rootProps } = localProps;

    const { handleDelete } = useDatagraphNodeTree(currentNode);

    const styles = treeView();

    const defaultExpandedValue: string[] = [];

    // recursively find currentNode in data and add each parent to defaultExpandedValue
    const findCurrentNode = (node: NodeWithChildren) => {
      if (node.slug === currentNode) {
        defaultExpandedValue.push(node.slug);
        return true;
      }

      if (node.children) {
        for (const child of node.children) {
          if (findCurrentNode(child)) {
            defaultExpandedValue.push(node.slug);
            return true;
          }
        }
      }

      return false;
    };

    data.children.forEach(findCurrentNode);

    const renderChild = (child: NodeWithChildren) => {
      return (
        <ArkTreeView.Branch
          key={child.id}
          value={child.slug}
          className={styles.branch}
        >
          <TreeBranch
            styles={styles}
            child={child}
            handleDelete={handleDelete}
          />

          <ArkTreeView.BranchContent className={styles.branchContent}>
            {child.children?.map((child) =>
              child.children ? (
                renderChild(child)
              ) : (
                <TreeItem
                  key={child.id}
                  styles={styles}
                  child={child}
                  handleDelete={handleDelete}
                />
              ),
            )}
          </ArkTreeView.BranchContent>
        </ArkTreeView.Branch>
      );
    };

    return (
      <ArkTreeView.Root
        ref={ref}
        aria-label={data.label}
        className={cx(styles.root, css(cssProps), className)}
        defaultExpandedValue={defaultExpandedValue}
        {...rootProps}
      >
        <ArkTreeView.Tree className={styles.tree}>
          {data.children.map(renderChild)}
        </ArkTreeView.Tree>
      </ArkTreeView.Root>
    );
  },
);

DatagraphNodeTree.displayName = "DatagraphNodeTree";

type BranchProps = {
  child: NodeWithChildren;
  styles: any;
  handleDelete: (slug: string) => void;
};

function TreeBranch({ styles, child, handleDelete }: BranchProps) {
  return (
    <ArkTreeView.BranchControl className={cx(styles.branchControl)}>
      <ArkTreeView.BranchIndicator className={styles.branchIndicator}>
        {child.children?.length ? <ChevronRightIcon /> : <BulletIcon />}
      </ArkTreeView.BranchIndicator>

      <ArkTreeView.BranchText asChild className={styles.branchText}>
        <Link href={`/directory/${child.slug}`}>{child.name}</Link>
      </ArkTreeView.BranchText>

      <DatagraphNodeMenu node={child} onDelete={handleDelete} />
    </ArkTreeView.BranchControl>
  );
}

function TreeItem({ styles, child }: BranchProps) {
  return (
    <ArkTreeView.Item value={child.slug} className={cx(styles.item)}>
      <ArkTreeView.ItemText className={styles.itemText}>
        <Link href={`/directory/${child.slug}`}>{child.name}</Link>
      </ArkTreeView.ItemText>
    </ArkTreeView.Item>
  );
}

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
