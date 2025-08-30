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
  Tabs,
  TabList,
  Tab,
  TabPanel,
} from "@mui/joy";

const TypewriterText = ({ text, speed = 50, onComplete }) => {
  const [displayText, setDisplayText] = useState("");
  const [currentIndex, setCurrentIndex] = useState(0);

  useEffect(() => {
    if (currentIndex < text.length) {
      const timer = setTimeout(() => {
        setDisplayText(text.slice(0, currentIndex + 1));
        setCurrentIndex(currentIndex + 1);
      }, speed);
      return () => clearTimeout(timer);
    } else if (onComplete) {
      onComplete();
    }
  }, [currentIndex, text, speed, onComplete]);

  return <span>{displayText}</span>;
};

const InstallationCommand = () => {
  const [currentStep, setCurrentStep] = useState(0);
  const [isAnimating, setIsAnimating] = useState(false);

  const installSteps = [
    {
      command: "curl -fsSL conveyor.dev/install | sh",
      description: "Download and install Conveyor CI",
      output: `✓ Downloading Conveyor CI v2.1.0...
✓ Verifying signature...
✓ Installing to /usr/local/bin/conveyor
✓ Installation completed successfully!

Conveyor CI is now ready to use!`,
    },
    {
      command: "conveyor init my-project",
      description: "Initialize a new project with Conveyor CI",
      output: `✓ Creating project structure...
✓ Generated conveyor.yml
✓ Created .conveyor/ directory
✓ Added example workflows

Project initialized! Edit conveyor.yml to get started.`,
    },
    {
      command: "conveyor start --dev",
      description: "Start Conveyor CI in development mode",
      output: `🚀 Starting Conveyor CI Engine...
✓ Engine started on :8080
✓ Worker pool initialized (3 workers)
✓ Event system ready
✓ Dashboard available at http://localhost:8080

Ready to process pipelines!`,
    },
  ];

  const runInstallAnimation = () => {
    setIsAnimating(true);
    setCurrentStep(0);

    const stepInterval = setInterval(() => {
      setCurrentStep((prev) => {
        if (prev < installSteps.length - 1) {
          return prev + 1;
        } else {
          clearInterval(stepInterval);
          setTimeout(() => setIsAnimating(false), 2000);
          return prev;
        }
      });
    }, 4000);
  };

  return (
    <Box
      sx={{
        p: 4,
        background:
          "linear-gradient(135deg, rgba(59, 130, 246, 0.1) 0%, rgba(37, 99, 235, 0.05) 100%)",
        border: "1px solid #3b82f6",
        borderRadius: 16,
        backdropFilter: "blur(10px)",
      }}
    >
      <Box
        sx={{
          display: "flex",
          justifyContent: "space-between",
          alignItems: "center",
          mb: 3,
        }}
      >
        <Typography level="h3" sx={{ color: "white", fontWeight: 600 }}>
          🚀 Quick Installation
        </Typography>
        <Button
          variant="outlined"
          size="sm"
          onClick={runInstallAnimation}
          disabled={isAnimating}
          sx={{
            borderColor: "primary.300",
            color: "primary.300",
            "&:hover": { bgcolor: "rgba(59, 130, 246, 0.1)" },
          }}
        >
          {isAnimating ? "⏳ Running..." : "▶️ Run Demo"}
        </Button>
      </Box>

      <Box
        sx={{
          bgcolor: "#0f172a",
          border: "1px solid #334155",
          borderRadius: 12,
          p: 3,
          fontFamily: 'Monaco, "Lucida Console", monospace',
          fontSize: "14px",
          minHeight: "200px",
        }}
      >
        {installSteps.map((step, index) => (
          <Box
            key={index}
            sx={{
              mb: index < installSteps.length - 1 ? 3 : 0,
              opacity: !isAnimating || currentStep >= index ? 1 : 0.3,
              transition: "opacity 0.5s ease",
            }}
          >
            <Box sx={{ display: "flex", alignItems: "center", mb: 1 }}>
              <Typography sx={{ color: "#10b981", mr: 1 }}>$</Typography>
              <Typography sx={{ color: "#e2e8f0" }}>
                {isAnimating && currentStep === index ? (
                  <TypewriterText text={step.command} speed={30} />
                ) : (
                  step.command
                )}
              </Typography>
            </Box>

            {(!isAnimating || currentStep >= index) && (
              <Box
                sx={{ ml: 2, color: "#94a3b8", whiteSpace: "pre-line", mb: 1 }}
              >
                {step.output}
              </Box>
            )}
          </Box>
        ))}
      </Box>

      <Box sx={{ mt: 3, display: "flex", gap: 2, flexWrap: "wrap" }}>
        <Chip variant="soft" color="success" size="sm">
          ✓ 3 minute setup
        </Chip>
        <Chip variant="soft" color="primary" size="sm">
          🔧 Zero configuration
        </Chip>
        <Chip variant="soft" color="warning" size="sm">
          📦 Single binary
        </Chip>
      </Box>
    </Box>
  );
};

