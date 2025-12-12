export default function RawJsonView({ inputJson, responseData }) {
  return (
    <div>
      <h3 className="font-semibold mb-2">Raw JSON</h3>

      <div className="mb-3">
        <label className="block text-sm mb-1">Request</label>
        <pre className="p-3 rounded bg-gray-800 text-white overflow-auto max-h-40">
          {inputJson}
        </pre>
      </div>

      <div>
        <label className="block text-sm mb-1">Response</label>
        <pre className="p-3 rounded bg-gray-800 text-white overflow-auto max-h-80">
          {responseData ? JSON.stringify(responseData, null, 2) : "No response yet"}
        </pre>
      </div>
    </div>
  );
}
