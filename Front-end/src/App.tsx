import { useState } from 'react';
import { Lightbulb, Loader2 } from 'lucide-react';
import FinalInsightForm from './components/FinalInsightForm';
import GenerateInsightForm from './components/GenerateInsightForm';
import ResultsDisplay from './components/ResultsDisplay';
import { InsightResponse } from './types/insights';

function App() {
  const [activeTab, setActiveTab] = useState<'final' | 'generate'>('final');
  const [results, setResults] = useState<InsightResponse | null>(null);
  const [isLoading, setIsLoading] = useState(false);

  const handleResults = (data: InsightResponse) => {
    setResults(data);
  };

  const handleLoading = (loading: boolean) => {
    setIsLoading(loading);
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 via-white to-blue-50">
      <div className="container mx-auto px-4 py-8 max-w-7xl">
        <div className="mb-8 text-center">
          <div className="flex items-center justify-center gap-3 mb-3">
            <div className="p-3 bg-blue-600 rounded-xl">
              <Lightbulb className="w-8 h-8 text-white" />
            </div>
            <h1 className="text-4xl font-bold text-gray-900">Call Insights Dashboard</h1>
          </div>
          <p className="text-gray-600 text-lg">
            Analyze call data and generate actionable insights
          </p>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
          <div className="space-y-6">
            <div className="bg-white rounded-xl shadow-lg p-6">
              <div className="flex gap-2 mb-6 p-1 bg-gray-100 rounded-lg">
                <button
                  onClick={() => setActiveTab('final')}
                  className={`flex-1 py-2.5 px-4 rounded-lg font-medium transition-all ${
                    activeTab === 'final'
                      ? 'bg-white text-blue-600 shadow-sm'
                      : 'text-gray-600 hover:text-gray-900'
                  }`}
                >
                  Final Insights
                </button>
                <button
                  onClick={() => setActiveTab('generate')}
                  className={`flex-1 py-2.5 px-4 rounded-lg font-medium transition-all ${
                    activeTab === 'generate'
                      ? 'bg-white text-blue-600 shadow-sm'
                      : 'text-gray-600 hover:text-gray-900'
                  }`}
                >
                  Generate Insights
                </button>
              </div>

              {isLoading && (
                <div className="mb-6 p-4 bg-blue-50 border border-blue-200 rounded-lg flex items-center gap-3">
                  <Loader2 className="w-5 h-5 text-blue-600 animate-spin" />
                  <span className="text-blue-700 font-medium">Processing your request...</span>
                </div>
              )}

              {activeTab === 'final' ? (
                <FinalInsightForm onResults={handleResults} onLoading={handleLoading} />
              ) : (
                <GenerateInsightForm onResults={handleResults} onLoading={handleLoading} />
              )}
            </div>
          </div>

          <div className="space-y-6">
            {results ? (
              <div className="bg-white rounded-xl shadow-lg p-6">
                <ResultsDisplay results={results} />
              </div>
            ) : (
              <div className="bg-white rounded-xl shadow-lg p-12 text-center">
                <div className="inline-flex p-4 bg-gray-100 rounded-full mb-4">
                  <Lightbulb className="w-12 h-12 text-gray-400" />
                </div>
                <h3 className="text-xl font-semibold text-gray-700 mb-2">
                  No Insights Yet
                </h3>
                <p className="text-gray-500">
                  Fill out the form and submit to see insights here
                </p>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}

export default App;
