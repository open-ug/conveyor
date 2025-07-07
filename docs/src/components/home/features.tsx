import {
  Box,
  Typography,
  Container,
  Grid,
  Card,
  CardContent,
  useTheme,
} from "@mui/joy";
import {
  FaStream,
  FaNetworkWired,
  FaCogs,
  FaCode,
  FaCloud,
} from "react-icons/fa";

const FeatureCard = ({ icon: Icon, title, description }) => {
  const theme = useTheme();

  return (
    <Card
      variant="outlined"
      sx={{
        height: "100%",
        transition: "all 0.3s",
        "&:hover": {
          boxShadow: "md",
          transform: "translateY(-5px)",
          borderColor: theme.palette.primary[400],
        },
      }}
    >
      <CardContent sx={{ p: 3 }}>
        <Box
          sx={{
            display: "flex",
            alignItems: "center",
            justifyContent: "center",
            width: 60,
            height: 60,
            borderRadius: "12px",
            mb: 2,
            background: `linear-gradient(135deg, ${theme.palette.primary[500]} 0%, ${theme.palette.primary[800]} 100%)`,
            color: "white",
            fontSize: "1.75rem",
          }}
        >
          <Icon />
        </Box>
        <Typography level="title-lg" sx={{ mb: 1.5 }}>
          {title}
        </Typography>
        <Typography level="body-md" sx={{ color: "text.secondary" }}>
          {description}
        </Typography>
      </CardContent>
    </Card>
  );
};

const FeaturesSection = () => {
  const features = [
    {
      icon: FaStream,
      title: "Streamlined Development Paradigm",
      description:
        "Simplify CI/CD development with intuitive workflows that eliminate boilerplate code and reduce complexity, allowing you to focus on building your pipeline logic.",
    },
    {
      icon: FaNetworkWired,
      title: "Distributed Architecture",
      description:
        "Scale effortlessly with a built-in distributed architecture that handles workload distribution, fault tolerance, and high availability out of the box.",
    },
    {
      icon: FaCogs,
      title: "Custom Configuration",
      description:
        "Tailor your CI/CD platform to your exact needs with powerful configuration options that support complex workflows without sacrificing developer experience.",
    },
    {
      icon: FaCode,
      title: "Multi-Language Support",
      description:
        "Build pipelines for any stack with first-class support for multiple programming languages and frameworks, enabling seamless integration across your entire tech ecosystem.",
    },
    {
      icon: FaCloud,
      title: "Cloud Native Support",
      description:
        "Deploy anywhere with native support for Kubernetes, containerization, and major cloud providers, ensuring your CI/CD platform works harmoniously with modern cloud infrastructure.",
    },
  ];

  return (
    <Box sx={{ py: 10 }}>
      <Container>
        <Box sx={{ textAlign: "center", mb: 8 }}>
          <Typography
            level="title-sm"
            sx={{
              color: "primary.500",
              fontWeight: "bold",
              textTransform: "uppercase",
              letterSpacing: "0.05em",
              mb: 1,
            }}
          >
            Features
          </Typography>
          <Typography
            level="h2"
            sx={{
              fontWeight: "bold",
              mb: 2,
              fontSize: { xs: "2rem", md: "2.5rem" },
            }}
          >
            What Makes Conveyor CI Different
          </Typography>
          <Typography
            level="body-lg"
            sx={{
              color: "text.secondary",
              maxWidth: "700px",
              mx: "auto",
            }}
          >
            Conveyor CI eliminates the pain points in building CI/CD platforms
            with features designed for developer productivity and operational
            reliability.
          </Typography>
        </Box>

        <Grid container spacing={4}>
          {features.map((feature, index) => (
            <Grid xs={12} md={6} lg={4} key={index}>
              <FeatureCard
                icon={feature.icon}
                title={feature.title}
                description={feature.description}
              />
            </Grid>
          ))}
        </Grid>
      </Container>
    </Box>
  );
};

export default FeaturesSection;
