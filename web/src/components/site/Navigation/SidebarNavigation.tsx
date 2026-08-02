import Link from "next/link";
import { type MouseEventHandler, type ReactNode } from "react";

type PaneProps = {
  label: string;
  title: string;
  children: ReactNode;
  footer?: ReactNode;
};

export function SidebarNavigationPane({
  label,
  title,
  children,
  footer,
}: PaneProps) {
  return (
    <header className="sidebar-navigation-pane navigation__surface">
      <nav className="sidebar-navigation" aria-label={label}>
        <div className="sidebar-navigation__header">
          <h2 className="sidebar-navigation__title">{title}</h2>
        </div>
        <div className="sidebar-navigation__body">{children}</div>
        {footer && <div className="sidebar-navigation__footer">{footer}</div>}
      </nav>
    </header>
  );
}

type SectionProps = {
  label?: string;
  children: ReactNode;
};

export function SidebarNavigationSection({ label, children }: SectionProps) {
  return (
    <section className="sidebar-navigation__section">
      {label && <h3 className="sidebar-navigation__section-label">{label}</h3>}
      <div className="sidebar-navigation__items">{children}</div>
    </section>
  );
}

type LinkProps = {
  href: string;
  label: string;
  icon: ReactNode;
  active?: boolean;
  description?: string;
  onClick?: MouseEventHandler<HTMLAnchorElement>;
};

export function SidebarNavigationLink({
  href,
  label,
  icon,
  active,
  description,
  onClick,
}: LinkProps) {
  return (
    <Link
      className="sidebar-navigation__link"
      href={href}
      aria-current={active ? "page" : undefined}
      onClick={onClick}
    >
      <span className="sidebar-navigation__icon" aria-hidden="true">
        {icon}
      </span>
      <span className="sidebar-navigation__link-content">
        <span className="sidebar-navigation__link-label">{label}</span>
        {description && (
          <span className="sidebar-navigation__link-description">
            {description}
          </span>
        )}
      </span>
    </Link>
  );
}
