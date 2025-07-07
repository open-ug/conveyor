import React, { useState } from "react";
import {
  Box,
  Card,
  CardContent,
  Typography,
  Chip,
  Button,
  Modal,
  ModalDialog,
  ModalClose,
  Stack,
  LinearProgress,
  Avatar,
  Divider,
  Sheet,
  IconButton,
  Tooltip,
} from "@mui/joy";
import { FaArrowRight, FaCalendar, FaCheckCircle } from "react-icons/fa";
import { FaClock } from "react-icons/fa6";
import { PiGitBranchBold, PiUsers } from "react-icons/pi";
import { FiTarget } from "react-icons/fi";

const roadmapData = [
  {
    id: 1,
    phase: "Q1 2025",
    title: "Core Infrastructure",
    status: "completed",
    progress: 100,
    date: "Jan 2025",
    description: "Foundational components and basic API",
    features: [
      { name: "Container Runtime", status: "completed", dependencies: [] },
      {
        name: "Basic Networking",
        status: "completed",
        dependencies: ["Container Runtime"],
      },
      {
        name: "Storage Interface",
        status: "completed",
        dependencies: ["Container Runtime"],
      },
    ],
    team: "Platform Team",
    impact: "High",
  },
  {
    id: 2,
    phase: "Q2 2025",
    title: "Security & Observability",
    status: "in-progress",
    progress: 65,
    date: "Apr 2025",
    description: "Enhanced security features and monitoring capabilities",
    features: [
      {
        name: "RBAC Implementation",
        status: "completed",
        dependencies: ["Basic Networking"],
      },
      {
        name: "Metrics Collection",
        status: "in-progress",
        dependencies: ["Storage Interface"],
      },
      {
        name: "Security Policies",
        status: "in-progress",
        dependencies: ["RBAC Implementation"],
      },
      {
        name: "Audit Logging",
        status: "planned",
        dependencies: ["Security Policies"],
      },
    ],
    team: "Security Team",
    impact: "High",
  },
  {
    id: 3,
    phase: "Q3 2025",
    title: "Advanced Orchestration",
    status: "in-progress",
    progress: 30,
    date: "Jul 2025",
    description: "Auto-scaling, load balancing, and service mesh integration",
    features: [
      {
        name: "Auto-scaling",
        status: "in-progress",
        dependencies: ["Metrics Collection"],
      },
      {
        name: "Load Balancer",
        status: "planned",
        dependencies: ["Auto-scaling"],
      },
      {
        name: "Service Mesh",
        status: "planned",
        dependencies: ["Security Policies", "Load Balancer"],
      },
      {
        name: "Traffic Management",
        status: "planned",
        dependencies: ["Service Mesh"],
      },
    ],
    team: "Platform Team",
    impact: "Medium",
  },
  {
    id: 4,
    phase: "Q4 2025",
    title: "Enterprise Features",
    status: "planned",
    progress: 0,
    date: "Oct 2025",
    description:
      "Multi-tenancy, advanced analytics, and enterprise integrations",
    features: [
      {
        name: "Multi-tenancy",
        status: "planned",
        dependencies: ["Audit Logging"],
      },
      {
        name: "Advanced Analytics",
        status: "planned",
        dependencies: ["Traffic Management"],
      },
      {
        name: "Enterprise SSO",
        status: "planned",
        dependencies: ["Multi-tenancy"],
      },
      {
        name: "Compliance Tools",
        status: "planned",
        dependencies: ["Advanced Analytics"],
      },
    ],
    team: "Enterprise Team",
    impact: "Medium",
  },
  {
    id: 5,
    phase: "Q1 2026",
    title: "AI/ML Integration",
    status: "planned",
    progress: 0,
    date: "Jan 2026",
    description: "Machine learning workflows and intelligent automation",
    features: [
      {
        name: "ML Pipeline",
        status: "planned",
        dependencies: ["Advanced Analytics"],
      },
      {
        name: "Intelligent Scaling",
        status: "planned",
        dependencies: ["ML Pipeline"],
      },
      {
        name: "Predictive Analytics",
        status: "planned",
        dependencies: ["Intelligent Scaling"],
      },
      {
        name: "AutoML Integration",
        status: "planned",
        dependencies: ["Predictive Analytics"],
      },
    ],
    team: "AI/ML Team",
    impact: "High",
  },
];

