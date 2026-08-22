import { type ReactNode } from "react";

import { cx } from "@/styled-system/css";
import { styled } from "@/styled-system/jsx";
import { pageHeader } from "@/styled-system/recipes";
import type { HTMLStyledProps } from "@/styled-system/types";

import { PageHeading } from "../page-heading";
import { Text } from "../text";

export type PageHeaderProps = Omit<HTMLStyledProps<"header">, "title"> & {
  title: ReactNode;
  description?: ReactNode;
  badge?: ReactNode;
  back?: ReactNode;
  navigation?: ReactNode;
  actions?: ReactNode;
};

export function PageHeader({
  title,
  description,
  badge,
  back,
  navigation,
  actions,
  className,
  ...props
}: PageHeaderProps) {
  const styles = pageHeader();

  return (
    <styled.header className={cx(styles.root, className)} {...props}>
      {navigation && <div className={styles.navigation}>{navigation}</div>}
      <div className={styles.row}>
        <div className={styles.heading}>
          {back}
          <div className={styles.titleGroup}>
            <div className={styles.titleRow}>
              <PageHeading>{title}</PageHeading>
              {badge}
            </div>
            {description && <Text variant="supporting">{description}</Text>}
          </div>
        </div>
        {actions && <div className={styles.actions}>{actions}</div>}
      </div>
    </styled.header>
  );
}
