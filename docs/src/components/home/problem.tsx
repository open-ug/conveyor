import { Chip, Typography, Stack, Button, Box, Container } from "@mui/joy";
import { useState } from "react";
import { FaGlobe } from "react-icons/fa";
import { GiElephant } from "react-icons/gi";
import { GrConfigure, GrDocumentConfig, GrKubernetes } from "react-icons/gr";
import { IoFlash } from "react-icons/io5";
import { SiKubernetes } from "react-icons/si";
import { VscGraph, VscGraphLine } from "react-icons/vsc";

export const ProblemSolution = () => {
  const [activeTab, setActiveTab] = useState(0);
  const [hoveredProblem, setHoveredProblem] = useState(null);

  const problems = [
    {
      icon: (
        <GiElephant
          style={{
            color: "white",
          }}
        />
      ),
      title: "Jenkins: Too Heavy",
      description: "Complex setup, resource-hungry, plugin hell",
      details:
        "Jenkins requires significant infrastructure, complex plugin management, and constant maintenance overhead.",
    },
    {
      icon: (
        <SiKubernetes
          style={{
            color: "#3b82f6",
          }}
        />
      ),
      title: "Tekton/Argo: K8s Lock-in",
      description: "Kubernetes dependency limits flexibility",
      details:
        "Platform teams shouldn't be forced into Kubernetes just for CI/CD. What about edge computing, hybrid clouds, or simpler deployments?",
    },
    {
      icon: (
        <GrDocumentConfig
          style={{
            color: "#3b82f6",
          }}
        />
      ),
      title: "Complex Configuration",
      description: "YAML hell and steep learning curves",
      details:
        "Existing solutions require deep expertise and complex configuration files that are hard to maintain and debug.",
    },
    {
      icon: (
        <VscGraph
          style={{
            color: "whitesmoke",
          }}
        />
      ),
      title: "Poor Observability",
      description: "Limited insights into pipeline performance",
      details:
        "Most tools provide basic logging but lack comprehensive observability, making troubleshooting and optimization difficult.",
    },
  ];

  const solutions = [
    {
      icon: <IoFlash style={{ color: "#10b981" }} />,
      title: "Lightweight Core",
      metric: "< 50MB",
      description: "Minimal resource footprint, fast startup times",
    },
    {
      icon: <FaGlobe style={{ color: "#3b82f6" }} />,
      title: "Platform Agnostic",
      metric: "Any Environment",
      description: "Run anywhere - cloud, on-prem, edge, containers",
    },
    {
      icon: <GrConfigure style={{ color: "#3b82f6" }} />,
      title: "Declarative Simplicity",
      metric: "5 Min Setup",
      description:
        "Intuitive YAML/JSON configuration, minimal learning curve and you can define your own DSL",
    },
    {
      icon: <VscGraphLine style={{ color: "whitesmoke" }} />,
      title: "Built-in Observability",
      metric: "Real-time",
      description: "Comprehensive metrics, logs, and performance insights",
    },
  ];

  return (
    <Box
      sx={{
        py: 10,
        background:
          "linear-gradient(180deg, #0f172a 0%, #1e293b 50%, #0f172a 100%)",
        position: "relative",
      }}
    >
      <Container maxWidth="xl">
        {/* Section Header */}
        <Box sx={{ textAlign: "center", mb: 8 }}>
          <Chip
            variant="outlined"
            sx={{
              borderColor: "primary.300",
              color: "primary.300",
              bgcolor: "rgba(59, 130, 246, 0.1)",
              mb: 3,
            }}
          >
            The Platform Developer's Dilemma
          </Chip>
          <Typography
            level="h2"
            sx={{
              fontSize: { xs: "2rem", md: "3rem" },
              fontWeight: 700,
              color: "white",
              mb: 3,
            }}
          >
            Existing CI/CD Solutions Fall Short
          </Typography>
          <Typography
            level="body-lg"
            sx={{
              color: "#94a3b8",
              maxWidth: "600px",
              mx: "auto",
              fontSize: "1.2rem",
              lineHeight: 1.6,
            }}
          >
            Platform developers need CI/CD tools that are simple, flexible, and
            powerful. Unfortunately, current options force difficult trade-offs.
          </Typography>
        </Box>

        {/* Problems Section */}
        <Box sx={{ mb: 10 }}>
          <Typography
            level="h3"
            sx={{
              fontSize: "1.5rem",
              fontWeight: 600,
              color: "#f87171",
              mb: 4,
              textAlign: "center",
            }}
          >
            😤 The Problems We Face
          </Typography>

          <Box
            sx={{
              display: "grid",
              gridTemplateColumns: { xs: "1fr", md: "repeat(2, 1fr)" },
              gap: 3,
            }}
          >
            {problems.map((problem, index) => (
              <Box
                key={index}
                onMouseEnter={() => setHoveredProblem(index)}
                onMouseLeave={() => setHoveredProblem(null)}
                sx={{
                  p: 4,
                  background:
                    hoveredProblem === index
                      ? "linear-gradient(135deg, rgba(248, 113, 113, 0.1) 0%, rgba(239, 68, 68, 0.05) 100%)"
                      : "rgba(15, 23, 42, 0.8)",
                  border: "1px solid",
                  borderColor: hoveredProblem === index ? "#f87171" : "#334155",
                  borderRadius: 12,
                  cursor: "pointer",
                  transform:
                    hoveredProblem === index
                      ? "translateY(-4px)"
                      : "translateY(0)",
                  transition: "all 0.3s ease",
                  backdropFilter: "blur(10px)",
                }}
              >
                <Box sx={{ display: "flex", alignItems: "flex-start", gap: 3 }}>
                  <Typography sx={{ fontSize: "2rem" }}>
                    {problem.icon}
                  </Typography>
                  <Box sx={{ flex: 1 }}>
                    <Typography
                      level="h4"
                      sx={{
                        color: "white",
                        fontWeight: 600,
                        mb: 1,
                      }}
                    >
                      {problem.title}
                    </Typography>
                    <Typography
                      level="body-md"
                      sx={{
                        color: "#f87171",
                        mb: 2,
                        fontWeight: 500,
                      }}
                    >
                      {problem.description}
                    </Typography>
                    <Typography
                      level="body-sm"
                      sx={{
                        color: "#94a3b8",
                        lineHeight: 1.6,
                      }}
                    >
                      {problem.details}
                    </Typography>
                  </Box>
                </Box>
              </Box>
            ))}
          </Box>
        </Box>

        {/* Solution Divider */}
        <Box
          sx={{
            display: "flex",
            alignItems: "center",
            justifyContent: "center",
            mb: 10,
            position: "relative",
          }}
        >
          <Box
            sx={{
              width: "100%",
              height: "1px",
              background:
                "linear-gradient(90deg, transparent 0%, #3b82f6 50%, transparent 100%)",
            }}
          />
          <Box
            sx={{
              position: "absolute",
              bgcolor: "#1e293b",
              px: 4,
              py: 2,
              borderRadius: 20,
              border: "2px solid #3b82f6",
              color: "#3b82f6",
              fontWeight: 600,
              fontSize: "1.1rem",
            }}
          >
            <Typography
              sx={{
                color: "#3b82f6",
                fontWeight: 600,
                fontSize: "1.1rem",
              }}
              startDecorator={<IoFlash />}
            >
              Enter Conveyor CI
            </Typography>
          </Box>
        </Box>

        {/* Solutions Section */}
        <Box>
          <Typography
            level="h3"
            sx={{
              fontSize: "1.5rem",
              fontWeight: 600,
              color: "#10b981",
              mb: 4,
              textAlign: "center",
            }}
          >
            ✨ How Conveyor CI Solves This
          </Typography>

          <Box
            sx={{
              display: "grid",
              gridTemplateColumns: { xs: "1fr", md: "repeat(2, 1fr)" },
              gap: 3,
              mb: 8,
            }}
          >
            {solutions.map((solution, index) => (
              <Box
                key={index}
                sx={{
                  p: 4,
                  background:
                    "linear-gradient(135deg, rgba(16, 185, 129, 0.1) 0%, rgba(5, 150, 105, 0.05) 100%)",
                  border: "1px solid #10b981",
                  borderRadius: 12,
                  position: "relative",
                  overflow: "hidden",
                  "&:hover": {
                    transform: "translateY(-4px)",
                    boxShadow: "0 20px 40px rgba(16, 185, 129, 0.2)",
                  },
                  transition: "all 0.3s ease",
                }}
              >
                <Box
                  sx={{
                    position: "absolute",
                    top: -2,
                    right: -2,
                    bgcolor: "#10b981",
                    color: "white",
                    px: 2,
                    py: 0.5,
                    borderRadius: "0 12px 0 12px",
                    fontSize: "0.75rem",
                    fontWeight: 600,
                  }}
                >
                  {solution.metric}
                </Box>

                <Box sx={{ display: "flex", alignItems: "flex-start", gap: 3 }}>
                  <Typography sx={{ fontSize: "2rem" }}>
                    {solution.icon}
                  </Typography>
                  <Box sx={{ flex: 1 }}>
                    <Typography
                      level="h4"
                      sx={{
                        color: "white",
                        fontWeight: 600,
                        mb: 2,
                      }}
                    >
                      {solution.title}
                    </Typography>
                    <Typography
                      level="body-md"
                      sx={{
                        color: "#94a3b8",
                        lineHeight: 1.6,
                      }}
                    >
                      {solution.description}
                    </Typography>
                  </Box>
                </Box>
              </Box>
            ))}
          </Box>
        </Box>
      </Container>
    </Box>
  );
};
