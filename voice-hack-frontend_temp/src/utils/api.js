export const callInsightsAPI = async (jsonString) => {
  try {
    const parsed = JSON.parse(jsonString);

    const res = await fetch("http://localhost:8080/insights/generate/", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(parsed),
    });

    return await res.json();
  } catch (err) {
    return { error: "Invalid JSON or API Error" };
  }
};
