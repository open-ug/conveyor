import React, { type ReactNode } from "react";
import clsx from "clsx";
import { ThemeClassNames } from "@docusaurus/theme-common";
import { useSidebarBreadcrumbs } from "@docusaurus/plugin-content-docs/client";
import { useHomePageRoute } from "@docusaurus/theme-common/internal";
import Link from "@docusaurus/Link";
import { translate } from "@docusaurus/Translate";
import HomeBreadcrumbItem from "@theme/DocBreadcrumbs/Items/Home";
import { Link as MuiLink } from "@mui/joy";

import styles from "./styles.module.css";
import { Breadcrumbs, Typography } from "@mui/joy";
import { MdChevronRight } from "react-icons/md";

// TODO move to design system folder
function BreadcrumbsItemLink({
  children,
  href,
  isLast,
}: {
  children: ReactNode;
  href: string | undefined;
  isLast: boolean;
}): ReactNode {
  const className = "breadcrumbs__link";
  if (isLast) {
    return (
      <span className={className} itemProp="name">
        {children}
      </span>
    );
  }
  return href ? (
    <Link className={className} href={href} itemProp="item">
      <span itemProp="name">{children}</span>
    </Link>
  ) : (
    // TODO Google search console doesn't like breadcrumb items without href.
    // The schema doesn't seem to require `id` for each `item`, although Google
    // insist to infer one, even if it's invalid. Removing `itemProp="item
    // name"` for now, since I don't know how to properly fix it.
    // See https://github.com/facebook/docusaurus/issues/7241
    <span className={className}>{children}</span>
  );
}

// TODO move to design system folder
function BreadcrumbsItem({
  children,
  active,
  index,
  addMicrodata,
}: {
  children: ReactNode;
  active?: boolean;
  index: number;
  addMicrodata: boolean;
}): ReactNode {
  return (
    <li
      {...(addMicrodata && {
        itemScope: true,
        itemProp: "itemListElement",
        itemType: "https://schema.org/ListItem",
      })}
      className={clsx("breadcrumbs__item", {
        "breadcrumbs__item--active": active,
      })}
    >
      {children}
      <meta itemProp="position" content={String(index + 1)} />
    </li>
  );
}

export default function DocBreadcrumbs(): ReactNode {
  const breadcrumbs = useSidebarBreadcrumbs();
  const homePageRoute = useHomePageRoute();

  if (!breadcrumbs) {
    return null;
  }

  return (
    <nav
      className={clsx(
        ThemeClassNames.docs.docBreadcrumbs,
        styles.breadcrumbsContainer
      )}
      aria-label={translate({
        id: "theme.docs.breadcrumbs.navAriaLabel",
        message: "Breadcrumbs",
        description: "The ARIA label for the breadcrumbs",
      })}
    >
      <Breadcrumbs
        size="sm"
        aria-label="breadcrumbs"
        separator={<MdChevronRight fontSize="sm" />}
        itemScope
        itemType="https://schema.org/BreadcrumbList"
      >
        {homePageRoute && <HomeBreadcrumbItem />}
        {breadcrumbs.map((item, idx) => {
          const isLast = idx === breadcrumbs.length - 1;
          const href =
            item.type === "category" && item.linkUnlisted
              ? undefined
              : item.href;
          return (
            <>
              {isLast ? (
                <Typography color="primary" fontWeight={500} fontSize={12}>
                  {item.label}
                </Typography>
              ) : (
                <MuiLink
                  underline="hover"
                  color="neutral"
                  fontSize={12}
                  fontWeight={500}
                  component={Link}
                  href={href}
                >
                  {item.label}
                </MuiLink>
              )}
            </>
          );
        })}
      </Breadcrumbs>
    </nav>
  );
}
