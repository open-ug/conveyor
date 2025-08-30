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
} from "@mui/joy";
import { FaCode } from "react-icons/fa";
import { GrNodes } from "react-icons/gr";
import { IoFlash } from "react-icons/io5";
import { PiGraph } from "react-icons/pi";

const AnimatedMetric = ({ value, label, suffix = "", duration = 2000 }) => {
  const [currentValue, setCurrentValue] = useState(0);

  useEffect(() => {
    const targetValue = parseInt(value);
    const startTime = Date.now();

    const animate = () => {
      const elapsed = Date.now() - startTime;
      const progress = Math.min(elapsed / duration, 1);
      const easeOutQuart = 1 - Math.pow(1 - progress, 4);

      setCurrentValue(Math.floor(targetValue * easeOutQuart));

      if (progress < 1) {
        requestAnimationFrame(animate);
      }
    };

    animate();
  }, [value, duration]);

  return (
    <Box sx={{ textAlign: "center" }}>
      <Typography
        level="h2"
        sx={{
          fontSize: "2.5rem",
          fontWeight: 800,
          color: "primary.400",
          mb: 0.5,
        }}
      >
        {currentValue}
        {suffix}
      </Typography>
      <Typography
        level="body-sm"
        sx={{ color: "#94a3b8", textTransform: "uppercase", letterSpacing: 1 }}
      >
        {label}
      </Typography>
    </Box>
  );
};

const FeatureCard = ({ feature, index, isActive, onClick }) => {
  const [isHovered, setIsHovered] = useState(false);

  return (
    <Card
      onClick={onClick}
      onMouseEnter={() => setIsHovered(true)}
      onMouseLeave={() => setIsHovered(false)}
      sx={{
        p: 3,
        cursor: "pointer",
        background: isActive
          ? "linear-gradient(135deg, rgba(59, 130, 246, 0.2) 0%, rgba(37, 99, 235, 0.1) 100%)"
          : "rgba(15, 23, 42, 0.8)",
        border: "1px solid",
        borderColor: isActive ? "primary.400" : "#334155",
        borderRadius: 16,
        transform:
          isHovered || isActive
            ? "translateY(-8px) scale(1.02)"
            : "translateY(0) scale(1)",
        transition: "all 0.4s cubic-bezier(0.4, 0, 0.2, 1)",
        backdropFilter: "blur(10px)",
        boxShadow: isActive
          ? "0 20px 40px rgba(59, 130, 246, 0.3), 0 0 0 1px rgba(59, 130, 246, 0.2)"
          : isHovered
          ? "0 12px 24px rgba(0, 0, 0, 0.3)"
          : "0 4px 8px rgba(0, 0, 0, 0.1)",
        position: "relative",
        overflow: "hidden",
      }}
    >
      {/* Animated background gradient */}
      {isActive && (
        <Box
          sx={{
            position: "absolute",
            top: 0,
            left: -100,
            right: -100,
            bottom: 0,
            background:
              "linear-gradient(45deg, transparent 30%, rgba(59, 130, 246, 0.1) 50%, transparent 70%)",
            animation: "shimmer 3s ease-in-out infinite",
            zIndex: 0,
          }}
        />
      )}

      <CardContent sx={{ position: "relative", zIndex: 1 }}>
        <Box sx={{ display: "flex", alignItems: "center", gap: 2, mb: 2 }}>
          <Box
            sx={{
              fontSize: "2rem",
              p: 2,
              bgcolor: isActive
                ? "rgba(59, 130, 246, 0.2)"
                : "rgba(51, 65, 85, 0.5)",
              borderRadius: 12,
              border: "1px solid",
              borderColor: isActive ? "primary.300" : "#475569",
              transition: "all 0.3s ease",
            }}
          >
            {feature.icon}
          </Box>
          <Box>
            <Typography
              level="h4"
              sx={{
                color: "white",
                fontWeight: 600,
                mb: 0.5,
              }}
            >
              {feature.title}
            </Typography>
            <Chip
              sx={{
                borderColor: "primary.300",
                color: "primary.300",
                bgcolor: "rgba(59, 130, 246, 0.1)",
                mb: 2,
              }}
              color={"primary"}
              size="sm"
            >
              {feature.category}
            </Chip>
          </Box>
        </Box>

        <Typography
          level="body-md"
          sx={{
            color: isActive ? "#cbd5e1" : "#94a3b8",
            lineHeight: 1.6,
            mb: 2,
          }}
        >
          {feature.description}
        </Typography>

        <Stack direction="row" spacing={1} sx={{ flexWrap: "wrap" }}>
          {feature.tags.map((tag, tagIndex) => (
            <Chip
              key={tagIndex}
              variant="outlined"
              size="sm"
              sx={{
                borderColor: "primary.300",
                color: "primary.300",
                bgcolor: "rgba(59, 130, 246, 0.1)",
                mb: 3,
                fontSize: "0.75rem",
              }}
            >
              {tag}
            </Chip>
          ))}
        </Stack>
      </CardContent>
    </Card>
  );
};

