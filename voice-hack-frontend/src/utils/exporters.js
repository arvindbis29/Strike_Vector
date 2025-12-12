import jsPDF from "jspdf";

export const saveAsJson = (data) => {
  const blob = new Blob([JSON.stringify(data, null, 2)], { type: "application/json" });
  const url = URL.createObjectURL(blob);
  const a = document.createElement("a");
  a.href = url;
  a.download = "insights-response.json";
  a.click();
  URL.revokeObjectURL(url);
};

export const downloadPdf = (data) => {
  const doc = new jsPDF();
  const title = "Insights Report";
  doc.setFontSize(16);
  doc.text(title, 10, 10);

  const text = JSON.stringify(data, null, 2);
  const lines = doc.splitTextToSize(text, 180);
  doc.setFontSize(10);
  doc.text(lines, 10, 20);
  doc.save("insights-report.pdf");
};
