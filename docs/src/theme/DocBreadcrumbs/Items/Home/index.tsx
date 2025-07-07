import React, { type ReactNode } from "react";
import Link from "@docusaurus/Link";
import useBaseUrl from "@docusaurus/useBaseUrl";
import { translate } from "@docusaurus/Translate";
import IconHome from "@theme/Icon/Home";
import { Link as MuiLink } from "@mui/joy";

import styles from "./styles.module.css";
import { MdHome } from "react-icons/md";

export default function HomeBreadcrumbItem(): ReactNode {
  const homeHref = useBaseUrl("/");

  return (
    <MuiLink
      aria-label={translate({
        id: "theme.docs.breadcrumbs.home",
        message: "Home page",
        description: "The ARIA label for the home page in the breadcrumbs",
      })}
      component={Link}
      href={homeHref}
      underline="none"
      color="neutral"
    >
      <MdHome />
    </MuiLink>
  );
}
