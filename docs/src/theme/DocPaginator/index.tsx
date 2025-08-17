import React, { type ReactNode } from "react";
import Translate, { translate } from "@docusaurus/Translate";
import PaginatorNavLink from "@theme/PaginatorNavLink";
import type { Props } from "@theme/DocPaginator";
import Link from "@docusaurus/Link";
import { Button } from "@mui/joy";
import { MdArrowBack, MdNavigateBefore, MdNavigateNext } from "react-icons/md";

export default function DocPaginator(props: Props): ReactNode {
  const { previous, next } = props;
  return (
    <nav
      className="pagination-nav docusaurus-mt-lg"
      aria-label={translate({
        id: "theme.docs.paginator.navAriaLabel",
        message: "Docs pages",
        description: "The ARIA label for the docs pagination",
      })}
    >
      {previous && (
        <>
          <Button
            variant="soft"
            color="neutral"
            size="lg"
            sx={{
              textDecoration: "none",
              "&:hover": {
                textDecoration: "none",
              },
            }}
            startDecorator={<MdNavigateBefore />}
            component={Link}
            to={previous.permalink}
          >
            {previous.title}
          </Button>
        </>
      )}
      {next && (
        <Button
          variant="soft"
          color="neutral"
          size="lg"
          sx={{
            textDecoration: "none",
            "&:hover": {
              textDecoration: "none",
            },
          }}
          endDecorator={<MdNavigateNext />}
          component={Link}
          to={next.permalink}
        >
          {next.title}
        </Button>
      )}
    </nav>
  );
}
