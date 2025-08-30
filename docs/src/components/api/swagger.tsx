import React, { useEffect, useState } from "react";
import {
  Box,
  Typography,
  Accordion,
  AccordionDetails,
  AccordionSummary,
  AccordionGroup,
  Card,
  CardContent,
  Chip,
  Divider,
  Sheet,
  Badge,
  Button,
  Modal,
  ModalDialog,
  ModalClose,
  Textarea,
  Input,
  FormControl,
  FormLabel,
  Stack,
  List,
  ListItem,
  ListItemContent,
  Table,
  CircularProgress,
} from "@mui/joy";
import { MdSend } from "react-icons/md";
import { Route } from "lucide-react";
import { FaCode } from "react-icons/fa";

const SwaggerUIComponent = ({ swaggerData }) => {
  const [selectedEndpoint, setSelectedEndpoint] = useState(null);
  const [tryItOutOpen, setTryItOutOpen] = useState(false);
  const [requestBody, setRequestBody] = useState("");
  const [pathParams, setPathParams] = useState({});

  // Method color mapping
  const getMethodColor = (method) => {
    const colors = {
      get: "primary",
      post: "success",
      put: "warning",
      delete: "danger",
      patch: "neutral",
    };
    return colors[method.toLowerCase()] || "neutral";
  };

  // Group endpoints by tags
  const groupedEndpoints = React.useMemo(() => {
    if (!swaggerData?.paths) return {};

    const groups = {};
    Object.entries(swaggerData.paths).forEach(([path, methods]) => {
      Object.entries(methods).forEach(([method, details]) => {
        const tag = details.tags?.[0] || "default";
        if (!groups[tag]) groups[tag] = [];
        groups[tag].push({ path, method, ...details });
      });
    });
    return groups;
  }, [swaggerData]);

  // Render parameter table
  const renderParameters = (parameters = []) => (
    <Table size="sm" sx={{ mt: 2 }}>
      <thead>
        <tr>
          <th style={{ width: "20%" }}>Name</th>
          <th style={{ width: "15%" }}>In</th>
          <th style={{ width: "15%" }}>Type</th>
          <th style={{ width: "10%" }}>Required</th>
          <th style={{ width: "40%" }}>Description</th>
        </tr>
      </thead>
      <tbody>
        {parameters.map((param, idx) => (
          <tr key={idx}>
            <td>
              <Typography level="body-sm" fontFamily="monospace">
                {param.name}
              </Typography>
            </td>
            <td>
              <Chip size="sm" variant="soft" color="neutral">
                {param.in}
              </Chip>
            </td>
            <td>
              <Typography level="body-sm">
                {param.type || param.schema?.type || "object"}
              </Typography>
            </td>
            <td>
              {param.required && (
                <Chip size="sm" color="danger" variant="soft">
                  Required
                </Chip>
              )}
            </td>
            <td>
              <Typography level="body-sm">
                {param.description || "-"}
              </Typography>
            </td>
          </tr>
        ))}
      </tbody>
    </Table>
  );

  // Render response codes
  const renderResponses = (responses = {}) => (
    <Stack spacing={2} sx={{ mt: 2 }}>
      {Object.entries(responses).map(([code, response]) => (
        <Card key={code} variant="outlined" size="sm">
          <CardContent>
            <Stack direction="row" spacing={2} alignItems="center">
              <Chip
                color={
                  code.startsWith("2")
                    ? "success"
                    : code.startsWith("4")
                    ? "warning"
                    : "danger"
                }
                variant="soft"
                size="sm"
              >
                {code}
              </Chip>
              <Typography level="body-sm">{response.description}</Typography>
            </Stack>
            {response.schema && (
              <Box
                sx={{
                  mt: 1,
                  p: 1,
                  bgcolor: "background.level1",
                  borderRadius: "sm",
                }}
              >
                <Typography level="body-xs" fontFamily="monospace">
                  Schema: {JSON.stringify(response.schema, null, 2)}
                </Typography>
              </Box>
            )}
          </CardContent>
        </Card>
      ))}
    </Stack>
  );

  // Try it out modal
  const TryItOutModal = () => (
    <Modal open={tryItOutOpen} onClose={() => setTryItOutOpen(false)}>
      <ModalDialog size="lg" sx={{ width: "90vw", maxWidth: 800 }}>
        <ModalClose />
        <Typography level="h4" sx={{ mb: 2 }}>
          Try it out: {selectedEndpoint?.method.toUpperCase()}{" "}
          {selectedEndpoint?.path}
        </Typography>

        <Stack spacing={3}>
          {/* Path Parameters */}
          {selectedEndpoint?.parameters?.filter((p) => p.in === "path").length >
            0 && (
            <Box>
              <Typography level="title-sm" sx={{ mb: 1 }}>
                Path Parameters
              </Typography>
              {selectedEndpoint.parameters
                .filter((p) => p.in === "path")
                .map((param) => (
                  <FormControl key={param.name} sx={{ mb: 1 }}>
                    <FormLabel>
                      {param.name} {param.required && "*"}
                    </FormLabel>
                    <Input
                      placeholder={param.description}
                      onChange={(e) =>
                        setPathParams((prev) => ({
                          ...prev,
                          [param.name]: e.target.value,
                        }))
                      }
                    />
                  </FormControl>
                ))}
            </Box>
          )}

          {/* Request Body */}
          {selectedEndpoint?.method !== "get" && (
            <Box>
              <Typography level="title-sm" sx={{ mb: 1 }}>
                Request Body
              </Typography>
              <Textarea
                minRows={8}
                placeholder="Enter JSON request body..."
                value={requestBody}
                onChange={(e) => setRequestBody(e.target.value)}
                sx={{ fontFamily: "monospace" }}
              />
            </Box>
          )}

          {/* Execute Button */}
          <Button
            startDecorator={<MdSend />}
            color="primary"
            onClick={() => {
              // In a real implementation, you'd make the actual API call here
              alert("This would execute the API call in a real implementation");
            }}
          >
            Execute
          </Button>
        </Stack>
      </ModalDialog>
    </Modal>
  );

  if (!swaggerData) {
    return (
      <Box sx={{ p: 3, textAlign: "center" }}>
        <Typography level="h4">No Swagger data provided</Typography>
      </Box>
    );
  }

  return (
    <Box
      sx={{
        minHeight: "100vh",
        bgcolor: "background.body",
        "--joy-palette-primary-500": "#1976d2",
        "--joy-palette-primary-600": "#1565c0",
        "--joy-palette-primary-700": "#0d47a1",
      }}
    >
      {/* Header */}
      <Sheet
        sx={{
          p: 4,
          bgcolor: "primary.500",
          color: "primary.50",
          background: "linear-gradient(135deg, #1976d2 0%, #1565c0 100%)",
        }}
      >
        <Stack spacing={2}>
          <Stack direction="row" alignItems="center" spacing={2}>
            <Route style={{ fontSize: 32 }} />
            <Typography level="h1" textColor="inherit">
              {swaggerData.info?.title || "API Documentation"}
            </Typography>
          </Stack>
          <Typography level="body-lg" textColor="inherit" sx={{ opacity: 0.9 }}>
            {swaggerData.info?.description}
          </Typography>
          <Stack direction="row" spacing={2}>
            <Chip
              variant="soft"
              color="neutral"
              sx={{ bgcolor: "rgba(255,255,255,0.2)", color: "inherit" }}
            >
              Version {swaggerData.info?.version}
            </Chip>
            <Chip
              variant="soft"
              color="neutral"
              sx={{ bgcolor: "rgba(255,255,255,0.2)", color: "inherit" }}
            >
              {swaggerData.host}
            </Chip>
            <Chip
              variant="soft"
              color="neutral"
              sx={{ bgcolor: "rgba(255,255,255,0.2)", color: "inherit" }}
            >
              {swaggerData.schemes?.join(", ")}
            </Chip>
          </Stack>
        </Stack>
      </Sheet>

      {/* API Documentation */}
      <Box sx={{ p: 3 }}>
        <AccordionGroup>
          {Object.entries(groupedEndpoints).map(([tag, endpoints]) => (
            <Accordion key={tag}>
              <AccordionSummary>
                <Typography level="h3" sx={{ color: "primary.600" }}>
                  {tag.charAt(0).toUpperCase() + tag.slice(1)}
                </Typography>
                <Badge
                  badgeContent={endpoints.length}
                  color="primary"
                  sx={{ ml: 2 }}
                />
              </AccordionSummary>
              <AccordionDetails>
                <Stack spacing={2}>
                  {endpoints.map((endpoint, idx) => (
                    <Card
                      key={idx}
                      variant="outlined"
                      sx={{
                        "&:hover": { boxShadow: "sm" },
                        transition: "box-shadow 0.2s",
                      }}
                    >
                      <CardContent>
                        <Stack spacing={2}>
                          {/* Method and Path Header */}
                          <Stack
                            direction="row"
                            alignItems="center"
                            spacing={2}
                          >
                            <Chip
                              color={getMethodColor(endpoint.method)}
                              variant="solid"
                              size="sm"
                              sx={{
                                minWidth: 80,
                                fontWeight: "bold",
                                textTransform: "uppercase",
                              }}
                            >
                              {endpoint.method}
                            </Chip>
                            <Typography
                              level="title-md"
                              fontFamily="monospace"
                              sx={{
                                flexGrow: 1,
                                wordBreak: "break-all",
                              }}
                            >
                              {endpoint.path}
                            </Typography>
                            <Button
                              size="sm"
                              variant="soft"
                              color="primary"
                              onClick={() => {
                                setSelectedEndpoint(endpoint);
                                setTryItOutOpen(true);
                              }}
                            >
                              Try it out
                            </Button>
                          </Stack>

                          {/* Summary and Description */}
                          <Box>
                            <Typography level="title-sm" sx={{ mb: 1 }}>
                              {endpoint.summary}
                            </Typography>
                            <Typography
                              level="body-sm"
                              sx={{ color: "text.secondary" }}
                            >
                              {endpoint.description}
                            </Typography>
                          </Box>

                          <Divider />

                          <AccordionGroup>
                            {/* Parameters */}
                            {endpoint.parameters?.length > 0 && (
                              <Accordion>
                                <AccordionSummary>
                                  <Typography level="title-sm">
                                    Parameters ({endpoint.parameters.length})
                                  </Typography>
                                </AccordionSummary>
                                <AccordionDetails>
                                  {renderParameters(endpoint.parameters)}
                                </AccordionDetails>
                              </Accordion>
                            )}

                            {/* Responses */}
                            <Accordion>
                              <AccordionSummary>
                                <Typography level="title-sm">
                                  Responses (
                                  {Object.keys(endpoint.responses || {}).length}
                                  )
                                </Typography>
                              </AccordionSummary>
                              <AccordionDetails>
                                {renderResponses(endpoint.responses)}
                              </AccordionDetails>
                            </Accordion>

                            {/* Request Body Schema */}
                            {endpoint.parameters?.find(
                              (p) => p.in === "body"
                            ) && (
                              <Accordion>
                                <AccordionSummary>
                                  <Typography level="title-sm">
                                    Request Body
                                  </Typography>
                                </AccordionSummary>
                                <AccordionDetails>
                                  <Box
                                    sx={{
                                      p: 2,
                                      bgcolor: "background.level1",
                                      borderRadius: "sm",
                                      border: "1px solid",
                                      borderColor: "divider",
                                    }}
                                  >
                                    <Typography
                                      level="body-xs"
                                      fontFamily="monospace"
                                    >
                                      {JSON.stringify(
                                        endpoint.parameters.find(
                                          (p) => p.in === "body"
                                        )?.schema,
                                        null,
                                        2
                                      )}
                                    </Typography>
                                  </Box>
                                </AccordionDetails>
                              </Accordion>
                            )}
                          </AccordionGroup>
                        </Stack>
                      </CardContent>
                    </Card>
                  ))}
                </Stack>
              </AccordionDetails>
            </Accordion>
          ))}
        </AccordionGroup>

        {/* Definitions */}
        {swaggerData.definitions && (
          <Card sx={{ mt: 4 }} variant="outlined">
            <CardContent>
              <Typography level="h2" sx={{ mb: 3, color: "primary.600" }}>
                <FaCode sx={{ mr: 1, verticalAlign: "middle" }} />
                Models
              </Typography>
              <AccordionGroup>
                {Object.entries(swaggerData.definitions).map(
                  ([name, definition]) => (
                    <Accordion key={name}>
                      <AccordionSummary>
                        <Typography level="title-md" fontFamily="monospace">
                          {name}
                        </Typography>
                      </AccordionSummary>
                      <AccordionDetails>
                        <Box
                          sx={{
                            p: 2,
                            bgcolor: "background.level1",
                            borderRadius: "sm",
                            border: "1px solid",
                            borderColor: "divider",
                          }}
                        >
                          <Typography
                            level="body-xs"
                            fontFamily="monospace"
                            sx={{ whiteSpace: "pre-wrap" }}
                          >
                            {JSON.stringify(definition, null, 2)}
                          </Typography>
                        </Box>
                      </AccordionDetails>
                    </Accordion>
                  )
                )}
              </AccordionGroup>
            </CardContent>
          </Card>
        )}
      </Box>

      <TryItOutModal />
    </Box>
  );
};

export default function APIReference() {
  const [swaggerData, setSwaggerData] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    // Fetch Swagger data from the API
    const fetchSwaggerData = async () => {
      setLoading(true);
      try {
        const response = await fetch("/swagger.json");
        const data = await response.json();
        setSwaggerData(data);
      } catch (error) {
        console.error("Error fetching Swagger data:", error);
        setError(error);
      }
      setLoading(false);
    };

    fetchSwaggerData();
  }, []);

  return (
    <>
      {loading ? (
        <CircularProgress />
      ) : error ? (
        <Typography color="danger">Error fetching Swagger data</Typography>
      ) : (
        <SwaggerUIComponent swaggerData={swaggerData} />
      )}
    </>
  );
}
