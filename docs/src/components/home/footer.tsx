import React, { useState } from "react";
import {
  Box,
  Typography,
  Button,
  Container,
  Stack,
  Chip,
  Divider,
  Input,
  IconButton,
} from "@mui/joy";
import { FaGithub } from "react-icons/fa";

const Footer = () => {
  const [email, setEmail] = useState("");
  const [subscribed, setSubscribed] = useState(false);

  const handleSubscribe = () => {
    if (email) {
      setSubscribed(true);
      setEmail("");
      setTimeout(() => setSubscribed(false), 3000);
    }
  };

  const footerSections = [
    {
      title: "Quick Links",
      links: [
        { label: "Tutorial", href: "/docs/category/quick-start-tutorial" },
        { label: "Installation", href: "/docs/installation" },
        { label: "Documentation", href: "/docs/introduction" },
        { label: "API Reference", href: "/docs/usage/using-api" },
      ],
    },
    {
      title: "Community",
      links: [
        {
          label: "GitHub",
          href: "https://github.com/open-ug/conveyor",
          external: true,
        },
        { label: "Contributing", href: "/docs/contributing/how-to-contribute" },
        {
          label: "Code of Conduct",
          href: "https://github.com/open-ug/conveyor/blob/main/CODE_OF_CONDUCT.md",
          external: true,
        },
        { label: "Roadmap", href: "/docs/contributing/roadmap" },
      ],
    },
  ];

  return (
    <Box
      component="footer"
      sx={{
        background: "linear-gradient(180deg, #0f172a 0%, #020617 100%)",
        color: "white",
        position: "relative",
        overflow: "hidden",
      }}
    >
      {/* Background decoration */}
      <Box
        sx={{
          position: "absolute",
          top: 0,
          left: 0,
          right: 0,
          bottom: 0,
          opacity: 0.05,
          background: `
            radial-gradient(circle at 20% 20%, #3b82f6 0%, transparent 50%),
            radial-gradient(circle at 80% 80%, #1d4ed8 0%, transparent 50%)
          `,
        }}
      />

      <Container maxWidth="xl" sx={{ position: "relative", zIndex: 1 }}>
        {/* Main Footer Content */}
        <Box sx={{ py: 8 }}>
          <Box
            sx={{
              display: "grid",
              gridTemplateColumns: {
                xs: "1fr",
                sm: "repeat(2, 1fr)",
                md: "repeat(3, 1fr)",
                lg: "repeat(3, 1fr)",
              },
              gap: 4,
              mb: 8,
            }}
          >
            {/* Brand Column */}
            <Box sx={{ gridColumn: { lg: "span 1" } }}>
              <Box
                sx={{ display: "flex", alignItems: "center", gap: 2, mb: 3 }}
              >
                <Box
                  sx={{
                    width: 40,
                    height: 40,
                    borderRadius: 8,
                    display: "flex",
                    alignItems: "center",
                    justifyContent: "center",
                    fontSize: "1.2rem",
                  }}
                >
                  <img src="/logos/icon.svg" alt="Conveyor CI" />
                </Box>
                <Typography
                  level="h4"
                  sx={{
                    fontWeight: 700,
                    background:
                      "linear-gradient(135deg, #ffffff 0%, #3b82f6 100%)",
                    backgroundClip: "text",
                    WebkitBackgroundClip: "text",
                    WebkitTextFillColor: "transparent",
                  }}
                >
                  Conveyor CI
                </Typography>
              </Box>
              <Typography
                level="body-sm"
                sx={{
                  color: "#94a3b8",
                  lineHeight: 1.6,
                  mb: 3,
                }}
              >
                The lightweight, distributed CI/CD engine built for platform
                developers who demand simplicity without compromise.
              </Typography>
              <Stack spacing={1}>
                <Chip
                  //variant="outlined"
                  size="sm"
                  sx={{
                    color: "#10b981",
                    width: "fit-content",
                    bgcolor: "rgba(16, 185, 129, 0.1)",
                  }}
                >
                  Apache 2.0 License
                </Chip>
              </Stack>
            </Box>

            {/* Footer Sections */}
            {footerSections.map((section, index) => (
              <Box key={index}>
                <Typography
                  level="title-sm"
                  sx={{
                    color: "white",
                    fontWeight: 600,
                    mb: 2,
                  }}
                >
                  {section.title}
                </Typography>
                <Stack spacing={1}>
                  {section.links.map((link, linkIndex) => (
                    <Box
                      key={linkIndex}
                      component="a"
                      href={link.href}
                      target={link.external ? "_blank" : undefined}
                      rel={link.external ? "noopener noreferrer" : undefined}
                      sx={{
                        color: "#94a3b8",
                        textDecoration: "none",
                        fontSize: "14px",
                        cursor: "pointer",
                        display: "flex",
                        alignItems: "center",
                        gap: 0.5,
                        "&:hover": {
                          color: "primary.300",
                        },
                        transition: "color 0.2s ease",
                      }}
                    >
                      {link.label}
                      {link.external && (
                        <Typography sx={{ fontSize: "10px", color: "#64748b" }}>
                          ↗
                        </Typography>
                      )}
                    </Box>
                  ))}
                </Stack>
              </Box>
            ))}
          </Box>
        </Box>

        <Divider sx={{ borderColor: "#334155" }} />

        {/* Bottom Bar */}
        <Box
          sx={{
            py: 4,
            display: "flex",
            flexDirection: { xs: "column", md: "row" },
            justifyContent: "space-between",
            alignItems: "center",
            gap: 2,
          }}
        >
          <Typography
            level="body-sm"
            sx={{ color: "#64748b", textAlign: { xs: "center", md: "left" } }}
          >
            © {new Date().getFullYear()} Conveyor CI Contributors. All rights
            reserved. Built with ❤️ for platform developers.
          </Typography>

          <Stack
            direction={{ xs: "column", sm: "row" }}
            spacing={2}
            sx={{ alignItems: "center" }}
          >
            <Box sx={{ display: "flex", alignItems: "center", gap: 1 }}>
              <Typography level="body-xs" sx={{ color: "#64748b" }}>
                Version:
              </Typography>
              <Chip
                sx={{
                  borderColor: "primary.300",
                  color: "primary.300",
                  bgcolor: "rgba(59, 130, 246, 0.1)",
                }}
                variant="outlined"
                size="sm"
              >
                v0.4.0
              </Chip>
            </Box>

            <Button
              variant="outlined"
              size="sm"
              sx={{
                borderColor: "primary.700",
                color: "primary.300",
                fontSize: "12px",
                "&:hover": {
                  bgcolor: "rgba(59, 130, 246, 0.1)",
                  borderColor: "primary.500",
                },
              }}
              startDecorator={<FaGithub />}
              component="a"
              href="https://github.com/open-ug/conveyor"
            >
              Star on GitHub
            </Button>
          </Stack>
        </Box>
      </Container>

      <style>{`
        @keyframes pulse {
          0%,
          100% {
            opacity: 1;
          }
          50% {
            opacity: 0.5;
          }
        }
      `}</style>
    </Box>
  );
};

export default Footer;
