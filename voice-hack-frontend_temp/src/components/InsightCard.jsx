import { Lightbulb, MessageCircle, CheckCircle, AlertTriangle } from "lucide-react";

export default function InsightCard({ item }) {
  return (
    <div className="border rounded-xl p-4 mb-4 bg-gray-50">
      <h3 className="text-lg font-semibold flex items-center gap-2">
        <Lightbulb className="text-yellow-600" /> {item.EnsightType}
      </h3>

      <div className="mt-3 space-y-2 text-sm">
        <p><strong>Concern:</strong> {item.Concerns}</p>
        <p><strong>Resolution:</strong> {item.Resolution}</p>
        <p><strong>Next Steps:</strong> {item.NextSteps}</p>

        {item.Alert && (
          <p className="flex items-center gap-2 text-red-700 font-medium">
            <AlertTriangle size={16} /> {item.Alert}
          </p>
        )}

        {item.Sentiment && (
          <p className="flex items-center gap-2 text-green-600 font-medium">
            <MessageCircle size={16} /> Sentiment: {item.Sentiment}
          </p>
        )}

        {item.KeyPoints && (
          <p className="flex items-center gap-2 text-gray-700">
            <CheckCircle size={16} /> {item.KeyPoints}
          </p>
        )}
      </div>
    </div>
  );
}