const getStatusIcon = (status) => {
  switch (status) {
    case "completed":
      return <FaCheckCircle size={20} style={{ color: "#10b981" }} />;
    case "in-progress":
      return <FaClock size={20} style={{ color: "#f59e0b" }} />;
    default:
      return <FaCalendar size={20} style={{ color: "#6b7280" }} />;
  }
};

const getStatusColor = (status) => {
  switch (status) {
    case "completed":
      return "success";
    case "in-progress":
      return "warning";
    default:
      return "neutral";
  }
};

const getImpactColor = (impact) => {
  switch (impact) {
    case "High":
      return "danger";
    case "Medium":
      return "warning";
    default:
      return "primary";
  }
};

function FeatureNode({ feature, isSelected, onSelect, position }) {
  const nodeStyle = {
    position: "absolute",
    left: `${position.x}px`,
    top: `${position.y}px`,
    transform: "translate(-50%, -50%)",
    transition: "all 0.3s ease",
    cursor: "pointer",
    zIndex: isSelected ? 10 : 1,
  };

  return (
    <Box sx={nodeStyle}>
      <IconButton
        variant={isSelected ? "solid" : "outlined"}
        color={getStatusColor(feature.status)}
        onClick={() => onSelect(feature)}
        sx={{
          width: 40,
          height: 40,
          borderRadius: "50%",
          transform: isSelected ? "scale(1.2)" : "scale(1)",
          boxShadow: isSelected ? "lg" : "sm",
        }}
      >
        {getStatusIcon(feature.status)}
      </IconButton>
      <Typography
        level="body-xs"
        sx={{
          position: "absolute",
          top: "50px",
          left: "50%",
          transform: "translateX(-50%)",
          whiteSpace: "nowrap",
          fontSize: "10px",
          fontWeight: isSelected ? "bold" : "normal",
        }}
      >
        {feature.name}
      </Typography>
    </Box>
  );
}

function ConnectionLine({ from, to }) {
  const dx = to.x - from.x;
  const dy = to.y - from.y;
  const length = Math.sqrt(dx * dx + dy * dy);
  const angle = (Math.atan2(dy, dx) * 180) / Math.PI;

  return (
    <Box
      sx={{
        position: "absolute",
        left: `${from.x}px`,
        top: `${from.y}px`,
        width: `${length}px`,
        height: "2px",
        backgroundColor: "#e5e7eb",
        transformOrigin: "0 50%",
        transform: `rotate(${angle}deg)`,
        zIndex: 0,
      }}
    />
  );
}

