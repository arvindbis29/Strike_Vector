import { CheckCircle, XCircle } from 'lucide-react';
import { InsightResponse } from '../types/insights';
import InsightCard from './InsightCard';

interface Props {
  results: InsightResponse | null;
}

export default function ResultsDisplay({ results }: Props) {
  if (!results) {
    return null;
  }

  const isSuccess = results.code === 200 && results.status === 'Success';

  return (
    <div className="space-y-6">
      <div className={`flex items-center gap-3 p-4 rounded-lg ${isSuccess ? 'bg-green-50 border border-green-200' : 'bg-red-50 border border-red-200'}`}>
        {isSuccess ? (
          <CheckCircle className="w-6 h-6 text-green-600" />
        ) : (
          <XCircle className="w-6 h-6 text-red-600" />
        )}
        <div>
          <p className={`font-semibold ${isSuccess ? 'text-green-900' : 'text-red-900'}`}>
            Status: {results.status}
          </p>
          <p className={`text-sm ${isSuccess ? 'text-green-700' : 'text-red-700'}`}>
            Code: {results.code}
          </p>
          {results.error && (
            <p className="text-sm text-red-700 mt-1">Error: {results.error}</p>
          )}
        </div>
      </div>

      {results.response?.ensights && results.response.ensights.length > 0 && (
        <div className="space-y-4">
          <h2 className="text-2xl font-bold text-gray-900">
            Insights ({results.response.ensights.length})
          </h2>
          <div className="space-y-4">
            {results.response.ensights.map((insight, index) => (
              <InsightCard key={index} insight={insight} />
            ))}
          </div>
        </div>
      )}
    </div>
  );
}
