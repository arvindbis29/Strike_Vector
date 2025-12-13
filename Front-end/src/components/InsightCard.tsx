import { AlertCircle, CheckCircle, TrendingUp, FileText, ArrowRight } from 'lucide-react';
import { Insight } from '../types/insights';

interface Props {
  insight: Insight;
}

export default function InsightCard({ insight }: Props) {
  const getSentimentColor = (sentiment: string) => {
    if (sentiment.toLowerCase().includes('positive')) return 'text-green-600 bg-green-50';
    if (sentiment.toLowerCase().includes('negative')) return 'text-red-600 bg-red-50';
    return 'text-gray-600 bg-gray-50';
  };

  const isFinalInsight = insight.EnsightType.toLowerCase() === 'final';

  return (
    <div className={`border rounded-xl p-6 space-y-5 ${isFinalInsight ? 'bg-blue-50 border-blue-200' : 'bg-white border-gray-200'}`}>
      <div className="flex items-start justify-between">
        <div className="flex items-center gap-3">
          <div className={`p-2 rounded-lg ${isFinalInsight ? 'bg-blue-600' : 'bg-blue-100'}`}>
            <FileText className={`w-5 h-5 ${isFinalInsight ? 'text-white' : 'text-blue-600'}`} />
          </div>
          <div>
            <h3 className="text-lg font-bold text-gray-900 capitalize">
              {insight.EnsightType.replace('_', ' ')}
            </h3>
            <span className={`inline-block mt-1 px-3 py-1 text-xs font-medium rounded-full ${getSentimentColor(insight.Sentiment)}`}>
              {insight.Sentiment}
            </span>
          </div>
        </div>
      </div>

      <div className="space-y-4">
        <div>
          <div className="flex items-center gap-2 mb-2">
            <AlertCircle className="w-4 h-4 text-orange-600" />
            <h4 className="text-sm font-semibold text-gray-700">Concerns</h4>
          </div>
          <p className="text-sm text-gray-600 leading-relaxed pl-6">{insight.Concerns}</p>
        </div>

        <div>
          <div className="flex items-center gap-2 mb-2">
            <CheckCircle className="w-4 h-4 text-green-600" />
            <h4 className="text-sm font-semibold text-gray-700">Resolution</h4>
          </div>
          <p className="text-sm text-gray-600 leading-relaxed pl-6">{insight.Resolution}</p>
        </div>

        <div>
          <div className="flex items-center gap-2 mb-2">
            <ArrowRight className="w-4 h-4 text-blue-600" />
            <h4 className="text-sm font-semibold text-gray-700">Next Steps</h4>
          </div>
          <p className="text-sm text-gray-600 leading-relaxed pl-6">{insight.NextSteps}</p>
        </div>

        {insight.Alert && insight.Alert.trim() && insight.Alert.toLowerCase() !== 'none identified in the provided calls.' && (
          <div className="bg-yellow-50 border border-yellow-200 rounded-lg p-4">
            <div className="flex items-start gap-2">
              <AlertCircle className="w-5 h-5 text-yellow-600 flex-shrink-0 mt-0.5" />
              <div>
                <h4 className="text-sm font-semibold text-yellow-800 mb-1">Alert</h4>
                <p className="text-sm text-yellow-700">{insight.Alert}</p>
              </div>
            </div>
          </div>
        )}

        <div>
          <div className="flex items-center gap-2 mb-2">
            <TrendingUp className="w-4 h-4 text-blue-600" />
            <h4 className="text-sm font-semibold text-gray-700">Key Points</h4>
          </div>
          <div className="pl-6 flex flex-wrap gap-2">
            {insight.KeyPoints.split(',').map((point, index) => (
              <span
                key={index}
                className="inline-block px-3 py-1 bg-blue-100 text-blue-700 text-xs font-medium rounded-full"
              >
                {point.trim()}
              </span>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
}
