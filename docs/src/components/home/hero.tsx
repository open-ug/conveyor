import { useState } from "react";
import {
  Box,
  Typography,
  Button,
  Container,
  Grid,
  Stack,
  useTheme,
} from "@mui/joy";

const HeroSection = () => {
  const theme = useTheme();
  const [isHovered, setIsHovered] = useState(false);

  return (
    <Box
      sx={{
        background: "linear-gradient(135deg, #0A1929 0%, #164B85 100%)",
        pt: 12,
        pb: 12,
        color: "white",
        overflow: "hidden",
        position: "relative",
      }}
    >
      {/* Decorative background elements */}
      <Box
        sx={{
          position: "absolute",
          top: 0,
          right: 0,
          width: "30%",
          height: "100%",
          opacity: 0.06,
          background:
            "radial-gradient(circle, rgba(255,255,255,0.8) 0%, rgba(255,255,255,0) 70%)",
        }}
      />

      <Container maxWidth="lg">
        <Grid container spacing={6} alignItems="center">
          <Grid xs={12} md={7}>
            <Stack spacing={4}>
              <Typography
                level="h1"
                fontWeight="bold"
                fontSize={{ xs: "2.5rem", md: "3.5rem" }}
                lineHeight={1.2}
                sx={{
                  color: "whitesmoke",
                }}
              >
                Build Your CI/CD Platform,
                <Typography
                  component="span"
                  sx={{
                    color: theme.palette.primary[300],
                    display: "block",
                  }}
                >
                  Without The Complexity
                </Typography>
              </Typography>

              <Typography
                level="body-md"
                fontSize="lg"
                sx={{ maxWidth: "600px", color: "white" }}
              >
                Conveyor CI simplifies developing scalable and reliable CI/CD
                platforms, eliminating common pain points with a streamlined,
                developer-friendly approach.
              </Typography>

              <Stack direction="row" spacing={2} sx={{ mt: 2 }}>
                <Button
                  variant="solid"
                  size="lg"
                  component="a"
                  href="/docs/tutorial/installation"
                  color="primary"
                  onMouseEnter={() => setIsHovered(true)}
                  onMouseLeave={() => setIsHovered(false)}
                  sx={{
                    px: 4,
                    py: 1.5,
                    fontWeight: 600,
                    textDecoration: "none",
                    "&:hover": {
                      textDecoration: "none",
                      color: "white",
                    },
                    transition: "all 0.2s",
                    transform: isHovered ? "translateY(-3px)" : "none",
                    boxShadow: isHovered ? "md" : "sm",
                  }}
                >
                  Get Started
                </Button>
                <Button
                  variant="outlined"
                  size="lg"
                  component="a"
                  href="/docs/introduction"
                  color="neutral"
                  sx={{
                    px: 4,
                    py: 1.5,
                    fontWeight: 600,
                    textDecoration: "none",
                    color: "white",
                    borderColor: "rgba(255,255,255,0.4)",
                    "&:hover": {
                      borderColor: "white",
                      backgroundColor: "rgba(255,255,255,0.04)",
                      textDecoration: "none",
                      color: "white",
                    },
                  }}
                >
                  View Documentation
                </Button>
              </Stack>
            </Stack>
          </Grid>

          <Grid xs={12} md={5}>
            <Box
              sx={{
                backgroundColor: "rgba(0,0,0,0.6)",
                borderRadius: "lg",
                p: 2,
                border: "1px solid rgba(255,255,255,0.1)",
                boxShadow: "lg",
                position: "relative",
                "&:before": {
                  content: '""',
                  position: "absolute",
                  top: 0,
                  left: 0,
                  right: 0,
                  height: "4px",
                  background: "linear-gradient(90deg, #3498db, #1abc9c)",
                  borderTopLeftRadius: "lg",
                  borderTopRightRadius: "lg",
                },
              }}
            >
              <Box
                sx={{ display: "flex", alignItems: "center", mb: 2, gap: 1 }}
              >
                <Box sx={{ display: "flex", gap: 1 }}>
                  <Box
                    sx={{
                      width: 12,
                      height: 12,
                      borderRadius: "50%",
                      backgroundColor: "#ff5f57",
                    }}
                  />
                  <Box
                    sx={{
                      width: 12,
                      height: 12,
                      borderRadius: "50%",
                      backgroundColor: "#febc2e",
                    }}
                  />
                  <Box
                    sx={{
                      width: 12,
                      height: 12,
                      borderRadius: "50%",
                      backgroundColor: "#28c840",
                    }}
                  />
                </Box>
                <Typography
                  level="body-sm"
                  sx={{ color: "gray", ml: 1, fontSize: "xs" }}
                >
                  conveyor_config.yaml
                </Typography>
              </Box>

              <Box
                sx={{
                  fontFamily: "monospace",
                  fontSize: "sm",
                  color: "#f8f8f2",
                  whiteSpace: "pre-line",
                }}
              >
                <Typography
                  component="div"
                  sx={{ color: "#ff79c6", fontFamily: "monospace" }}
                >
                  # Define your CI/CD pipeline with Conveyor
                </Typography>
                <Typography
                  component="div"
                  sx={{ color: "#8be9fd", fontFamily: "monospace" }}
                >
                  pipeline:
                </Typography>
                <Typography
                  component="div"
                  sx={{ color: "#f1fa8c", fontFamily: "monospace", ml: 2 }}
                >
                  name: build-and-deploy
                </Typography>
                <Typography
                  component="div"
                  sx={{ color: "#8be9fd", fontFamily: "monospace", ml: 2 }}
                >
                  stages:
                </Typography>
                <Typography
                  component="div"
                  sx={{ color: "#f1fa8c", fontFamily: "monospace", ml: 4 }}
                >
                  - test
                </Typography>
                <Typography
                  component="div"
                  sx={{ color: "#f1fa8c", fontFamily: "monospace", ml: 4 }}
                >
                  - build
                </Typography>
                <Typography
                  component="div"
                  sx={{ color: "#f1fa8c", fontFamily: "monospace", ml: 4 }}
                >
                  - deploy
                </Typography>
                <Typography
                  component="div"
                  sx={{ color: "#8be9fd", fontFamily: "monospace", ml: 2 }}
                >
                  distributed: true
                </Typography>
                <Typography
                  component="div"
                  sx={{ color: "#8be9fd", fontFamily: "monospace", ml: 2 }}
                >
                  runners:
                </Typography>
                <Typography
                  component="div"
                  sx={{ color: "#f1fa8c", fontFamily: "monospace", ml: 4 }}
                >
                  - cloud-native
                </Typography>
              </Box>
            </Box>
          </Grid>
        </Grid>
      </Container>
    </Box>
  );
};

export default HeroSection;
