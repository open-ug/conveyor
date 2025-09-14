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
import { ProblemSolution } from "../components/home/problem";
import Features from "../components/home/features";
import CodeSamplesDemo from "../components/home/install";
import CommunityProjectStatus from "../components/home/status";
import Footer from "../components/home/footer";

export default function Home(): ReactNode {
  const { siteConfig } = useDocusaurusContext();
  return (
    <Layout
      title={`Lightweight Engine for building Cloud Native CI/CD Platforms`}
      description="Conveyor CI is a lightweight engine for building distributed CI/CD
                    systems with ease."
    >
      <HeroSection />
      <ProblemSolution />
    </Layout>
  );
}
