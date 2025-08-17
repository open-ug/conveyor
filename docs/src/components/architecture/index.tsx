import React, { useState, useRef, useCallback, useEffect } from "react";
import {
  Box,
  Card,
  CardContent,
  Typography,
  Chip,
  Button,
  Stack,
  Divider,
  Badge,
  LinearProgress,
  Container,
  Grid,
  useTheme,
} from "@mui/joy";
import { GrNetwork } from "react-icons/gr";
import {
  BiGitBranch,
  BiGlobe,
  BiLock,
  BiServer,
  BiShield,
} from "react-icons/bi";
import { FaUsers } from "react-icons/fa";
import { BsActivity, BsDatabase, BsInfo } from "react-icons/bs";
import { CiMonitor } from "react-icons/ci";
import { FiZap } from "react-icons/fi";

const architectureData = {
  nodes: [
    {
      id: "load-balancer",
      name: "Load Balancer",
      type: "network",
      icon: "Network",
      position: { x: 0.5, y: 0.15 }, // Using percentage positions
      status: "healthy",
      description: "Distributes incoming requests across multiple instances",
      metrics: { cpu: 15, memory: 32, requests: 1250 },
      version: "v2.4.1",
      replicas: 2,
    },
    {
      id: "api-gateway",
      name: "API Gateway",
      type: "service",
      icon: "Globe",
      position: { x: 0.5, y: 0.35 },
      status: "healthy",
      description:
        "Central entry point for all API requests with authentication",
      metrics: { cpu: 45, memory: 68, requests: 2100 },
      version: "v3.2.0",
      replicas: 3,
    },
    {
      id: "auth-service",
      name: "Auth Service",
      type: "service",
      icon: "Shield",
      position: { x: 0.25, y: 0.55 },
      status: "healthy",
      description: "Handles authentication and authorization",
      metrics: { cpu: 25, memory: 45, requests: 450 },
      version: "v1.8.2",
      replicas: 2,
    },
    {
      id: "user-service",
      name: "User Service",
      type: "service",
      icon: "Users",
      position: { x: 0.5, y: 0.55 },
      status: "warning",
      description: "Manages user profiles and preferences",
      metrics: { cpu: 78, memory: 85, requests: 890 },
      version: "v2.1.0",
      replicas: 3,
    },
    {
      id: "order-service",
      name: "Order Service",
      type: "service",
      icon: "Container",
      position: { x: 0.75, y: 0.55 },
      status: "healthy",
      description: "Processes and manages orders",
      metrics: { cpu: 35, memory: 52, requests: 650 },
      version: "v1.9.5",
      replicas: 4,
    },
    {
      id: "database",
      name: "Database Cluster",
      type: "storage",
      icon: "Database",
      position: { x: 0.35, y: 0.75 },
      status: "healthy",
      description: "PostgreSQL cluster with read replicas",
      metrics: { cpu: 42, memory: 65, connections: 145 },
      version: "v13.8",
      replicas: 3,
    },
    {
      id: "cache",
      name: "Redis Cache",
      type: "storage",
      icon: "Zap",
      position: { x: 0.65, y: 0.75 },
      status: "healthy",
      description: "In-memory cache for fast data access",
      metrics: { cpu: 18, memory: 38, hits: 2840 },
      version: "v6.2.7",
      replicas: 2,
    },
    {
      id: "monitoring",
      name: "Monitoring",
      type: "observability",
      icon: "Monitor",
      position: { x: 0.15, y: 0.9 },
      status: "healthy",
      description: "Metrics collection and alerting system",
      metrics: { cpu: 22, memory: 45, alerts: 3 },
      version: "v2.35.0",
      replicas: 2,
    },
    {
      id: "logging",
      name: "Log Aggregator",
      type: "observability",
      icon: "Activity",
      position: { x: 0.5, y: 0.9 },
      status: "healthy",
      description: "Centralized logging and log analysis",
      metrics: { cpu: 35, memory: 78, logs: 15000 },
      version: "v7.17.0",
      replicas: 3,
    },
    {
      id: "security",
      name: "Security Scanner",
      type: "security",
      icon: "Lock",
      position: { x: 0.85, y: 0.9 },
      status: "healthy",
      description: "Continuous security scanning and compliance",
      metrics: { cpu: 15, memory: 25, scans: 42 },
      version: "v1.4.3",
      replicas: 1,
    },
  ],
  connections: [
    { from: "load-balancer", to: "api-gateway", type: "http" },
    { from: "api-gateway", to: "auth-service", type: "grpc" },
    { from: "api-gateway", to: "user-service", type: "http" },
    { from: "api-gateway", to: "order-service", type: "http" },
    { from: "user-service", to: "database", type: "sql" },
    { from: "order-service", to: "database", type: "sql" },
    { from: "auth-service", to: "database", type: "sql" },
    { from: "user-service", to: "cache", type: "redis" },
    { from: "order-service", to: "cache", type: "redis" },
    { from: "monitoring", to: "api-gateway", type: "metrics" },
    { from: "monitoring", to: "user-service", type: "metrics" },
    { from: "monitoring", to: "order-service", type: "metrics" },
    { from: "monitoring", to: "database", type: "metrics" },
    { from: "logging", to: "api-gateway", type: "logs" },
    { from: "logging", to: "user-service", type: "logs" },
    { from: "logging", to: "order-service", type: "logs" },
    { from: "security", to: "api-gateway", type: "scan" },
    { from: "security", to: "database", type: "scan" },
  ],
};