const CodeDemo = () => {
  const [activeTab, setActiveTab] = useState(0);
  const [executionState, setExecutionState] = useState("idle");

  const codeExamples = [
    {
      title: "Basic Pipeline",
      description: "Simple CI/CD pipeline with test and deploy stages",
      filename: "conveyor.yml",
      code: `# conveyor.yml - Basic Pipeline Configuration
version: "2.1"

pipeline:
  name: "web-app-deploy"
  description: "Build, test, and deploy web application"
  
  triggers:
    - push: [main, develop]
    - pull_request: [main]
    
  variables:
    NODE_VERSION: "18"
    APP_NAME: "my-web-app"

stages:
  # Test Stage
  - name: test
    image: node:$NODE_VERSION
    cache:
      - node_modules/
    commands:
      - npm ci
      - npm run lint
      - npm test -- --coverage
      - npm run build
    artifacts:
      - dist/
      - coverage/
    
  # Deploy Stage  
  - name: deploy
    depends: [test]
    image: alpine/kubectl:latest
    when:
      branch: main
    secrets:
      - KUBE_CONFIG
      - AWS_ACCESS_KEY
    commands:
      - kubectl apply -f k8s/deployment.yml
      - kubectl rollout status deployment/$APP_NAME
    notifications:
      slack: "#deployments"`,
      execution: {
        logs: [
          "🚀 Pipeline 'web-app-deploy' started",
          "📦 Pulling image node:18...",
          "⚡ Starting stage 'test'",
          "├── npm ci... ✓",
          "├── npm run lint... ✓",
          "├── npm test... ✓ (98% coverage)",
          "├── npm run build... ✓",
          "📤 Uploading artifacts...",
          "⚡ Starting stage 'deploy'",
          "├── kubectl apply... ✓",
          "├── rollout status... ✓",
          "🎉 Pipeline completed successfully!",
        ],
      },
    },
    {
      title: "Advanced Parallel",
      description: "Complex pipeline with parallel stages and matrix builds",
      filename: "conveyor.yml",
      code: `# conveyor.yml - Advanced Parallel Pipeline
version: "2.1"

pipeline:
  name: "microservices-deploy"
  
  # Matrix build across multiple environments
  matrix:
    environment: [staging, production]
    service: [api, frontend, worker]
    
  # Parallel execution configuration  
  parallel:
    max_concurrent: 6
    fail_fast: false

stages:
  # Parallel testing across services
  - name: test-services
    parallel: true
    matrix:
      service: [api, frontend, worker]
    image: node:18
    working_dir: "/app/services/$service"
    commands:
      - npm ci
      - npm test
      - npm run security-audit
    artifacts:
      - "services/$service/coverage/"
      
  # Integration tests (after all unit tests pass)
  - name: integration-test
    depends: [test-services]
    image: docker/compose:latest
    services:
      - postgres:13
      - redis:6
    commands:
      - docker-compose up -d
      - npm run test:integration
      - docker-compose down
      
  # Parallel deployment per environment
  - name: deploy
    depends: [integration-test]
    parallel: true
    matrix:
      environment: [staging, production]
    image: hashicorp/terraform:latest
    commands:
      - terraform init
      - terraform plan -var="env=$environment"
      - terraform apply -auto-approve
    notifications:
      teams: 
        channel: "#deployments"
        on_success: true
        on_failure: true`,
      execution: {
        logs: [
          "🚀 Matrix pipeline started (6 parallel jobs)",
          "┌── test-services [api] ✓",
          "├── test-services [frontend] ✓",
          "├── test-services [worker] ✓",
          "├── integration-test ⏳",
          "│   ├── Starting services... ✓",
          "│   ├── Running tests... ✓",
          "│   └── Cleanup... ✓",
          "├── deploy [staging] ⏳",
          "├── deploy [production] ⏳",
          "└── All jobs completed! 🎉",
        ],
      },
    },
    {
      title: "Plugin Integration",
      description: "Using plugins for enhanced functionality",
      filename: "conveyor.yml",
      code: `# conveyor.yml - Plugin-Enhanced Pipeline  
version: "2.1"

# Plugin configurations
plugins:
  - name: "slack-notify"
    version: "1.2.0"
    config:
      webhook: "$SLACK_WEBHOOK"
      channel: "#ci-cd"
      
  - name: "security-scanner"
    version: "2.0.1" 
    config:
      severity: "high"
      fail_on_critical: true
      
  - name: "performance-monitor"
    version: "1.0.3"
    config:
      baseline_branch: "main"
      threshold: 10 # 10% regression threshold

pipeline:
  name: "secure-deploy"
  
  # Global plugin hooks
  hooks:
    before_pipeline:
      - slack-notify: "🚀 Starting deployment pipeline"
    after_pipeline:
      - slack-notify: "✅ Pipeline completed"
      - performance-monitor: "analyze"

stages:
  - name: security-scan
    image: alpine:latest
    plugins:
      - security-scanner:
          scan_type: "dependency"
          output: "security-report.json"
    commands:
      - echo "Security scan handled by plugin"
    artifacts:
      - security-report.json
      
  - name: build
    depends: [security-scan]
    image: docker:latest
    commands:
      - docker build -t $APP_NAME:$BUILD_ID .
      - docker push $APP_NAME:$BUILD_ID
    plugins:
      - performance-monitor:
          metric: "build_time"
          
  - name: deploy
    depends: [build]
    image: kubectl:latest
    commands:
      - kubectl set image deployment/$APP_NAME app=$APP_NAME:$BUILD_ID
    plugins:
      - slack-notify:
          message: "🎉 Deployed $APP_NAME:$BUILD_ID to production"
          mention: "@platform-team"`,
      execution: {
        logs: [
          "🔌 Loading plugins...",
          "├── slack-notify v1.2.0 ✓",
          "├── security-scanner v2.0.1 ✓",
          "└── performance-monitor v1.0.3 ✓",
          "🚀 Pipeline started",
          "🔒 Running security scan...",
          "├── Dependency scan: 0 critical issues ✓",
          "🏗️ Building application...",
          "├── Docker build completed ✓",
          "├── Performance: 15% faster than baseline ✓",
          "🚀 Deploying to production...",
          "├── Deployment successful ✓",
          "📢 Slack notification sent ✓",
        ],
      },
    },
  ];

  const runExecution = () => {
    setExecutionState("running");
    setTimeout(() => setExecutionState("completed"), 5000);
  };

  const currentExample = codeExamples[activeTab];

  return (
    <Box>
      <Tabs value={activeTab} onChange={(e, value) => setActiveTab(value)}>
        <TabList
          sx={{
            bgcolor: "rgba(15, 23, 42, 0.8)",
            borderRadius: "12px 12px 0 0",
            border: "1px solid #334155",
            borderBottom: "none",
          }}
        >
          {codeExamples.map((example, index) => (
            <Tab
              key={index}
              sx={{
                color: activeTab === index ? "primary.300" : "#94a3b8",
                "&.Mui-selected": {
                  color: "primary.300",
                  bgcolor: "rgba(59, 130, 246, 0.1)",
                },
              }}
            >
              {example.title}
            </Tab>
          ))}
        </TabList>

        {codeExamples.map((example, index) => (
          <TabPanel key={index} value={index} sx={{ p: 0 }}>
            <Box
              sx={{
                display: "flex",
                flexDirection: { xs: "column", lg: "row" },
                gap: 3,
              }}
            >
              {/* Code Editor */}
              <Box sx={{ flex: 2 }}>
                <Box
                  sx={{
                    bgcolor: "#0f172a",
                    border: "1px solid #334155",
                    borderTop: "none",
                    borderRadius: "0 0 12px 12px",
                    overflow: "hidden",
                  }}
                >
                  {/* Editor Header */}
                  <Box
                    sx={{
                      display: "flex",
                      alignItems: "center",
                      justifyContent: "space-between",
                      p: 2,
                      bgcolor: "#1e293b",
                      borderBottom: "1px solid #334155",
                    }}
                  >
                    <Box sx={{ display: "flex", alignItems: "center", gap: 2 }}>
                      <Box sx={{ display: "flex", gap: 1 }}>
                        <Box
                          sx={{
                            w: 12,
                            h: 12,
                            borderRadius: "50%",
                            bgcolor: "#ef4444",
                          }}
                        />
                        <Box
                          sx={{
                            w: 12,
                            h: 12,
                            borderRadius: "50%",
                            bgcolor: "#f59e0b",
                          }}
                        />
                        <Box
                          sx={{
                            w: 12,
                            h: 12,
                            borderRadius: "50%",
                            bgcolor: "#10b981",
                          }}
                        />
                      </Box>
                      <Typography level="body-sm" sx={{ color: "#94a3b8" }}>
                        {example.filename}
                      </Typography>
                    </Box>
                    <Box sx={{ display: "flex", gap: 2 }}>
                      <Button
                        size="sm"
                        variant="outlined"
                        sx={{
                          borderColor: "primary.700",
                          color: "primary.300",
                          fontSize: "12px",
                          minHeight: "auto",
                          py: 0.5,
                        }}
                      >
                        📋 Copy
                      </Button>
                      <Button
                        size="sm"
                        color="success"
                        onClick={runExecution}
                        disabled={executionState === "running"}
                        sx={{ fontSize: "12px", minHeight: "auto", py: 0.5 }}
                      >
                        {executionState === "running"
                          ? "⏳ Running"
                          : "▶️ Execute"}
                      </Button>
                    </Box>
                  </Box>

                  {/* Code Content */}
                  <Box
                    sx={{
                      p: 3,
                      fontFamily: 'Monaco, "Lucida Console", monospace',
                      fontSize: "13px",
                      color: "#e2e8f0",
                      lineHeight: 1.5,
                      maxHeight: "500px",
                      overflow: "auto",
                    }}
                  >
                    <pre style={{ margin: 0, whiteSpace: "pre-wrap" }}>
                      {example.code}
                    </pre>
                  </Box>
                </Box>
              </Box>

              {/* Execution Panel */}
              <Box sx={{ flex: 1, minWidth: "300px" }}>
                <Card
                  sx={{
                    height: "100%",
                    bgcolor: "#0f172a",
                    border: "1px solid #334155",
                    borderRadius: 12,
                  }}
                >
                  <CardContent>
                    <Box
                      sx={{
                        display: "flex",
                        alignItems: "center",
                        gap: 2,
                        mb: 2,
                      }}
                    >
                      <Typography
                        level="h4"
                        sx={{ color: "white", fontWeight: 600 }}
                      >
                        📊 Live Execution
                      </Typography>
                      <Chip
                        variant="soft"
                        color={
                          executionState === "running"
                            ? "warning"
                            : executionState === "completed"
                            ? "success"
                            : "neutral"
                        }
                        size="sm"
                      >
                        {executionState === "running"
                          ? "⏳ Running"
                          : executionState === "completed"
                          ? "✅ Success"
                          : "⚪ Idle"}
                      </Chip>
                    </Box>

                    <Typography
                      level="body-sm"
                      sx={{ color: "#94a3b8", mb: 3 }}
                    >
                      {example.description}
                    </Typography>

                    <Box
                      sx={{
                        bgcolor: "#1e293b",
                        border: "1px solid #475569",
                        borderRadius: 8,
                        p: 2,
                        fontFamily: "Monaco, monospace",
                        fontSize: "12px",
                        minHeight: "200px",
                        maxHeight: "300px",
                        overflow: "auto",
                      }}
                    >
                      {example.execution.logs.map((log, logIndex) => (
                        <Box
                          key={logIndex}
                          sx={{
                            color: log.includes("✓")
                              ? "#10b981"
                              : log.includes("⏳")
                              ? "#f59e0b"
                              : log.includes("🎉")
                              ? "#3b82f6"
                              : "#e2e8f0",
                            mb: 0.5,
                            opacity:
                              executionState === "running" && logIndex > 3
                                ? 0.5
                                : 1,
                          }}
                        >
                          {log}
                        </Box>
                      ))}

                      {executionState === "running" && (
                        <Box sx={{ color: "#f59e0b", mt: 1 }}>
                          <span>⚡ Executing pipeline...</span>
                          <Box
                            sx={{
                              display: "inline-block",
                              ml: 1,
                              animation: "pulse 1s infinite",
                            }}
                          >
                            ●
                          </Box>
                        </Box>
                      )}
                    </Box>

                    <Box
                      sx={{ mt: 3, display: "flex", gap: 1, flexWrap: "wrap" }}
                    >
                      <Chip
                        variant="outlined"
                        size="sm"
                        sx={{ borderColor: "#475569", color: "#94a3b8" }}
                      >
                        ⚡ Fast execution
                      </Chip>
                      <Chip
                        variant="outlined"
                        size="sm"
                        sx={{ borderColor: "#475569", color: "#94a3b8" }}
                      >
                        📊 Real-time logs
                      </Chip>
                      <Chip
                        variant="outlined"
                        size="sm"
                        sx={{ borderColor: "#475569", color: "#94a3b8" }}
                      >
                        🔄 Live monitoring
                      </Chip>
                    </Box>
                  </CardContent>
                </Card>
              </Box>
            </Box>
          </TabPanel>
        ))}
      </Tabs>
    </Box>
  );
};