const Features = () => {
  const [activeFeature, setActiveFeature] = useState(0);

  const features = [
    {
      icon: <IoFlash style={{ color: "#10b981" }} />,
      title: "Lightweight Architecture",
      category: "Performance",
      description:
        "Built from the ground up for minimal resource consumption. Start pipelines in milliseconds, not minutes. Perfect for resource-constrained environments and edge computing.",
      tags: ["< 50MB", "Fast Startup", "Low Memory"],
      details: {
        metrics: [
          { value: "47", label: "MB Memory", suffix: "" },
          { value: "200", label: "MS Startup", suffix: "" },
          { value: "95", label: "Less Resources", suffix: "%" },
        ],
        codeExample: `# Resource usage comparison
Jenkins:     2GB+ RAM, 10+ seconds startup
Conveyor CI: 47MB RAM, 200ms startup

$ conveyor status
✓ Engine running (47MB RAM)
✓ 5 workers active
✓ Processing 50 pipelines/min`,
      },
    },
    {
      icon: <GrNodes style={{ color: "whitesmoke" }} />,
      title: "Distributed by Design",
      category: "Architecture",
      description:
        "Horizontal scaling with intelligent load balancing. Add workers dynamically across multiple machines, cloud regions, or edge locations without complex configuration.",
      tags: ["Auto-scaling", "Load Balancing", "Multi-region"],
      details: {
        metrics: [
          { value: "1000", label: "Pipelines/Min", suffix: "+" },
          { value: "99.9", label: "Uptime", suffix: "%" },
          { value: "5", label: "Sec Scale-up", suffix: "" },
        ],
        codeExample: `# Scale horizontally in seconds
$ conveyor scale --workers 10 --regions us-east,eu-west

✓ Scaling to 10 workers...
✓ Load balancer updated
✓ Cross-region sync enabled
✓ Ready for 1000+ concurrent pipelines`,
      },
    },
    {
      icon: <PiGraph style={{ color: "whitesmoke" }} />,
      title: "Built-in Observability",
      category: "Monitoring",
      description:
        "Comprehensive metrics, distributed tracing, and real-time monitoring out of the box. No additional tools needed - everything you need to understand your pipelines.",
      tags: ["Real-time Logs", "Metrics", "Tracing", "Dashboards"],
      details: {
        metrics: [
          { value: "100", label: "Metrics Types", suffix: "+" },
          { value: "1", label: "MS Trace Latency", suffix: "" },
          { value: "24", label: "Hours Retention", suffix: "/7" },
        ],
        codeExample: `# Built-in observability dashboard
http://localhost:8080/dashboard

Pipeline Performance:
├── Average execution: 1m 23s
├── Success rate: 98.5%
├── Queue depth: 12 jobs
└── Worker utilization: 76%

Real-time metrics available via API`,
      },
    },
    {
      icon: "🌐",
      title: "Platform Agnostic",
      category: "Flexibility",
      description:
        "Run anywhere - cloud, on-premises, containers, bare metal, or edge devices. No Kubernetes required. Your infrastructure, your choice.",
      tags: ["Cloud Native", "On-Prem", "Edge", "Containers"],
      details: {
        metrics: [
          { value: "15", label: "Platforms", suffix: "+" },
          { value: "0", label: "K8s Dependency", suffix: "" },
          { value: "3", label: "Min Install", suffix: "" },
        ],
        codeExample: `# Deploy anywhere
$ conveyor deploy --target docker
$ conveyor deploy --target kubernetes  
$ conveyor deploy --target bare-metal
$ conveyor deploy --target edge-device

✓ All platforms supported
✓ No vendor lock-in`,
      },
    },
    {
      icon: "🔌",
      title: "Extensible Plugin API",
      category: "Extensibility",
      description:
        "Rich plugin ecosystem with hot-swappable extensions. Create custom integrations without restarting your pipelines. TypeScript/Go SDK included.",
      tags: ["Plugin System", "Hot Reload", "SDK", "Custom Integrations"],
      details: {
        metrics: [
          { value: "50", label: "Built-in Plugins", suffix: "+" },
          { value: "0", label: "Downtime Updates", suffix: "" },
          { value: "2", label: "Supported SDKs", suffix: "" },
        ],
        codeExample: `# Create and deploy plugins
$ conveyor plugin create my-integration
$ conveyor plugin deploy --hot-reload

Available plugins:
├── Slack notifications
├── GitHub integration  
├── AWS deployment
└── Custom metrics collector`,
      },
    },
    {
      icon: "⏱️",
      title: "Real-time Event System",
      category: "Communication",
      description:
        "WebSocket-based event streaming for instant updates. Build reactive dashboards and integrations that respond to pipeline events in real-time.",
      tags: ["WebSockets", "Event Streaming", "Real-time", "Reactive"],
      details: {
        metrics: [
          { value: "1000", label: "Events/Sec", suffix: "+" },
          { value: "10", label: "MS Event Latency", suffix: "" },
          { value: "99.99", label: "Event Delivery", suffix: "%" },
        ],
        codeExample: `# Real-time event streaming
ws://localhost:8080/events

Event stream:
├── pipeline.started
├── stage.completed  
├── deployment.success
└── metrics.updated

Build reactive dashboards instantly!`,
      },
    },
  ];

  useEffect(() => {
    const interval = setInterval(() => {
      setActiveFeature((prev) => (prev + 1) % features.length);
    }, 8000);
    return () => clearInterval(interval);
  }, [features.length]);

  const currentFeature = features[activeFeature];

  return (
    <Box
      sx={{
        py: 10,
        background:
          "linear-gradient(180deg, #1e293b 0%, #0f172a 50%, #1e293b 100%)",
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
          opacity: 0.1,
          background: `
            radial-gradient(circle at 25% 25%, #3b82f6 0%, transparent 50%),
            radial-gradient(circle at 75% 75%, #1d4ed8 0%, transparent 50%)
          `,
        }}
      />

      <Container maxWidth="xl" sx={{ position: "relative", zIndex: 1 }}>
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
            startDecorator={<FaCode />}
          >
            Platform Developer Focused
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
            Features That Actually Matter
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
            Every feature designed with platform developers in mind. No bloat,
            no complexity - just powerful tools that work.
          </Typography>
        </Box>

        {/* Stats Overview */}
        <Box
          sx={{
            display: "grid",
            gridTemplateColumns: { xs: "repeat(2, 1fr)", md: "repeat(4, 1fr)" },
            gap: 4,
            mb: 8,
            p: 4,
            background: "rgba(15, 23, 42, 0.8)",
            borderRadius: 16,
            border: "1px solid #334155",
            backdropFilter: "blur(10px)",
          }}
        >
          <AnimatedMetric value="47" label="MB Memory" />
          <AnimatedMetric value="200" label="MS Startup" />
          <AnimatedMetric value="1000" label="Pipelines/Min" suffix="+" />
          <AnimatedMetric value="99.9" label="Uptime" suffix="%" />
        </Box>

        {/* Interactive Features Grid */}
        <Box
          sx={{
            display: "flex",
            flexDirection: { xs: "column", lg: "row" },
            gap: 4,
          }}
        >
          {/* Feature Cards */}
          <Box sx={{ flex: 1 }}>
            <Box
              sx={{
                display: "grid",
                gridTemplateColumns: { xs: "1fr", md: "repeat(2, 1fr)" },
                gap: 3,
              }}
            >
              {features.map((feature, index) => (
                <FeatureCard
                  key={index}
                  feature={feature}
                  index={index}
                  isActive={activeFeature === index}
                  onClick={() => setActiveFeature(index)}
                />
              ))}
            </Box>
          </Box>

          {/* Feature Details Panel */}
          <Box sx={{ flex: 1, maxWidth: { lg: "500px" } }}>
            <Card
              sx={{
                p: 4,
                background:
                  "linear-gradient(135deg, rgba(59, 130, 246, 0.1) 0%, rgba(37, 99, 235, 0.05) 100%)",
                border: "1px solid #3b82f6",
                borderRadius: 16,
                backdropFilter: "blur(10px)",
                height: "fit-content",
                position: "sticky",
                top: 20,
              }}
            >
              <CardContent>
                <Box
                  sx={{ display: "flex", alignItems: "center", gap: 2, mb: 3 }}
                >
                  <Typography sx={{ fontSize: "2.5rem" }}>
                    {currentFeature.icon}
                  </Typography>
                  <Box>
                    <Typography
                      level="h3"
                      sx={{ color: "white", fontWeight: 600, mb: 0.5 }}
                    >
                      {currentFeature.title}
                    </Typography>
                    <Chip
                      sx={{
                        borderColor: "primary.300",
                        color: "primary.300",
                        bgcolor: "rgba(59, 130, 246, 0.1)",
                        mb: 2,
                      }}
                      color="primary"
                      size="sm"
                    >
                      {currentFeature.category}
                    </Chip>
                  </Box>
                </Box>

                {/* Metrics */}
                <Box
                  sx={{
                    display: "grid",
                    gridTemplateColumns: "repeat(3, 1fr)",
                    gap: 2,
                    mb: 4,
                  }}
                >
                  {currentFeature.details.metrics.map((metric, index) => (
                    <Box
                      key={index}
                      sx={{
                        textAlign: "center",
                        p: 2,
                        bgcolor: "rgba(59, 130, 246, 0.1)",
                        borderRadius: 8,
                        border: "1px solid rgba(59, 130, 246, 0.2)",
                      }}
                    >
                      <Typography
                        level="h4"
                        sx={{ color: "primary.300", fontWeight: 700 }}
                      >
                        {metric.value}
                        {metric.suffix}
                      </Typography>
                      <Typography
                        level="body-xs"
                        sx={{ color: "#94a3b8", textTransform: "uppercase" }}
                      >
                        {metric.label}
                      </Typography>
                    </Box>
                  ))}
                </Box>

                {/* Code Example */}
                <Box
                  sx={{
                    bgcolor: "#0f172a",
                    border: "1px solid #334155",
                    borderRadius: 8,
                    p: 3,
                    fontFamily: "Monaco, monospace",
                    fontSize: "13px",
                    color: "#e2e8f0",
                    lineHeight: 1.5,
                    overflow: "auto",
                    maxHeight: "200px",
                    margin: 0,
                    whiteSpace: "pre-wrap",
                  }}
                  component="pre"
                >
                  {currentFeature.details.codeExample}
                </Box>
              </CardContent>
            </Card>
          </Box>
        </Box>

        {/* Bottom CTA */}
        <Box sx={{ textAlign: "center", mt: 10 }}>
          <Typography level="h3" sx={{ color: "white", mb: 2 }}>
            Ready to experience these features?
          </Typography>
          <Button
            size="lg"
            sx={{
              bgcolor: "primary.500",
              "&:hover": { bgcolor: "primary.600" },
              px: 6,
              py: 1.5,
            }}
          >
            Get Started with Conveyor CI →
          </Button>
        </Box>
      </Container>

      <style>{`
        @keyframes shimmer {
          0% {
            transform: translateX(-100%);
          }
          100% {
            transform: translateX(100%);
          }
        }
      `}</style>
    </Box>
  );
};

export default Features;