const getTypeColor = (type) => {
  switch (type) {
    case "network":
      return "primary";
    case "service":
      return "success";
    case "storage":
      return "warning";
    case "observability":
      return "neutral";
    case "security":
      return "danger";
    default:
      return "neutral";
  }
};

const getStatusColor = (status) => {
  switch (status) {
    case "healthy":
      return "success";
    case "warning":
      return "warning";
    case "error":
      return "danger";
    default:
      return "neutral";
  }
};

const getConnectionColor = (type) => {
  switch (type) {
    case "http":
      return "#3b82f6";
    case "grpc":
      return "#10b981";
    case "sql":
      return "#f59e0b";
    case "redis":
      return "#ef4444";
    case "metrics":
      return "#8b5cf6";
    case "logs":
      return "#06b6d4";
    case "scan":
      return "#f97316";
    default:
      return "#6b7280";
  }
};

function ConnectionLine({
  from,
  to,
  type,
  nodes,
  containerWidth,
  containerHeight,
}) {
  const fromNode = nodes.find((n) => n.id === from);
  const toNode = nodes.find((n) => n.id === to);

  if (!fromNode || !toNode) return null;

  // Convert percentage positions to actual pixels
  const fromX = fromNode.position.x * containerWidth;
  const fromY = fromNode.position.y * containerHeight;
  const toX = toNode.position.x * containerWidth;
  const toY = toNode.position.y * containerHeight;

  const dx = toX - fromX;
  const dy = toY - fromY;
  const angle = (Math.atan2(dy, dx) * 180) / Math.PI;

  // Calculate arrowhead position
  const arrowLength = 10;
  const arrowX = toX - arrowLength * Math.cos(Math.atan2(dy, dx));
  const arrowY = toY - arrowLength * Math.sin(Math.atan2(dy, dx));

  return (
    <g>
      <line
        x1={fromX}
        y1={fromY}
        x2={arrowX}
        y2={arrowY}
        stroke={getConnectionColor(type)}
        strokeWidth="2"
        strokeDasharray={
          type === "metrics" || type === "logs" || type === "scan"
            ? "5,5"
            : "none"
        }
      />
      <polygon
        points={`${toX},${toY} ${toX - 8},${toY - 4} ${toX - 8},${toY + 4}`}
        fill={getConnectionColor(type)}
        transform={`rotate(${angle} ${toX} ${toY})`}
      />
    </g>
  );
}

