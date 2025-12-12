import InsightCard from "./InsightCard";
import { saveAsJson, downloadPdf } from "../utils/exporters";

export default function InsightsDisplay({ responseData, loading, fullWidth }) {
  const hasData = responseData && responseData.response && responseData.response.ensights;

  return (
    <div>
      <div className="flex items-center justify-between mb-3">
        <h3 className="font-semibold">Insights</h3>

        <div className="flex gap-2">
          <button onClick={()=>saveAsJson(responseData)} className="px-3 py-1 border rounded">Export JSON</button>
          <button onClick={()=>downloadPdf(responseData)} className="px-3 py-1 border rounded">Download PDF</button>
        </div>
      </div>

      {loading && <div className="text-gray-500">Loading...</div>}

      {!loading && !hasData && (
        <div className="text-gray-500">No insights yet. Click "Generate Insights" after providing input JSON.</div>
      )}

      {hasData && (
        <div className={fullWidth ? "space-y-3" : "grid grid-cols-1 gap-4"}>
          {responseData.response.ensights.map((it, idx) => (
            <InsightCard key={idx} item={it} index={idx} />
          ))}
        </div>
      )}
    </div>
  );
}
