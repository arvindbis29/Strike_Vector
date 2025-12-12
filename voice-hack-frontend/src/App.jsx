import { useState } from "react";
import InputForm from "./components/InputForm";
import InsightsView from "./components/InsightsView";

export default function App() {
  const [responseData, setResponseData] = useState(null);

  return (
    <div className="min-h-screen bg-gray-50 p-6">
      <div className="max-w-5xl mx-auto">
        {!responseData ? (
          <InputForm setResponseData={setResponseData} />
        ) : (
          <InsightsView
            responseData={responseData}
            goBack={() => setResponseData(null)}
          />
        )}
      </div>
    </div>
  );
}
