import React, { useState, useEffect } from "react";
import { Box, Typography, Button, Container, Stack, Chip } from "@mui/joy";
import { IoFlash } from "react-icons/io5";
import { SiKubernetes } from "react-icons/si";
import { VscGraph } from "react-icons/vsc";
import { MdOutlineArrowForward } from "react-icons/md";
import { FaGithub } from "react-icons/fa";
import Link from "@docusaurus/Link";

const CodeBlock = () => {
  const [currentStep, setCurrentStep] = useState(0);

  const codeSteps = [
    {
      title: "Define Pipeline",
      code: `# conveyor.yml
pipeline:
  name: "build-deploy"
  triggers: [push, pr]
  
stages:
  - name: test
    image: node:18
    commands:
      - npm test
      
  - name: deploy
    depends: [test]
    image: alpine/kubectl
    commands:
      - kubectl apply -f k8s/`,
    },
    {
      title: "Execute & Monitor",
      code: `$ conveyor run --pipeline build-deploy

✓ Pipeline started successfully
⚡ Real-time logs streaming...
📊 Observability dashboard: localhost:8080

Stage [test] ████████████ 100% Complete
Stage [deploy] ██████░░░░░░ 60% Running

🚀 Deployment successful in 2m 34s`,
    },
    {
      title: "Scale Horizontally",
      code: `$ conveyor scale --workers 5

✓ Scaling to 5 workers...
✓ Load balancer configured
✓ Event system synchronized

Current capacity: 50 concurrent pipelines
Average execution time: 1m 12s
System health: 🟢 Optimal

Ready to handle enterprise workloads!`,
    },
  ];

  useEffect(() => {
    const interval = setInterval(() => {
      setCurrentStep((prev) => (prev + 1) % codeSteps.length);
    }, 4000);
    return () => clearInterval(interval);
  }, []);

  return (
    <Box
      sx={{
        background: "linear-gradient(135deg, #0f172a 0%, #1e293b 100%)",
        borderRadius: 12,
        p: 3,
        fontFamily: 'Monaco, "Lucida Console", monospace',
        fontSize: "14px",
        minHeight: "280px",
        position: "relative",
        border: "1px solid",
        borderColor: "primary.300",
        boxShadow: "0 8px 32px rgba(59, 130, 246, 0.15)",
        overflow: "hidden",
      }}
    >
      <Box
        sx={{
          display: "flex",
          alignItems: "center",
          mb: 2,
          opacity: 0.8,
        }}
      >
        <Box sx={{ display: "flex", gap: 1, mr: 2 }}>
          <Box sx={{ w: 12, h: 12, borderRadius: "50%", bgcolor: "#ef4444" }} />
          <Box sx={{ w: 12, h: 12, borderRadius: "50%", bgcolor: "#f59e0b" }} />
          <Box sx={{ w: 12, h: 12, borderRadius: "50%", bgcolor: "#10b981" }} />
        </Box>
        <Typography level="body-sm" sx={{ color: "#64748b" }}>
          {codeSteps[currentStep].title}
        </Typography>
      </Box>

      <Box
        sx={{
          position: "relative",
          height: "200px",
          overflow: "hidden",
        }}
      >
        {codeSteps.map((step, index) => (
          <Box
            key={index}
            sx={{
              position: "absolute",
              top: 0,
              left: 0,
              right: 0,
              opacity: currentStep === index ? 1 : 0,
              transform: `translateY(${currentStep === index ? 0 : 20}px)`,
              transition: "all 0.6s ease-in-out",
              whiteSpace: "pre-wrap",
              color: "#e2e8f0",
              lineHeight: 1.5,
            }}
          >
            {step.code}
          </Box>
        ))}
      </Box>

      <Box
        sx={{
          position: "absolute",
          bottom: 16,
          right: 16,
          display: "flex",
          gap: 1,
        }}
      >
        {codeSteps.map((_, index) => (
          <Box
            key={index}
            sx={{
              w: 8,
              h: 8,
              borderRadius: "50%",
              bgcolor: currentStep === index ? "primary.400" : "neutral.600",
              transition: "all 0.3s ease",
            }}
          />
        ))}
      </Box>
    </Box>
  );
};

