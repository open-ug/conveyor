import React from "react";

const REPO = "open-ug/conveyor";

const ContributorList = () => {
  const [loading, setLoading] = React.useState(true);
  const [contributors, setContributors] = React.useState<
    Array<{ login: string; avatar_url: string; html_url: string }>
  >([]);

  React.useEffect(() => {
    const fetchContributors = async () => {
      try {
        const response = await fetch(
          `https://api.github.com/repos/${REPO}/contributors`
        );
        const data = await response.json();
        setContributors(data);
      } catch (error) {
        console.error("Error fetching contributors:", error);
      } finally {
        setLoading(false);
      }
    };

    fetchContributors();
  }, []);

  if (loading) {
    return <div>Loading contributors...</div>;
  }
  return (
    <div style={{ display: "flex", flexWrap: "wrap", gap: "16px" }}>
      {contributors.map((contributor) => (
        <a
          key={contributor.login}
          href={contributor.html_url}
          target="_blank"
          rel="noopener noreferrer"
          style={{
            textAlign: "center",
            textDecoration: "none",
            color: "inherit",
          }}
        >
          <img
            src={contributor.avatar_url}
            alt={contributor.login}
            style={{ width: "80px", height: "80px", borderRadius: "50%" }}
          />
          <div>{contributor.login}</div>
        </a>
      ))}
    </div>
  );
};

export default ContributorList;
