import React, { type ReactNode } from "react";
import Translate from "@docusaurus/Translate";
import { ThemeClassNames } from "@docusaurus/theme-common";
import Link from "@docusaurus/Link";
import IconEdit from "@theme/Icon/Edit";
import type { Props } from "@theme/EditThisPage";
import { Button } from "@mui/joy";
import { MdEdit } from "react-icons/md";
import { FaGithub } from "react-icons/fa";

export default function EditThisPage({ editUrl }: Props): ReactNode {
  return (
    <Button
      variant="outlined"
      sx={{
        textDecoration: "none",
        "&:hover": {
          textDecoration: "none",
        },
      }}
      startDecorator={<FaGithub />}
      component={Link}
      to={editUrl}
    >
      Edit this page
    </Button>
  );
}
