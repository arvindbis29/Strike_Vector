export default function JsonInputBox({ inputJson, setInputJson, onSubmit, loading }) {
  const onPretty = () => {
    try {
      const p = JSON.stringify(JSON.parse(inputJson), null, 2);
      setInputJson(p);
    } catch (e) {
      alert("Invalid JSON");
    }
  };

  const onClear = () => {
    setInputJson("");
  };

  return (
    <div>
      <h3 className="font-semibold mb-2">Input JSON</h3>
      <textarea
        className="w-full h-64 p-3 border rounded-lg bg-transparent"
        value={inputJson}
        onChange={(e) => setInputJson(e.target.value)}
      />

      <div className="flex gap-2 mt-3">
        <button onClick={onPretty} className="px-3 py-1 border rounded">Pretty</button>
        <button onClick={onClear} className="px-3 py-1 border rounded">Clear</button>
        <button onClick={onSubmit} className="ml-auto px-4 py-2 bg-blue-600 text-white rounded">
          {loading ? "Generating..." : "Generate Insights"}
        </button>
      </div>
    </div>
  );
}