export default function CNCFRoadmap() {
  const [selectedMilestone, setSelectedMilestone] = useState(null);
  const [selectedFeature, setSelectedFeature] = useState(null);
  const [expandedPhases, setExpandedPhases] = useState(new Set(["Q2 2025"]));

  const togglePhase = (phase) => {
    const newExpanded = new Set(expandedPhases);
    if (newExpanded.has(phase)) {
      newExpanded.delete(phase);
    } else {
      newExpanded.add(phase);
    }
    setExpandedPhases(newExpanded);
  };

  const getFeaturePositions = (features) => {
    const positions = {};
    const radius = 80;
    const centerX = 200;
    const centerY = 100;

    features.forEach((feature, index) => {
      const angle = (index / features.length) * 2 * Math.PI;
      positions[feature.name] = {
        x: centerX + radius * Math.cos(angle),
        y: centerY + radius * Math.sin(angle),
      };
    });

    return positions;
  };

  const getConnections = (features, positions) => {
    const connections = [];
    features.forEach((feature) => {
      feature.dependencies.forEach((dep) => {
        const depFeature = features.find((f) => f.name === dep);
        if (depFeature && positions[dep] && positions[feature.name]) {
          connections.push({
            from: positions[dep],
            to: positions[feature.name],
          });
        }
      });
    });
    return connections;
  };

  return (
    <Box sx={{ maxWidth: 1200, mx: "auto", p: 3 }}>
      <Typography level="h1" sx={{ mb: 1, textAlign: "center" }}>
        Project Roadmap
      </Typography>
      <Typography
        level="body-lg"
        sx={{ mb: 4, textAlign: "center", color: "text.secondary" }}
      >
        Interactive timeline showing our development milestones and feature
        dependencies
      </Typography>

      <Box sx={{ position: "relative" }}>
        {/* Timeline Line */}
        <Box
          sx={{
            position: "absolute",
            left: "50px",
            top: 0,
            bottom: 0,
            width: "4px",
            backgroundColor: "primary.200",
            borderRadius: "2px",
          }}
        />

        {roadmapData.map((milestone, index) => (
          <Box key={milestone.id} sx={{ position: "relative", mb: 6 }}>
            {/* Timeline Dot */}
            <Box
              sx={{
                position: "absolute",
                left: "38px",
                top: "20px",
                width: "28px",
                height: "28px",
                borderRadius: "50%",
                backgroundColor: "background.surface",
                border: "4px solid",
                borderColor:
                  getStatusColor(milestone.status) === "success"
                    ? "success.400"
                    : getStatusColor(milestone.status) === "warning"
                    ? "warning.400"
                    : "neutral.400",
                display: "flex",
                alignItems: "center",
                justifyContent: "center",
                zIndex: 2,
              }}
            >
              {getStatusIcon(milestone.status)}
            </Box>

            {/* Milestone Card */}
            <Card
              variant="outlined"
              sx={{
                ml: 10,
                cursor: "pointer",
                transition: "all 0.3s ease",
                "&:hover": {
                  boxShadow: "lg",
                  transform: "translateY(-2px)",
                },
              }}
              onClick={() => togglePhase(milestone.phase)}
            >
              <CardContent>
                <Stack
                  direction="row"
                  justifyContent="space-between"
                  alignItems="flex-start"
                  sx={{ mb: 2 }}
                >
                  <Box>
                    <Typography level="h3" sx={{ mb: 1 }}>
                      {milestone.title}
                    </Typography>
                    <Typography
                      level="body-sm"
                      sx={{ color: "text.secondary", mb: 2 }}
                    >
                      {milestone.description}
                    </Typography>
                  </Box>
                  <Stack direction="row" spacing={1} alignItems="center">
                    <Chip
                      size="sm"
                      variant="soft"
                      color={getStatusColor(milestone.status)}
                    >
                      {milestone.status.replace("-", " ")}
                    </Chip>
                    <Chip
                      size="sm"
                      variant="soft"
                      color={getImpactColor(milestone.impact)}
                    >
                      {milestone.impact} Impact
                    </Chip>
                  </Stack>
                </Stack>

                <Stack direction="row" spacing={3} sx={{ mb: 2 }}>
                  <Stack direction="row" spacing={1} alignItems="center">
                    <FaCalendar size={16} />
                    <Typography level="body-sm">{milestone.date}</Typography>
                  </Stack>
                  <Stack direction="row" spacing={1} alignItems="center">
                    <PiUsers size={16} />
                    <Typography level="body-sm">{milestone.team}</Typography>
                  </Stack>
                  <Stack direction="row" spacing={1} alignItems="center">
                    <FiTarget size={16} />
                    <Typography level="body-sm">
                      {milestone.progress}% Complete
                    </Typography>
                  </Stack>
                </Stack>

                <LinearProgress
                  determinate
                  value={milestone.progress}
                  color={getStatusColor(milestone.status)}
                  sx={{ mb: 2 }}
                />

                <Button
                  variant="soft"
                  size="sm"
                  endDecorator={
                    expandedPhases.has(milestone.phase) ? (
                      <FaArrowRight style={{ transform: "rotate(90deg)" }} />
                    ) : (
                      <FaArrowRight />
                    )
                  }
                  onClick={(e) => {
                    e.stopPropagation();
                    togglePhase(milestone.phase);
                  }}
                >
                  {expandedPhases.has(milestone.phase) ? "Hide" : "Show"}{" "}
                  Feature Map
                </Button>
              </CardContent>
            </Card>

            {/* Interactive Feature Map */}
            {expandedPhases.has(milestone.phase) && (
              <Sheet
                variant="soft"
                sx={{
                  ml: 10,
                  mt: 2,
                  p: 3,
                  position: "relative",
                  //height: "300px",
                  borderRadius: "md",
                  overflow: "hidden",
                }}
              >
                <Typography level="h4" sx={{ mb: 2 }}>
                  Feature Dependencies
                </Typography>

                <Box sx={{ position: "relative", height: "200px" }}>
                  {(() => {
                    const positions = getFeaturePositions(milestone.features);
                    const connections = getConnections(
                      milestone.features,
                      positions
                    );

                    return (
                      <>
                        {/* Connection Lines */}
                        {connections.map((connection, idx) => (
                          <ConnectionLine
                            key={idx}
                            from={connection.from}
                            to={connection.to}
                          />
                        ))}

                        {/* Feature Nodes */}
                        {milestone.features.map((feature) => (
                          <FeatureNode
                            key={feature.name}
                            feature={feature}
                            isSelected={selectedFeature?.name === feature.name}
                            onSelect={setSelectedFeature}
                            position={positions[feature.name]}
                          />
                        ))}
                      </>
                    );
                  })()}
                </Box>

                {selectedFeature && (
                  <Card variant="outlined" sx={{ mt: 2 }}>
                    <CardContent>
                      <Stack
                        direction="row"
                        justifyContent="space-between"
                        alignItems="center"
                      >
                        <Typography level="h4">
                          {selectedFeature.name}
                        </Typography>
                        <Chip
                          size="sm"
                          color={getStatusColor(selectedFeature.status)}
                        >
                          {selectedFeature.status}
                        </Chip>
                      </Stack>
                      {selectedFeature.dependencies.length > 0 && (
                        <Typography level="body-sm" sx={{ mt: 1 }}>
                          <strong>Dependencies:</strong>{" "}
                          {selectedFeature.dependencies.join(", ")}
                        </Typography>
                      )}
                    </CardContent>
                  </Card>
                )}
              </Sheet>
            )}
          </Box>
        ))}
      </Box>

      {/* Legend */}
      <Card variant="outlined" sx={{ mt: 4 }}>
        <CardContent>
          <Typography level="h4" sx={{ mb: 2 }}>
            Legend
          </Typography>
          <Stack direction="row" spacing={4} flexWrap="wrap">
            <Stack direction="row" spacing={1} alignItems="center">
              <FaCheckCircle size={16} style={{ color: "#10b981" }} />
              <Typography level="body-sm">Completed</Typography>
            </Stack>
            <Stack direction="row" spacing={1} alignItems="center">
              <FaClock size={16} style={{ color: "#f59e0b" }} />
              <Typography level="body-sm">In Progress</Typography>
            </Stack>
            <Stack direction="row" spacing={1} alignItems="center">
              <FaCalendar size={16} style={{ color: "#6b7280" }} />
              <Typography level="body-sm">Planned</Typography>
            </Stack>
            <Stack direction="row" spacing={1} alignItems="center">
              <PiGitBranchBold size={16} />
              <Typography level="body-sm">
                Click nodes to see dependencies
              </Typography>
            </Stack>
          </Stack>
        </CardContent>
      </Card>
    </Box>
  );
}
