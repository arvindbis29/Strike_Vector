import InsightCard from "./InsightCard";
import { ArrowLeft } from "lucide-react";

export default function InsightsView({ responseData, goBack }) {
  return (
    <div className="bg-white shadow-md rounded-xl p-6">
      <button
        onClick={goBack}
        className="flex items-center gap-2 text-blue-600 hover:underline mb-4"
      >
        <ArrowLeft size={18} /> Back
      </button>

      <h2 className="text-2xl font-bold mb-4">Insights</h2>

      {responseData.response?.ensights?.map((ins, idx) => (
        <InsightCard key={idx} item={ins} />
      ))}
    </div>
  );
}
