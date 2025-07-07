import { useState } from "react";
import {
  Box,
  Typography,
  Container,
  Button,
  Card,
  Stack,
  Divider,
  useTheme,
} from "@mui/joy";
import { FaRocket, FaGithub, FaSlack, FaArrowRight } from "react-icons/fa";

const CallToActionSection = () => {
  const theme = useTheme();
  const [hovered, setHovered] = useState(null);

  const ctaOptions = [
    {
      id: "docs",
      title: "Read the Docs",
      description: "Comprehensive guides and API references to get you started",
      icon: <FaRocket />,
      buttonText: "View Documentation",
      buttonAction: "/docs/introduction",
      primary: true,
    },
    {
      id: "github",
      title: "Explore Source",
      description:
        "Check out our GitHub repository and contribute to the project",
      icon: <FaGithub />,
      buttonText: "Visit GitHub",
      buttonAction: "https://github.com/open-ug/conveyor",
      primary: false,
    },
    {
      id: "community",
      title: "Join Community",
      description: "Connect with other developers using Conveyor CI",
      icon: <FaSlack />,
      buttonText: "Join Slack",
      buttonAction:
        "https://open-ug.slack.com/join/shared_invite/zt-382rwaqzr-aeB5ES_S94fYVPKCGHW~Qg",
      primary: false,
    },
  ];

  return (
    <Box
      sx={{
        py: 10,
        background: `linear-gradient(180deg, ${theme.palette.background.body} 0%, ${theme.palette.background.level1} 100%)`,
        position: "relative",
        overflow: "hidden",
      }}
    >
      {/* Decorative background elements */}
      <Box
        sx={{
          position: "absolute",
          top: -50,
          left: -100,
          width: 300,
          height: 300,
          borderRadius: "50%",
          background: `radial-gradient(circle, ${theme.palette.primary[100]} 0%, transparent 70%)`,
          opacity: 0.4,
          zIndex: 0,
        }}
      />

      <Box
        sx={{
          position: "absolute",
          bottom: -80,
          right: -80,
          width: 300,
          height: 300,
          borderRadius: "50%",
          background: `radial-gradient(circle, ${theme.palette.primary[100]} 0%, transparent 70%)`,
          opacity: 0.3,
          zIndex: 0,
        }}
      />

      <Container maxWidth="lg" sx={{ position: "relative", zIndex: 1 }}>
        <Box sx={{ textAlign: "center", mb: 6 }}>
          <Typography
            level="h2"
            sx={{
              fontWeight: "bold",
              mb: 2,
              fontSize: { xs: "2rem", md: "2.5rem" },
            }}
          >
            Ready to Transform Your CI/CD Experience?
          </Typography>
          <Typography
            level="body-lg"
            sx={{
              color: "text.secondary",
              maxWidth: "700px",
              mx: "auto",
            }}
          >
            Get started with Conveyor CI today and build a robust, scalable
            pipeline in minutes.
          </Typography>
        </Box>

        <Box
          sx={{
            display: "flex",
            flexDirection: { xs: "column", md: "row" },
            gap: 3,
            justifyContent: "center",
          }}
        >
          {ctaOptions.map((option) => (
            <Card
              key={option.id}
              variant={option.primary ? "solid" : "outlined"}
              color={option.primary ? "primary" : undefined}
              sx={{
                flex: 1,
                p: 3,
                display: "flex",
                flexDirection: "column",
                minWidth: "240px",
                maxWidth: { xs: "100%", md: "400px" },
                transition: "all 0.3s",
                "&:hover": {
                  boxShadow: "md",
                  transform: "translateY(-5px)",
                },
                backgroundColor: option.primary
                  ? undefined
                  : "background.surface",
              }}
              onMouseEnter={() => setHovered(option.id)}
              onMouseLeave={() => setHovered(null)}
            >
              <Box
                sx={{
                  fontSize: "2rem",
                  mb: 2,
                  color: option.primary ? "primary.solidBg" : "primary.500",
                }}
              >
                {option.icon}
              </Box>
              <Typography level="title-lg" sx={{ mb: 1, fontWeight: "bold" }}>
                {option.title}
              </Typography>
              <Typography level="body-md" sx={{ mb: 3, flexGrow: 1 }}>
                {option.description}
              </Typography>
              <Button
                variant={option.primary ? "solid" : "outlined"}
                color={option.primary ? "neutral" : "primary"}
                fullWidth
                size="lg"
                component="a"
                href={option.buttonAction}
                endDecorator={<FaArrowRight />}
                sx={{
                  mt: "auto",
                  borderWidth: option.primary ? 0 : 2,
                  backgroundColor: option.primary ? "white" : undefined,
                  color: option.primary ? "primary.600" : undefined,
                  "&:hover": {
                    backgroundColor: option.primary ? "white" : undefined,
                    opacity: option.primary ? 0.9 : undefined,
                  },
                  transition: "all 0.2s",
                  transform:
                    hovered === option.id ? "translateY(-2px)" : "none",
                }}
              >
                {option.buttonText}
              </Button>
            </Card>
          ))}
        </Box>

        <Divider sx={{ my: 8 }} />

        {/*<Stack
          direction={{ xs: "column", md: "row" }}
          spacing={3}
          alignItems="center"
          justifyContent="center"
          sx={{ textAlign: { xs: "center", md: "left" } }}
        >
          <Box sx={{ flex: 1, maxWidth: "600px" }}>
            <Typography level="h3" sx={{ fontWeight: "bold", mb: 2 }}>
              Need Enterprise Support?
            </Typography>
            <Typography level="body-md" sx={{ color: "text.secondary" }}>
              Our team provides dedicated support, training, and customization
              options for enterprise teams looking to implement Conveyor CI at
              scale.
            </Typography>
          </Box>
          <Button
            variant="outlined"
            color="neutral"
            size="lg"
            sx={{
              px: 4,
              borderColor: theme.palette.divider,
              "&:hover": {
                borderColor: theme.palette.primary[400],
              },
            }}
          >
            Contact Sales
          </Button>
        </Stack>*/}
      </Container>
    </Box>
  );
};

export default CallToActionSection;
