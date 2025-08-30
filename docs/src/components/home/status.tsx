import React, { useState, useEffect } from "react";
import {
  Box,
  Typography,
  Button,
  Container,
  Stack,
  Chip,
  Card,
  CardContent,
  Avatar,
  LinearProgress,
} from "@mui/joy";

const AnimatedCounter = ({ target, duration = 2000, suffix = "" }) => {
  const [count, setCount] = useState(0);

  useEffect(() => {
    const startTime = Date.now();
    const animate = () => {
      const elapsed = Date.now() - startTime;
      const progress = Math.min(elapsed / duration, 1);
      const easeOut = 1 - Math.pow(1 - progress, 3);

      setCount(Math.floor(target * easeOut));

      if (progress < 1) {
        requestAnimationFrame(animate);
      }
    };
    animate();
  }, [target, duration]);

  return (
    <Typography
      level="h2"
      sx={{
        fontSize: "2.5rem",
        fontWeight: 800,
        color: "primary.400",
      }}
    >
      {count.toLocaleString()}
      {suffix}
    </Typography>
  );
};

const ContributorAvatar = ({ name, avatar, contributions, delay = 0 }) => {
  const [isVisible, setIsVisible] = useState(false);

  useEffect(() => {
    const timer = setTimeout(() => setIsVisible(true), delay);
    return () => clearTimeout(timer);
  }, [delay]);

  return (
    <Box
      sx={{
        textAlign: "center",
        opacity: isVisible ? 1 : 0,
        transform: isVisible ? "translateY(0)" : "translateY(20px)",
        transition: "all 0.6s ease",
      }}
    >
      <Avatar
        src={avatar}
        sx={{
          width: 60,
          height: 60,
          mx: "auto",
          mb: 1,
          border: "2px solid",
          borderColor: "primary.300",
          "&:hover": {
            transform: "scale(1.1)",
            borderColor: "primary.400",
          },
          transition: "all 0.3s ease",
        }}
      >
        {name.charAt(0)}
      </Avatar>
      <Typography level="body-sm" sx={{ color: "white", fontWeight: 500 }}>
        {name}
      </Typography>
      <Typography level="body-xs" sx={{ color: "#94a3b8" }}>
        {contributions} commits
      </Typography>
    </Box>
  );
};

const RoadmapItem = ({ title, status, description, quarter, features }) => {
  const getStatusColor = (status) => {
    switch (status) {
      case "completed":
        return "success";
      case "in-progress":
        return "warning";
      case "planned":
        return "neutral";
      default:
        return "neutral";
    }
  };

  const getStatusIcon = (status) => {
    switch (status) {
      case "completed":
        return "✅";
      case "in-progress":
        return "🚧";
      case "planned":
        return "📋";
      default:
        return "📋";
    }
  };

  return (
    <Card
      sx={{
        p: 3,
        background:
          status === "completed"
            ? "linear-gradient(135deg, rgba(16, 185, 129, 0.1) 0%, rgba(5, 150, 105, 0.05) 100%)"
            : status === "in-progress"
            ? "linear-gradient(135deg, rgba(245, 158, 11, 0.1) 0%, rgba(217, 119, 6, 0.05) 100%)"
            : "rgba(15, 23, 42, 0.8)",
        border: "1px solid",
        borderColor:
          status === "completed"
            ? "#10b981"
            : status === "in-progress"
            ? "#f59e0b"
            : "#334155",
        borderRadius: 12,
        position: "relative",
        "&:hover": {
          transform: "translateY(-4px)",
          boxShadow: "0 12px 24px rgba(0, 0, 0, 0.2)",
        },
        transition: "all 0.3s ease",
      }}
    >
      <Box
        sx={{
          position: "absolute",
          top: -8,
          left: 16,
          bgcolor:
            getStatusColor(status) === "success"
              ? "#10b981"
              : getStatusColor(status) === "warning"
              ? "#f59e0b"
              : "#64748b",
          color: "white",
          px: 2,
          py: 0.5,
          borderRadius: 8,
          fontSize: "0.75rem",
          fontWeight: 600,
        }}
      >
        {quarter}
      </Box>

      <CardContent>
        <Box sx={{ display: "flex", alignItems: "center", gap: 2, mb: 2 }}>
          <Typography sx={{ fontSize: "1.5rem" }}>
            {getStatusIcon(status)}
          </Typography>
          <Box>
            <Typography
              level="h4"
              sx={{ color: "white", fontWeight: 600, mb: 0.5 }}
            >
              {title}
            </Typography>
            <Chip variant="soft" color={getStatusColor(status)} size="sm">
              {status.replace("-", " ").toUpperCase()}
            </Chip>
          </Box>
        </Box>

        <Typography
          level="body-md"
          sx={{ color: "#94a3b8", mb: 3, lineHeight: 1.6 }}
        >
          {description}
        </Typography>

        <Box>
          <Typography
            level="body-sm"
            sx={{ color: "#cbd5e1", fontWeight: 500, mb: 1 }}
          >
            Key Features:
          </Typography>
          <Stack spacing={0.5}>
            {features.map((feature, index) => (
              <Box
                key={index}
                sx={{ display: "flex", alignItems: "center", gap: 1 }}
              >
                <Box
                  sx={{
                    width: 6,
                    height: 6,
                    borderRadius: "50%",
                    bgcolor:
                      status === "completed"
                        ? "#10b981"
                        : status === "in-progress"
                        ? "#f59e0b"
                        : "#64748b",
                  }}
                />
                <Typography level="body-sm" sx={{ color: "#94a3b8" }}>
                  {feature}
                </Typography>
              </Box>
            ))}
          </Stack>
        </Box>
      </CardContent>
    </Card>
  );
};

