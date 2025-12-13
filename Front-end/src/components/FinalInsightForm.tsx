import { useState } from 'react';
import { Search } from 'lucide-react';
import { FinalInsightRequest, GenerateInsightRequest, CallData, InsightResponse } from '../types/insights';

interface Props {
  onResults: (results: InsightResponse) => void;
  onLoading: (loading: boolean) => void;
}

export default function FinalInsightForm({ onResults, onLoading }: Props) {
  const [maxCallLimit, setMaxCallLimit] = useState<number>(100);
  const [error, setError] = useState<string>('');
  // Add the same fields the generate endpoint expects so the final endpoint
  // receives a complete payload. These have sensible defaults matching
  // `GenerateInsightForm` so Final Insights can be requested independently.
  const [glid] = useState<number>(18888);
  const [executiveId] = useState<string>('78910');
  const [customerType] = useState<string>('New');
  const [customerCityName] = useState<string>('New Delhi');
  const [callData] = useState<CallData[]>([
    {
      call_recording_url: '',
      call_type: 'PNS',
      call_date: new Date().toISOString().slice(0, 16).replace('T', ' ') + ':00',
    },
  ]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    onLoading(true);

    try {
      // Send a full payload that satisfies the server's required fields.
      const requestData: GenerateInsightRequest & FinalInsightRequest = {
        glid,
        executive_id: executiveId,
        customer_type: customerType,
        customer_city_name: customerCityName,
        max_call_limit: maxCallLimit,
        call_data: callData,
      };

      const response = await fetch('http://localhost:8080/insights/final/', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(requestData),
      });

      if (!response.ok) {
        const text = await response.text();
        throw new Error(`API Error: ${response.status} ${text}`);
      }

      const data: any = await response.json();

      // Normalize like generate form: make sure `response.ensights` exists and
      // insight objects include `EnsightType` expected by the UI.
      const rawInsights: any[] = data?.response?.ensights || data?.response?.Insights || data?.response?.insights || [];

      const normalizedInsights = rawInsights.map((insight: any) => ({
        ...insight,
        EnsightType: insight.EnsightType || insight.InsightType || '',
        Concerns: insight.Concerns || insight.concerns || '',
        Resolution: insight.Resolution || insight.resolution || '',
        NextSteps: insight.NextSteps || insight.nextSteps || '',
        Alert: insight.Alert || insight.alert || '',
        Sentiment: insight.Sentiment || insight.sentiment || '',
        KeyPoints: insight.KeyPoints || insight.keyPoints || '',
      }));

      const normalized: any = {
        ...data,
        response: { ensights: normalizedInsights },
      };

      onResults(normalized as any);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch insights');
    } finally {
      onLoading(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-2">
          Max Call Limit
        </label>
        <input
          type="number"
          value={maxCallLimit}
          onChange={(e) => setMaxCallLimit(Number(e.target.value))}
          className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
          min="1"
          required
        />
      </div>

      {error && (
        <div className="p-4 bg-red-50 border border-red-200 rounded-lg text-red-700 text-sm">
          {error}
        </div>
      )}

      <button
        type="submit"
        className="w-full bg-blue-600 text-white py-3 px-4 rounded-lg hover:bg-blue-700 transition-colors flex items-center justify-center gap-2 font-medium"
      >
        <Search className="w-5 h-5" />
        Get Final Insights
      </button>
    </form>
  );
}