function ArchitectureNode({
  node,
  isSelected,
  onSelect,
  onDrag,
  isDragging,
  containerWidth,
  containerHeight,
  isMobile,
}) {
  const getIcon = (iconName) => {
    switch (iconName) {
      case "Network":
        return GrNetwork;
      case "Globe":
        return BiGlobe;
      case "Shield":
        return BiShield;
      case "Users":
        return FaUsers;
      case "Container":
        return BiServer;
      case "Database":
        return BsDatabase;
      case "Zap":
        return FiZap;
      case "Monitor":
        return CiMonitor;
      case "Activity":
        return BsActivity;
      case "Lock":
        return BiLock;
      default:
        return BiServer;
    }
  };

  const Icon = getIcon(node.icon);

  // Convert percentage to pixels
  const pixelX = node.position.x * containerWidth;
  const pixelY = node.position.y * containerHeight;

  // Responsive node size
  const nodeSize = isMobile
    ? { width: 80, height: 50 }
    : { width: 100, height: 60 };

  const handleMouseDown = (e) => {
    e.preventDefault();
    const startX = e.clientX - pixelX;
    const startY = e.clientY - pixelY;

    const handleMouseMove = (e) => {
      const newX = Math.max(
        0,
        Math.min(1, (e.clientX - startX) / containerWidth)
      );
      const newY = Math.max(
        0,
        Math.min(1, (e.clientY - startY) / containerHeight)
      );
      onDrag(node.id, { x: newX, y: newY });
    };

    const handleMouseUp = () => {
      document.removeEventListener("mousemove", handleMouseMove);
      document.removeEventListener("mouseup", handleMouseUp);
    };

    document.addEventListener("mousemove", handleMouseMove);
    document.addEventListener("mouseup", handleMouseUp);
  };

  // Touch handling for mobile
  const handleTouchStart = (e) => {
    e.preventDefault();
    const touch = e.touches[0];
    const startX = touch.clientX - pixelX;
    const startY = touch.clientY - pixelY;

    const handleTouchMove = (e) => {
      const touch = e.touches[0];
      const newX = Math.max(
        0,
        Math.min(1, (touch.clientX - startX) / containerWidth)
      );
      const newY = Math.max(
        0,
        Math.min(1, (touch.clientY - startY) / containerHeight)
      );
      onDrag(node.id, { x: newX, y: newY });
    };

    const handleTouchEnd = () => {
      document.removeEventListener("touchmove", handleTouchMove);
      document.removeEventListener("touchend", handleTouchEnd);
    };

    document.addEventListener("touchmove", handleTouchMove);
    document.addEventListener("touchend", handleTouchEnd);
  };

  return (
    <Box
      sx={{
        position: "absolute",
        left: pixelX - nodeSize.width / 2,
        top: pixelY - nodeSize.height / 2,
        width: nodeSize.width,
        height: nodeSize.height,
        cursor: isDragging ? "grabbing" : "grab",
        transform: isSelected ? "scale(1.05)" : "scale(1)",
        transition: isDragging ? "none" : "transform 0.2s ease",
        zIndex: isSelected ? 10 : 1,
        touchAction: "none", // Prevent default touch behaviors
      }}
      onMouseDown={handleMouseDown}
      onTouchStart={handleTouchStart}
      onClick={() => onSelect(node)}
    >
      <Card
        variant={isSelected ? "solid" : "outlined"}
        color={isSelected ? getTypeColor(node.type) : "neutral"}
        sx={{
          width: "100%",
          height: "100%",
          display: "flex",
          alignItems: "center",
          justifyContent: "center",
          position: "relative",
          "&:hover": {
            boxShadow: "lg",
          },
        }}
      >
        <Badge
          badgeContent=""
          color={getStatusColor(node.status)}
          sx={{
            "& .MuiBadge-badge": {
              width: 8,
              height: 8,
              borderRadius: "50%",
              top: 4,
              right: 4,
            },
          }}
        >
          <Stack alignItems="center" spacing={0.5}>
            <Icon size={isMobile ? 16 : 20} />
            <Typography
              level={isMobile ? "body-xs" : "body-sm"}
              textAlign="center"
              sx={{
                fontSize: isMobile ? "10px" : "12px",
                lineHeight: 1.2,
                fontWeight: 500,
              }}
            >
              {node.name}
            </Typography>
          </Stack>
        </Badge>
      </Card>
    </Box>
  );
}