const FloatingCard = ({ children, delay = 0 }) => {
  const [isVisible, setIsVisible] = useState(false);

  useEffect(() => {
    const timer = setTimeout(() => setIsVisible(true), delay);
    return () => clearTimeout(timer);
  }, [delay]);

  return (
    <Box
      sx={{
        transform: isVisible ? "translateY(0px)" : "translateY(30px)",
        opacity: isVisible ? 1 : 0,
        transition: "all 0.8s cubic-bezier(0.4, 0, 0.2, 1)",
      }}
    >
      {children}
    </Box>
  );
};

const Hero = () => {
  const [mousePosition, setMousePosition] = useState({ x: 0, y: 0 });

  useEffect(() => {
    const handleMouseMove = (e) => {
      setMousePosition({
        x: (e.clientX / window.innerWidth) * 100,
        y: (e.clientY / window.innerHeight) * 100,
      });
    };

    window.addEventListener("mousemove", handleMouseMove);
    return () => window.removeEventListener("mousemove", handleMouseMove);
  }, []);

  return (
    <Box
      sx={{
        minHeight: "100vh",
        background: `
          radial-gradient(circle at ${mousePosition.x}% ${mousePosition.y}%, 
            rgba(59, 130, 246, 0.1) 0%, 
            transparent 50%),
          linear-gradient(135deg, 
            #0f172a 0%, 
            #1e293b 50%, 
            #0f172a 100%)
        `,
        position: "relative",
        display: "flex",
        alignItems: "center",
        overflow: "hidden",
        transition: "background 0.3s ease",
      }}
    >
      {/* Animated background elements */}
      <Box
        sx={{
          position: "absolute",
          top: 0,
          left: 0,
          right: 0,
          bottom: 0,
          opacity: 0.1,
          background: `
            radial-gradient(circle at 20% 30%, #3b82f6 0%, transparent 30%),
            radial-gradient(circle at 80% 70%, #1d4ed8 0%, transparent 30%),
            radial-gradient(circle at 60% 20%, #2563eb 0%, transparent 25%)
          `,
          animation: "pulse 4s ease-in-out infinite",
        }}
      />

      <Container maxWidth="xl" sx={{ position: "relative", zIndex: 1 }}>
        <Stack
          direction={{ xs: "column", lg: "row" }}
          spacing={6}
          alignItems="center"
          sx={{ py: 8 }}
        >
          {/* Left Content */}
          <Box sx={{ flex: 1, maxWidth: { lg: "50%" } }}>
            <FloatingCard>
              <Stack spacing={3}>
                <Box>
                  <Chip
                    variant="outlined"
                    sx={{
                      borderColor: "primary.300",
                      color: "primary.300",
                      bgcolor: "rgba(59, 130, 246, 0.1)",
                      mb: 2,
                    }}
                    startDecorator={<IoFlash />}
                  >
                    Next Gen CI/CD Engine
                  </Chip>

                  <Typography
                    level="h1"
                    sx={{
                      fontSize: { xs: "2.5rem", md: "3.5rem", lg: "4rem" },
                      fontWeight: 800,
                      background:
                        "linear-gradient(135deg, #ffffff 0%, #3b82f6 100%)",
                      backgroundClip: "text",
                      WebkitBackgroundClip: "text",
                      WebkitTextFillColor: "transparent",
                      lineHeight: 1.1,
                      mb: 2,
                    }}
                  >
                    Conveyor CI
                  </Typography>

                  <Typography
                    level="h2"
                    sx={{
                      fontSize: { xs: "1.5rem", md: "2rem" },
                      fontWeight: 400,
                      color: "#94a3b8",
                      mb: 3,
                      maxWidth: "600px",
                    }}
                  >
                    The lightweight, distributed CI/CD engine built for platform
                    developers who demand simplicity without compromise.
                  </Typography>
                </Box>

                <Stack direction="row" spacing={2} sx={{ flexWrap: "wrap" }}>
                  <Chip
                    sx={{
                      borderColor: "primary.300",
                      color: "primary.300",
                      bgcolor: "rgba(59, 130, 246, 0.1)",
                      mb: 2,
                    }}
                    color="primary"
                    size="lg"
                    startDecorator={<IoFlash />}
                  >
                    Lightning Fast
                  </Chip>
                  <Chip
                    sx={{
                      borderColor: "primary.300",
                      color: "primary.300",
                      bgcolor: "rgba(59, 130, 246, 0.1)",
                      mb: 2,
                    }}
                    color="primary"
                    size="lg"
                    startDecorator={<SiKubernetes />}
                  >
                    Kubernetes-Free
                  </Chip>
                  <Chip
                    sx={{
                      borderColor: "primary.300",
                      color: "primary.300",
                      bgcolor: "rgba(59, 130, 246, 0.1)",
                      mb: 2,
                    }}
                    color="primary"
                    size="lg"
                    startDecorator={<VscGraph />}
                  >
                    Built-in Observability
                  </Chip>
                </Stack>

                <Typography
                  level="body-lg"
                  sx={{
                    color: "#cbd5e1",
                    lineHeight: 1.7,
                    fontSize: "1.1rem",
                    maxWidth: "540px",
                  }}
                >
                  Unlike heavyweight solutions like Jenkins or Kubernetes-bound
                  tools like Tekton, Conveyor CI delivers enterprise-grade CI/CD
                  with minimal overhead and maximum flexibility.
                </Typography>

                <Stack
                  direction={{ xs: "column", sm: "row" }}
                  spacing={2}
                  sx={{ pt: 2 }}
                >
                  <Button
                    size="lg"
                    sx={{
                      bgcolor: "primary.500",
                      "&:hover": { bgcolor: "primary.600" },
                      px: 4,
                      py: 1.5,
                      fontSize: "1.1rem",
                      borderRadius: 8,
                    }}
                    component={Link}
                    href="docs/introduction"
                    endDecorator={<MdOutlineArrowForward />}
                  >
                    Get Started
                  </Button>
                  <Button
                    variant="outlined"
                    size="lg"
                    sx={{
                      borderColor: "primary.300",
                      color: "primary.300",
                      "&:hover": {
                        bgcolor: "rgba(59, 130, 246, 0.1)",
                        borderColor: "primary.400",
                      },
                      px: 4,
                      py: 1.5,
                      fontSize: "1.1rem",
                      borderRadius: 8,
                    }}
                    component={Link}
                    href="https://github.com/open-ug/conveyor"
                    startDecorator={<FaGithub />}
                  >
                    View on GitHub
                  </Button>
                </Stack>

                <Box sx={{ pt: 3 }}>
                  <Typography level="body-sm" sx={{ color: "#64748b", mb: 1 }}>
                    Quick Install:
                  </Typography>
                  <Box
                    sx={{
                      bgcolor: "#0f172a",
                      border: "1px solid",
                      borderColor: "primary.700",
                      borderRadius: 6,
                      p: 2,
                      fontFamily: "Monaco, monospace",
                      color: "#e2e8f0",
                      fontSize: "14px",
                      maxWidth: "400px",
                    }}
                  >
                    $ curl -fsSL conveyor.dev/install | sh
                  </Box>
                </Box>
              </Stack>
            </FloatingCard>
          </Box>

          {/* Right Content - Code Demo */}
          <Box sx={{ flex: 1, maxWidth: { lg: "50%" }, width: "100%" }}>
            <FloatingCard delay={400}>
              <CodeBlock />
            </FloatingCard>
          </Box>
        </Stack>
      </Container>

      <style>{`
        @keyframes pulse {
          0%,
          100% {
            opacity: 0.1;
          }
          50% {
            opacity: 0.2;
          }
        }
      `}</style>
    </Box>
  );
};

export default Hero;
