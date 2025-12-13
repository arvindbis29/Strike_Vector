import { useState } from 'react';
import { Search, Plus, Trash2 } from 'lucide-react';
import { GenerateInsightRequest, CallData, InsightResponse } from '../types/insights';

interface Props {
  onResults: (results: InsightResponse) => void;
  onLoading: (loading: boolean) => void;
}

export default function GenerateInsightForm({ onResults, onLoading }: Props) {
  const [glid, setGlid] = useState<number>(18888);
  const [executiveId, setExecutiveId] = useState<string>('78910');
  const [customerType, setCustomerType] = useState<string>('Star');
  const [customerCityName, setCustomerCityName] = useState<string>('New Delhi');
  const [maxCallLimit, setMaxCallLimit] = useState<number>(10);
  const [callData, setCallData] = useState<CallData[]>([
    {
      call_recording_url: '',
      call_type: 'PNS',
      call_date: new Date().toISOString().slice(0, 16).replace('T', ' ') + ':00',
    },
  ]);
  const [error, setError] = useState<string>('');

  const handleAddCall = () => {
    setCallData([
      ...callData,
      {
        call_recording_url: '',
        call_type: 'PNS',
        call_date: new Date().toISOString().slice(0, 16).replace('T', ' ') + ':00',
      },
    ]);
  };

  const handleRemoveCall = (index: number) => {
    setCallData(callData.filter((_, i) => i !== index));
  };

  const handleCallDataChange = (index: number, field: keyof CallData, value: string) => {
    const newCallData = [...callData];
    newCallData[index] = { ...newCallData[index], [field]: value };
    setCallData(newCallData);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    onLoading(true);

    try {
      const requestData: GenerateInsightRequest = {
        glid,
        executive_id: executiveId,
        customer_type: customerType,
        customer_city_name: customerCityName,
        max_call_limit: maxCallLimit,
        call_data: callData,
      };

      const response = await fetch('http://localhost:8080/insights/generate/', {
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

      // Normalize response shape: backend sometimes returns `Insights` (capitalized)
      // or `ensights`. Ensure `response.ensights` exists and insight objects have
      // the `EnsightType` key expected by the UI.
      const rawInsights: any[] = data?.response?.ensights || data?.response?.Insights || data?.response?.insights || [];

      const normalizedInsights = rawInsights.map((insight: any) => ({
        // Keep existing keys, but ensure UI-friendly keys exist
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
      setError(err instanceof Error ? err.message : 'Failed to generate insights');
    } finally {
      onLoading(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        {/* <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">GLID</label>
          <input
            type="number"
            value={glid}
            onChange={(e) => setGlid(Number(e.target.value))}
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            required
          />
        </div> */}

        {/* <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">
            Executive ID
          </label>
          <input
            type="text"
            value={executiveId}
            onChange={(e) => setExecutiveId(e.target.value)}
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            required
          />
        </div> */}

        {/* <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">
            Customer Type
          </label>
           <input
            type="text"
            value={customerType}
            onChange={(e) => setCustomerType(e.target.value)}
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            required
          />
        </div> */}

        {/* <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">
            Customer City
          </label>
          <input
            type="text"
            value={customerCityName}
            onChange={(e) => setCustomerCityName(e.target.value)}
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            required
          />
        </div> */}

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">
            Max Sample Insights from DB 
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
      </div>

      <div className="border-t pt-6">
        <div className="flex items-center justify-between mb-4">
          <h3 className="text-lg font-semibold text-gray-900">Call Data</h3>
          <button
            type="button"
            onClick={handleAddCall}
            className="flex items-center gap-2 px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 transition-colors text-sm"
          >
            <Plus className="w-4 h-4" />
            Add Call
          </button>
        </div>

        <div className="space-y-4">
          {callData.map((call, index) => (
            <div key={index} className="p-4 border border-gray-200 rounded-lg space-y-3 bg-gray-50">
              <div className="flex items-center justify-between mb-2">
                <span className="text-sm font-medium text-gray-700">Call {index + 1}</span>
                {callData.length > 1 && (
                  <button
                    type="button"
                    onClick={() => handleRemoveCall(index)}
                    className="text-red-600 hover:text-red-700 transition-colors"
                  >
                    <Trash2 className="w-4 h-4" />
                  </button>
                )}
              </div>

              <div>
                <label className="block text-xs font-medium text-gray-600 mb-1">
                  Recording URL
                </label>
                <input
                  type="url"
                  value={call.call_recording_url}
                  onChange={(e) =>
                    handleCallDataChange(index, 'call_recording_url', e.target.value)
                  }
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent text-sm"
                  placeholder="https://example.com/recording.mp3"
                  required
                />
              </div>

              <div className="grid grid-cols-2 gap-3">
                <div>
                  <label className="block text-xs font-medium text-gray-600 mb-1">
                    Call Type
                  </label>
                  <input
                    type="text"
                    value={call.call_type}
                    onChange={(e) => handleCallDataChange(index, 'call_type', e.target.value)}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent text-sm"
                    required
                  />
                </div>

                <div>
                  <label className="block text-xs font-medium text-gray-600 mb-1">
                    Call Date
                  </label>
                  <input
                    type="datetime-local"
                    value={call.call_date.slice(0, 16).replace(' ', 'T')}
                    onChange={(e) =>
                      handleCallDataChange(
                        index,
                        'call_date',
                        e.target.value.replace('T', ' ') + ':00'
                      )
                    }
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent text-sm"
                    required
                  />
                </div>
              </div>
            </div>
          ))}
        </div>
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
        Generate Insights
      </button>
    </form>
  );
}