export default function InteractiveArchitecture() {
  const [nodes, setNodes] = useState(architectureData.nodes);
  const [selectedNode, setSelectedNode] = useState(null);
  const [draggingNode, setDraggingNode] = useState(null);
  const [showConnections, setShowConnections] = useState(true);
  const [filterType, setFilterType] = useState("all");
  const [containerDimensions, setContainerDimensions] = useState({
    width: 800,
    height: 600,
  });
  const [isMobile, setIsMobile] = useState(false);
  const containerRef = useRef(null);

  // Handle responsive behavior
  useEffect(() => {
    const handleResize = () => {
      if (containerRef.current) {
        const rect = containerRef.current.getBoundingClientRect();
        const newWidth = rect.width;
        const newHeight =
          window.innerWidth < 768 ? 400 : window.innerWidth < 1024 ? 500 : 600;
        setContainerDimensions({ width: newWidth, height: newHeight });
        setIsMobile(window.innerWidth < 768);
      }
    };

    // Initial call with a small delay to ensure DOM is ready
    setTimeout(handleResize, 100);
    window.addEventListener("resize", handleResize);
    return () => window.removeEventListener("resize", handleResize);
  }, []);

  const handleNodeDrag = useCallback((nodeId, newPosition) => {
    setDraggingNode(nodeId);
    setNodes((prevNodes) =>
      prevNodes.map((node) =>
        node.id === nodeId ? { ...node, position: newPosition } : node
      )
    );
  }, []);

  const handleNodeSelect = useCallback((node) => {
    setSelectedNode(node);
    setDraggingNode(null);
  }, []);

  const filteredNodes = nodes.filter(
    (node) => filterType === "all" || node.type === filterType
  );

  const filteredConnections = architectureData.connections.filter(
    (conn) =>
      filteredNodes.some((node) => node.id === conn.from) &&
      filteredNodes.some((node) => node.id === conn.to)
  );

  const nodeTypes = [...new Set(nodes.map((node) => node.type))];

  return (
    <Container maxWidth="xl" sx={{ py: 3 }}>
      <Stack spacing={3}>
        {/* Header */}
        <Box textAlign="center">
          <Typography level="h1" sx={{ mb: 1 }}>
            Interactive Architecture Diagram
          </Typography>
          <Typography level="body-lg" sx={{ color: "text.secondary" }}>
            Drag nodes to rearrange, click to view details, and explore system
            connections
          </Typography>
        </Box>

        {/* Controls */}
        <Card variant="outlined">
          <CardContent>
            <Stack
              direction={{ xs: "column", sm: "row" }}
              spacing={2}
              alignItems={{ xs: "stretch", sm: "center" }}
              flexWrap="wrap"
              useFlexGap
            >
              <Typography level="body-sm" fontWeight="bold">
                Filters:
              </Typography>
              <Stack direction="row" spacing={1} flexWrap="wrap" useFlexGap>
                <Button
                  variant={filterType === "all" ? "solid" : "outlined"}
                  size="sm"
                  onClick={() => setFilterType("all")}
                >
                  All
                </Button>
                {nodeTypes.map((type) => (
                  <Button
                    key={type}
                    variant={filterType === type ? "solid" : "outlined"}
                    color={getTypeColor(type)}
                    size="sm"
                    onClick={() => setFilterType(type)}
                    sx={{ textTransform: "capitalize" }}
                  >
                    {type}
                  </Button>
                ))}
              </Stack>
              <Divider
                orientation="vertical"
                sx={{ display: { xs: "none", sm: "block" } }}
              />
              <Button
                variant={showConnections ? "solid" : "outlined"}
                size="sm"
                onClick={() => setShowConnections(!showConnections)}
                startDecorator={<BiGitBranch size={16} />}
              >
                {showConnections ? "Hide" : "Show"} Connections
              </Button>
            </Stack>
          </CardContent>
        </Card>

        {/* Main Content */}
        <Stack spacing={3}>
          {/* Architecture Diagram */}
          <Box>
            <Card
              variant="outlined"
              sx={{
                position: "relative",
                height: containerDimensions.height,
                overflow: "hidden",
                bgcolor: "background.level1",
              }}
            >
              <CardContent sx={{ p: 0, height: "100%" }}>
                <Box
                  ref={containerRef}
                  sx={{
                    position: "relative",
                    width: "100%",
                    height: "100%",
                    minHeight: containerDimensions.height,
                  }}
                >
                  {/* SVG for connections */}
                  {showConnections && (
                    <svg
                      style={{
                        position: "absolute",
                        top: 0,
                        left: 0,
                        width: "100%",
                        height: "100%",
                        pointerEvents: "none",
                        zIndex: 0,
                      }}
                    >
                      {filteredConnections.map((connection, index) => (
                        <ConnectionLine
                          key={index}
                          from={connection.from}
                          to={connection.to}
                          type={connection.type}
                          nodes={filteredNodes}
                          containerWidth={containerDimensions.width}
                          containerHeight={containerDimensions.height}
                        />
                      ))}
                    </svg>
                  )}

                  {/* Nodes */}
                  {filteredNodes.map((node) => (
                    <ArchitectureNode
                      key={node.id}
                      node={node}
                      isSelected={selectedNode?.id === node.id}
                      onSelect={handleNodeSelect}
                      onDrag={handleNodeDrag}
                      isDragging={draggingNode === node.id}
                      containerWidth={containerDimensions.width}
                      containerHeight={containerDimensions.height}
                      isMobile={isMobile}
                    />
                  ))}
                </Box>
              </CardContent>
            </Card>
          </Box>

          {/* Details Panel */}
          <Box>
            <Stack spacing={2}>
              {selectedNode ? (
                <Card variant="outlined">
                  <CardContent>
                    <Stack spacing={2}>
                      <Stack
                        direction="row"
                        justifyContent="space-between"
                        alignItems="flex-start"
                        flexWrap="wrap"
                        spacing={1}
                      >
                        <Typography level="h3">{selectedNode.name}</Typography>
                        <Stack direction="row" spacing={1}>
                          <Chip
                            size="sm"
                            color={getTypeColor(selectedNode.type)}
                          >
                            {selectedNode.type}
                          </Chip>
                          <Chip
                            size="sm"
                            color={getStatusColor(selectedNode.status)}
                          >
                            {selectedNode.status}
                          </Chip>
                        </Stack>
                      </Stack>

                      <Typography
                        level="body-sm"
                        sx={{ color: "text.secondary" }}
                      >
                        {selectedNode.description}
                      </Typography>

                      <Divider />

                      <Stack spacing={1}>
                        <Typography level="body-sm" fontWeight="bold">
                          Version & Replicas
                        </Typography>
                        <Stack direction="row" justifyContent="space-between">
                          <Typography level="body-sm">Version:</Typography>
                          <Typography level="body-sm">
                            {selectedNode.version}
                          </Typography>
                        </Stack>
                        <Stack direction="row" justifyContent="space-between">
                          <Typography level="body-sm">Replicas:</Typography>
                          <Typography level="body-sm">
                            {selectedNode.replicas}
                          </Typography>
                        </Stack>
                      </Stack>

                      <Divider />

                      <Stack spacing={2}>
                        <Typography level="body-sm" fontWeight="bold">
                          Metrics
                        </Typography>

                        <Stack spacing={1}>
                          <Stack direction="row" justifyContent="space-between">
                            <Typography level="body-xs">CPU Usage</Typography>
                            <Typography level="body-xs">
                              {selectedNode.metrics.cpu}%
                            </Typography>
                          </Stack>
                          <LinearProgress
                            determinate
                            value={selectedNode.metrics.cpu}
                            color={
                              selectedNode.metrics.cpu > 70
                                ? "danger"
                                : selectedNode.metrics.cpu > 50
                                ? "warning"
                                : "success"
                            }
                            sx={{ height: 6 }}
                          />
                        </Stack>

                        <Stack spacing={1}>
                          <Stack direction="row" justifyContent="space-between">
                            <Typography level="body-xs">
                              Memory Usage
                            </Typography>
                            <Typography level="body-xs">
                              {selectedNode.metrics.memory}%
                            </Typography>
                          </Stack>
                          <LinearProgress
                            determinate
                            value={selectedNode.metrics.memory}
                            color={
                              selectedNode.metrics.memory > 80
                                ? "danger"
                                : selectedNode.metrics.memory > 60
                                ? "warning"
                                : "success"
                            }
                            sx={{ height: 6 }}
                          />
                        </Stack>

                        <Stack direction="row" justifyContent="space-between">
                          <Typography level="body-sm">
                            {selectedNode.type === "storage"
                              ? "Connections"
                              : "Requests/min"}
                            :
                          </Typography>
                          <Typography level="body-sm" fontWeight="bold">
                            {selectedNode.metrics.requests ||
                              selectedNode.metrics.connections ||
                              selectedNode.metrics.hits ||
                              selectedNode.metrics.alerts ||
                              selectedNode.metrics.logs ||
                              selectedNode.metrics.scans}
                          </Typography>
                        </Stack>
                      </Stack>

                      <Divider />

                      <Stack spacing={1}>
                        <Typography level="body-sm" fontWeight="bold">
                          Connections
                        </Typography>
                        {architectureData.connections
                          .filter(
                            (conn) =>
                              conn.from === selectedNode.id ||
                              conn.to === selectedNode.id
                          )
                          .map((conn, idx) => (
                            <Stack
                              key={idx}
                              direction="row"
                              justifyContent="space-between"
                              alignItems="center"
                            >
                              <Typography level="body-xs">
                                {conn.from === selectedNode.id ? "→" : "←"}
                                {conn.from === selectedNode.id
                                  ? nodes.find((n) => n.id === conn.to)?.name
                                  : nodes.find((n) => n.id === conn.from)?.name}
                              </Typography>
                              <Chip
                                size="sm"
                                variant="soft"
                                sx={{ fontSize: "10px" }}
                              >
                                {conn.type}
                              </Chip>
                            </Stack>
                          ))}
                      </Stack>
                    </Stack>
                  </CardContent>
                </Card>
              ) : (
                <Card variant="outlined">
                  <CardContent>
                    <Stack alignItems="center" spacing={2} sx={{ py: 4 }}>
                      <BsInfo size={48} color="#888" />
                      <Typography level="body-lg" textAlign="center">
                        Click on any node to view detailed information
                      </Typography>
                      <Typography
                        level="body-sm"
                        sx={{ color: "text.secondary" }}
                        textAlign="center"
                      >
                        Drag nodes to rearrange the architecture diagram
                      </Typography>
                    </Stack>
                  </CardContent>
                </Card>
              )}

              {/* Legend */}
              <Card variant="outlined">
                <CardContent>
                  <Typography level="h4" sx={{ mb: 2 }}>
                    Connection Types
                  </Typography>
                  <Stack spacing={1}>
                    {[
                      { type: "http", label: "HTTP/REST" },
                      { type: "grpc", label: "gRPC" },
                      { type: "sql", label: "Database" },
                      { type: "redis", label: "Cache" },
                      { type: "metrics", label: "Metrics" },
                      { type: "logs", label: "Logs" },
                      { type: "scan", label: "Security" },
                    ].map(({ type, label }) => (
                      <Stack
                        key={type}
                        direction="row"
                        alignItems="center"
                        spacing={1}
                      >
                        <Box
                          sx={{
                            width: 16,
                            height: 2,
                            backgroundColor: getConnectionColor(type),
                            ...(type === "metrics" ||
                            type === "logs" ||
                            type === "scan"
                              ? {
                                  backgroundImage: `repeating-linear-gradient(90deg, ${getConnectionColor(
                                    type
                                  )} 0px, ${getConnectionColor(
                                    type
                                  )} 4px, transparent 4px, transparent 8px)`,
                                }
                              : {}),
                          }}
                        />
                        <Typography level="body-xs">{label}</Typography>
                      </Stack>
                    ))}
                  </Stack>
                </CardContent>
              </Card>
            </Stack>
          </Box>
        </Stack>
      </Stack>
    </Container>
  );
}
