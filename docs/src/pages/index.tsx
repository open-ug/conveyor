import type { ReactNode } from "react";
import useDocusaurusContext from "@docusaurus/useDocusaurusContext";
import Layout from "@theme/Layout";

import HeroSection from "../components/home/hero";

export default function Home(): ReactNode {
  const { siteConfig } = useDocusaurusContext();
  return (
    <Layout
      title={`Headless, cloud-native CI/CD orchestration engine`}
      description="Conveyor CI is a headless, cloud-native CI/CD orchestration engine for building distributed CI/CD systems with ease."
    >
      <HeroSection />
    </Layout>
  );
}
