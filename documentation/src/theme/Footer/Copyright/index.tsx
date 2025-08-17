import React, { type ReactNode } from "react";
import type { Props } from "@theme/Footer/Copyright";
import { Link, Typography } from "@mui/joy";

export default function FooterCopyright({ copyright }: Props): ReactNode {
  console.log("FooterCopyright", copyright);
  return (
    <Typography
      level="body-sm"
      sx={{
        color: "white",
        textAlign: "center",
        marginTop: 2,
      }}
    >
      Conveyor CI is a free and open-source project. Copyrighted under the
      Apache License 2.0. Copyright © {new Date().getFullYear()}{" "}
      <Link component="a" href="https://open.ug">
        Open UG Labs
      </Link>{" "}
      & Conveyor CI contributors.
    </Typography>
  );
}