const CodeSamplesDemo = () => {
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
            💻 See It In Action
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
            From Installation to Production
          </Typography>
          <Typography
            level="body-lg"
            sx={{
              color: "#94a3b8",
              maxWidth: "700px",
              mx: "auto",
              fontSize: "1.2rem",
              lineHeight: 1.6,
            }}
          >
            Get up and running in minutes. Explore real pipeline configurations
            and see live execution examples.
          </Typography>
        </Box>

        {/* Installation Section */}
        <Box sx={{ mb: 10 }}>
          <InstallationCommand />
        </Box>

        {/* Interactive Code Demo */}
        <Box sx={{ mb: 8 }}>
          <Typography
            level="h3"
            sx={{
              color: "white",
              fontWeight: 600,
              mb: 4,
              textAlign: "center",
            }}
          >
            🚀 Interactive Pipeline Examples
          </Typography>
          <CodeDemo />
        </Box>

        {/* Quick Start CTA */}
        <Box
          sx={{
            textAlign: "center",
            p: 6,
            background:
              "linear-gradient(135deg, rgba(59, 130, 246, 0.1) 0%, rgba(37, 99, 235, 0.05) 100%)",
            border: "1px solid #3b82f6",
            borderRadius: 16,
            backdropFilter: "blur(10px)",
          }}
        >
          <Typography
            level="h4"
            sx={{ color: "white", fontWeight: 600, mb: 2 }}
          >
            Ready to build your first pipeline?
          </Typography>
          <Typography
            level="body-lg"
            sx={{ color: "#94a3b8", mb: 4, maxWidth: "600px", mx: "auto" }}
          >
            Start with our comprehensive documentation and example templates.
            Get your CI/CD pipeline running in under 5 minutes.
          </Typography>
          <Stack
            direction={{ xs: "column", sm: "row" }}
            spacing={2}
            sx={{ justifyContent: "center" }}
          >
            <Button
              size="lg"
              sx={{
                bgcolor: "primary.500",
                "&:hover": { bgcolor: "primary.600" },
                px: 4,
              }}
            >
              📚 Documentation
            </Button>
            <Button
              variant="outlined"
              size="lg"
              sx={{
                borderColor: "primary.300",
                color: "primary.300",
                "&:hover": { bgcolor: "rgba(59, 130, 246, 0.1)" },
                px: 4,
              }}
            >
              🏗️ Example Templates
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

export default CodeSamplesDemo;