const CommunityProjectStatus = () => {
  const [activeSection, setActiveSection] = useState("community");

  const communityStats = [
    { label: "GitHub Stars", value: 12500, suffix: "+" },
    { label: "Contributors", value: 89, suffix: "" },
    { label: "Forks", value: 2100, suffix: "+" },
    { label: "Discord Members", value: 1800, suffix: "+" },
  ];

  const topContributors = [
    { name: "Alex Chen", avatar: "", contributions: 342 },
    { name: "Sarah Kim", avatar: "", contributions: 287 },
    { name: "Marcus Doe", avatar: "", contributions: 219 },
    { name: "Lisa Wang", avatar: "", contributions: 198 },
    { name: "David Brown", avatar: "", contributions: 156 },
    { name: "Emma Wilson", avatar: "", contributions: 134 },
  ];

  const roadmapItems = [
    {
      title: "Core Engine v2.0",
      status: "completed",
      quarter: "Q4 2024",
      description:
        "Complete rewrite of the core engine with improved performance and reliability. New distributed architecture with horizontal scaling capabilities.",
      features: [
        "Distributed task execution",
        "Improved error handling",
        "Plugin API v2",
        "Performance optimizations",
      ],
    },
    {
      title: "Enhanced Observability",
      status: "in-progress",
      quarter: "Q1 2025",
      description:
        "Advanced monitoring, metrics collection, and distributed tracing. Real-time dashboards and alerting system.",
      features: [
        "Distributed tracing",
        "Custom metrics dashboard",
        "Advanced alerting",
        "Performance analytics",
      ],
    },
    {
      title: "Enterprise Features",
      status: "planned",
      quarter: "Q2 2025",
      description:
        "Enterprise-grade features including RBAC, audit logging, and advanced security controls for large organizations.",
      features: [
        "Role-based access control",
        "Audit logging",
        "SAML/OIDC integration",
        "Policy as code",
      ],
    },
    {
      title: "Cloud-Native Extensions",
      status: "planned",
      quarter: "Q3 2025",
      description:
        "Enhanced cloud integrations and serverless deployment options. Native support for major cloud providers.",
      features: [
        "Serverless execution",
        "Cloud provider integrations",
        "Auto-scaling policies",
        "Cost optimization",
      ],
    },
  ];

  return (
    <Box
      sx={{
        py: 10,
        background:
          "linear-gradient(180deg, #1e293b 0%, #0f172a 50%, #1e293b 100%)",
        position: "relative",
      }}
    >
      <Container maxWidth="xl">
        {/* Section Navigation */}
        <Box sx={{ display: "flex", justifyContent: "center", mb: 8 }}>
          <Stack direction="row" spacing={1}>
            <Button
              variant={activeSection === "community" ? "solid" : "outlined"}
              onClick={() => setActiveSection("community")}
              sx={{
                borderColor: "primary.300",
                color: activeSection === "community" ? "white" : "primary.300",
                bgcolor:
                  activeSection === "community" ? "primary.500" : "transparent",
                "&:hover": {
                  bgcolor:
                    activeSection === "community"
                      ? "primary.600"
                      : "rgba(59, 130, 246, 0.1)",
                },
              }}
            >
              👥 Community
            </Button>
            <Button
              variant={activeSection === "status" ? "solid" : "outlined"}
              onClick={() => setActiveSection("status")}
              sx={{
                borderColor: "primary.300",
                color: activeSection === "status" ? "white" : "primary.300",
                bgcolor:
                  activeSection === "status" ? "primary.500" : "transparent",
                "&:hover": {
                  bgcolor:
                    activeSection === "status"
                      ? "primary.600"
                      : "rgba(59, 130, 246, 0.1)",
                },
              }}
            >
              📊 Project Status
            </Button>
          </Stack>
        </Box>

        {/* Community Section */}
        {activeSection === "community" && (
          <Box>
            {/* Community Header */}
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
                👥 Join the Community
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
                Built by Developers, for Developers
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
                Conveyor CI thrives thanks to our amazing open-source community.
                Join thousands of platform developers building the future of
                CI/CD.
              </Typography>
            </Box>

            {/* Community Stats */}
            <Box
              sx={{
                display: "grid",
                gridTemplateColumns: {
                  xs: "repeat(2, 1fr)",
                  md: "repeat(4, 1fr)",
                },
                gap: 4,
                mb: 10,
              }}
            >
              {communityStats.map((stat, index) => (
                <Card
                  key={index}
                  sx={{
                    p: 3,
                    textAlign: "center",
                    background: "rgba(15, 23, 42, 0.8)",
                    border: "1px solid #334155",
                    borderRadius: 12,
                    backdropFilter: "blur(10px)",
                    "&:hover": {
                      transform: "translateY(-4px)",
                      borderColor: "primary.400",
                      boxShadow: "0 8px 32px rgba(59, 130, 246, 0.2)",
                    },
                    transition: "all 0.3s ease",
                  }}
                >
                  <CardContent>
                    <AnimatedCounter target={stat.value} suffix={stat.suffix} />
                    <Typography
                      level="body-sm"
                      sx={{
                        color: "#94a3b8",
                        textTransform: "uppercase",
                        letterSpacing: 1,
                        mt: 0.5,
                      }}
                    >
                      {stat.label}
                    </Typography>
                  </CardContent>
                </Card>
              ))}
            </Box>

            {/* Top Contributors */}
            <Box sx={{ mb: 8 }}>
              <Typography
                level="h3"
                sx={{
                  color: "white",
                  fontWeight: 600,
                  textAlign: "center",
                  mb: 4,
                }}
              >
                🌟 Top Contributors
              </Typography>
              <Box
                sx={{
                  display: "grid",
                  gridTemplateColumns: {
                    xs: "repeat(3, 1fr)",
                    md: "repeat(6, 1fr)",
                  },
                  gap: 3,
                  mb: 4,
                }}
              >
                {topContributors.map((contributor, index) => (
                  <ContributorAvatar
                    key={index}
                    {...contributor}
                    delay={index * 100}
                  />
                ))}
              </Box>
              <Box sx={{ textAlign: "center" }}>
                <Button
                  variant="outlined"
                  sx={{
                    borderColor: "primary.300",
                    color: "primary.300",
                    "&:hover": { bgcolor: "rgba(59, 130, 246, 0.1)" },
                  }}
                >
                  👀 View All Contributors
                </Button>
              </Box>
            </Box>

            {/* Ways to Contribute */}
            <Box
              sx={{
                display: "grid",
                gridTemplateColumns: { xs: "1fr", md: "repeat(3, 1fr)" },
                gap: 4,
                mb: 8,
              }}
            >
              {[
                {
                  icon: "💻",
                  title: "Code Contributions",
                  description: "Submit PRs, fix bugs, add features",
                  action: "Start Contributing",
                  link: "GitHub Repository",
                },
                {
                  icon: "📝",
                  title: "Documentation",
                  description: "Improve docs, write tutorials, examples",
                  action: "Help with Docs",
                  link: "Documentation Site",
                },
                {
                  icon: "🐛",
                  title: "Issue Reports",
                  description: "Report bugs, suggest features",
                  action: "Report Issues",
                  link: "Issue Tracker",
                },
              ].map((item, index) => (
                <Card
                  key={index}
                  sx={{
                    p: 4,
                    textAlign: "center",
                    background:
                      "linear-gradient(135deg, rgba(59, 130, 246, 0.1) 0%, rgba(37, 99, 235, 0.05) 100%)",
                    border: "1px solid #3b82f6",
                    borderRadius: 16,
                    "&:hover": {
                      transform: "translateY(-8px)",
                      boxShadow: "0 20px 40px rgba(59, 130, 246, 0.2)",
                    },
                    transition: "all 0.3s ease",
                  }}
                >
                  <CardContent>
                    <Typography sx={{ fontSize: "3rem", mb: 2 }}>
                      {item.icon}
                    </Typography>
                    <Typography
                      level="h4"
                      sx={{ color: "white", fontWeight: 600, mb: 2 }}
                    >
                      {item.title}
                    </Typography>
                    <Typography
                      level="body-md"
                      sx={{ color: "#94a3b8", mb: 3, lineHeight: 1.6 }}
                    >
                      {item.description}
                    </Typography>
                    <Button
                      variant="outlined"
                      sx={{
                        borderColor: "primary.300",
                        color: "primary.300",
                        "&:hover": { bgcolor: "rgba(59, 130, 246, 0.1)" },
                      }}
                    >
                      {item.action}
                    </Button>
                  </CardContent>
                </Card>
              ))}
            </Box>

            {/* Communication Channels */}
            <Box
              sx={{
                p: 6,
                background: "rgba(15, 23, 42, 0.8)",
                border: "1px solid #334155",
                borderRadius: 16,
                textAlign: "center",
              }}
            >
              <Typography
                level="h3"
                sx={{ color: "white", fontWeight: 600, mb: 3 }}
              >
                💬 Stay Connected
              </Typography>
              <Stack
                direction={{ xs: "column", md: "row" }}
                spacing={3}
                sx={{ justifyContent: "center", mb: 4 }}
              >
                {[
                  { icon: "💬", label: "Discord", members: "1.8K+" },
                  { icon: "🐦", label: "Twitter", followers: "5.2K+" },
                  { icon: "📧", label: "Newsletter", subscribers: "3.1K+" },
                  { icon: "📺", label: "YouTube", subscribers: "892" },
                ].map((channel, index) => (
                  <Box
                    key={index}
                    sx={{
                      p: 2,
                      border: "1px solid #475569",
                      borderRadius: 8,
                      minWidth: 120,
                      "&:hover": {
                        borderColor: "primary.400",
                        bgcolor: "rgba(59, 130, 246, 0.05)",
                      },
                      transition: "all 0.3s ease",
                      cursor: "pointer",
                    }}
                  >
                    <Typography sx={{ fontSize: "1.5rem", mb: 0.5 }}>
                      {channel.icon}
                    </Typography>
                    <Typography level="body-sm" sx={{ color: "white" }}>
                      {channel.label}
                    </Typography>
                    <Typography level="body-xs" sx={{ color: "#94a3b8" }}>
                      {channel.members ||
                        channel.followers ||
                        channel.subscribers}
                    </Typography>
                  </Box>
                ))}
              </Stack>
            </Box>
          </Box>
        )}

        {/* Project Status Section */}
        {activeSection === "status" && (
          <Box>
            {/* Status Header */}
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
                📊 Project Transparency
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
                Roadmap & Project Health
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
                Full transparency on where we are, where we're going, and how
                you can get involved.
              </Typography>
            </Box>

            {/* Current Status Cards */}
            <Box
              sx={{
                display: "grid",
                gridTemplateColumns: { xs: "1fr", md: "repeat(2, 1fr)" },
                gap: 4,
                mb: 8,
              }}
            >
              <Card
                sx={{
                  p: 4,
                  background:
                    "linear-gradient(135deg, rgba(16, 185, 129, 0.1) 0%, rgba(5, 150, 105, 0.05) 100%)",
                  border: "1px solid #10b981",
                  borderRadius: 16,
                }}
              >
                <CardContent>
                  <Box
                    sx={{
                      display: "flex",
                      alignItems: "center",
                      gap: 2,
                      mb: 3,
                    }}
                  >
                    <Typography sx={{ fontSize: "2rem" }}>🏥</Typography>
                    <Typography
                      level="h3"
                      sx={{ color: "white", fontWeight: 600 }}
                    >
                      Project Health
                    </Typography>
                  </Box>
                  <Stack spacing={2}>
                    <Box>
                      <Box
                        sx={{
                          display: "flex",
                          justifyContent: "space-between",
                          mb: 1,
                        }}
                      >
                        <Typography level="body-sm" sx={{ color: "#cbd5e1" }}>
                          Test Coverage
                        </Typography>
                        <Typography level="body-sm" sx={{ color: "#10b981" }}>
                          94%
                        </Typography>
                      </Box>
                      <LinearProgress
                        determinate
                        value={94}
                        sx={{
                          bgcolor: "rgba(16, 185, 129, 0.2)",
                          "& .MuiLinearProgress-bar": { bgcolor: "#10b981" },
                        }}
                      />
                    </Box>
                    <Box>
                      <Box
                        sx={{
                          display: "flex",
                          justifyContent: "space-between",
                          mb: 1,
                        }}
                      >
                        <Typography level="body-sm" sx={{ color: "#cbd5e1" }}>
                          CI/CD Pipeline
                        </Typography>
                        <Typography level="body-sm" sx={{ color: "#10b981" }}>
                          ✅ Passing
                        </Typography>
                      </Box>
                      <LinearProgress
                        determinate
                        value={100}
                        sx={{
                          bgcolor: "rgba(16, 185, 129, 0.2)",
                          "& .MuiLinearProgress-bar": { bgcolor: "#10b981" },
                        }}
                      />
                    </Box>
                    <Box>
                      <Box
                        sx={{
                          display: "flex",
                          justifyContent: "space-between",
                          mb: 1,
                        }}
                      >
                        <Typography level="body-sm" sx={{ color: "#cbd5e1" }}>
                          Security Score
                        </Typography>
                        <Typography level="body-sm" sx={{ color: "#10b981" }}>
                          A+
                        </Typography>
                      </Box>
                      <LinearProgress
                        determinate
                        value={96}
                        sx={{
                          bgcolor: "rgba(16, 185, 129, 0.2)",
                          "& .MuiLinearProgress-bar": { bgcolor: "#10b981" },
                        }}
                      />
                    </Box>
                  </Stack>
                </CardContent>
              </Card>

              <Card
                sx={{
                  p: 4,
                  background:
                    "linear-gradient(135deg, rgba(59, 130, 246, 0.1) 0%, rgba(37, 99, 235, 0.05) 100%)",
                  border: "1px solid #3b82f6",
                  borderRadius: 16,
                }}
              >
                <CardContent>
                  <Box
                    sx={{
                      display: "flex",
                      alignItems: "center",
                      gap: 2,
                      mb: 3,
                    }}
                  >
                    <Typography sx={{ fontSize: "2rem" }}>📜</Typography>
                    <Typography
                      level="h3"
                      sx={{ color: "white", fontWeight: 600 }}
                    >
                      License & Legal
                    </Typography>
                  </Box>
                  <Stack spacing={2}>
                    <Box
                      sx={{ display: "flex", justifyContent: "space-between" }}
                    >
                      <Typography level="body-md" sx={{ color: "#cbd5e1" }}>
                        License
                      </Typography>
                      <Chip variant="soft" color="primary" size="sm">
                        MIT License
                      </Chip>
                    </Box>
                    <Box
                      sx={{ display: "flex", justifyContent: "space-between" }}
                    >
                      <Typography level="body-md" sx={{ color: "#cbd5e1" }}>
                        Copyright
                      </Typography>
                      <Typography level="body-md" sx={{ color: "#94a3b8" }}>
                        2024 Conveyor CI
                      </Typography>
                    </Box>
                    <Box
                      sx={{ display: "flex", justifyContent: "space-between" }}
                    >
                      <Typography level="body-md" sx={{ color: "#cbd5e1" }}>
                        Commercial Use
                      </Typography>
                      <Chip variant="soft" color="success" size="sm">
                        ✅ Allowed
                      </Chip>
                    </Box>
                    <Box
                      sx={{ display: "flex", justifyContent: "space-between" }}
                    >
                      <Typography level="body-md" sx={{ color: "#cbd5e1" }}>
                        Modification
                      </Typography>
                      <Chip variant="soft" color="success" size="sm">
                        ✅ Allowed
                      </Chip>
                    </Box>
                  </Stack>
                </CardContent>
              </Card>
            </Box>

            {/* Roadmap */}
            <Box sx={{ mb: 8 }}>
              <Typography
                level="h3"
                sx={{
                  color: "white",
                  fontWeight: 600,
                  textAlign: "center",
                  mb: 6,
                }}
              >
                🗺️ Development Roadmap
              </Typography>
              <Box
                sx={{
                  display: "grid",
                  gridTemplateColumns: { xs: "1fr", lg: "repeat(2, 1fr)" },
                  gap: 4,
                }}
              >
                {roadmapItems.map((item, index) => (
                  <RoadmapItem key={index} {...item} />
                ))}
              </Box>
            </Box>

            {/* Release Information */}
            <Box
              sx={{
                p: 6,
                background: "rgba(15, 23, 42, 0.8)",
                border: "1px solid #334155",
                borderRadius: 16,
                textAlign: "center",
              }}
            >
              <Typography
                level="h3"
                sx={{ color: "white", fontWeight: 600, mb: 2 }}
              >
                📦 Current Release
              </Typography>
              <Box
                sx={{
                  display: "flex",
                  justifyContent: "center",
                  gap: 4,
                  mb: 4,
                }}
              >
                <Box>
                  <Typography
                    level="h2"
                    sx={{ color: "primary.400", fontWeight: 800 }}
                  >
                    v2.1.0
                  </Typography>
                  <Typography level="body-sm" sx={{ color: "#94a3b8" }}>
                    Latest Stable
                  </Typography>
                </Box>
                <Box>
                  <Typography
                    level="h4"
                    sx={{ color: "#10b981", fontWeight: 700 }}
                  >
                    Dec 2024
                  </Typography>
                  <Typography level="body-sm" sx={{ color: "#94a3b8" }}>
                    Released
                  </Typography>
                </Box>
              </Box>
              <Stack
                direction={{ xs: "column", sm: "row" }}
                spacing={2}
                sx={{ justifyContent: "center" }}
              >
                <Button
                  variant="solid"
                  sx={{
                    bgcolor: "primary.500",
                    "&:hover": { bgcolor: "primary.600" },
                  }}
                >
                  📥 Download v2.1.0
                </Button>
                <Button
                  variant="outlined"
                  sx={{
                    borderColor: "primary.300",
                    color: "primary.300",
                    "&:hover": { bgcolor: "rgba(59, 130, 246, 0.1)" },
                  }}
                >
                  📋 Release Notes
                </Button>
              </Stack>
            </Box>
          </Box>
        )}
      </Container>
    </Box>
  );
};

export default CommunityProjectStatus;
