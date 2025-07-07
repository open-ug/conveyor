import type { ReactNode } from "react";
import clsx from "clsx";
import Link from "@docusaurus/Link";
import useDocusaurusContext from "@docusaurus/useDocusaurusContext";
import Layout from "@theme/Layout";
import Heading from "@theme/Heading";

import styles from "./index.module.css";
import HeroSection from "../components/home/hero";
import FeaturesSection from "../components/home/features";
import CallToActionSection from "../components/home/calltoaction";

export default function Home(): ReactNode {
  const { siteConfig } = useDocusaurusContext();
  return (
    <Layout
      title={`The Software framework for building Cloud Native CI/CD Platforms`}
      description="Conveyor CI is a Software framework for building Custom CI/CD Pipelines"
    >
      <HeroSection />
      <FeaturesSection />
      <CallToActionSection />
    </Layout>
  );
}
