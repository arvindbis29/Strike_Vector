import { useState } from "react";
import { callInsightsAPI } from "../utils/api";
import { Loader2 } from "lucide-react";

export default function InputForm({ setResponseData }) {
  const [form, setForm] = useState({
    glid: "",
    executiveId: "",
    customerType: "",
    customerCity: "",
    callData: [
      { text: "" }
    ]
  });

  const [loading, setLoading] = useState(false);

  const updateField = (field, value) => {
    setForm({ ...form, [field]: value });
  };

  const updateCallText = (index, value) => {
    const updated = [...form.callData];
    updated[index].text = value;
    setForm({ ...form, callData: updated });
  };

  const addCallRow = () => {
    setForm({ ...form, callData: [...form.callData, { text: "" }] });
  };

  const submitForm = async () => {
    setLoading(true);

    const payload = {
      glid: Number(form.glid),
      executive_id: form.executiveId,
      customer_type: form.customerType,
      customer_city_name: form.customerCity,
      call_data: form.callData.map((c) => c.text)
    };

    const response = await callInsightsAPI(JSON.stringify(payload));
    setResponseData(response);
    setLoading(false);
  };

  return (
    <div className="bg-white shadow-md rounded-xl p-6">
      <h2 className="text-2xl font-semibold mb-4">Generate Insights</h2>

      <div className="space-y-4">
        <div>
          <label className="block text-sm font-medium">GLID</label>
          <input
            type="number"
            className="input"
            value={form.glid}
            onChange={(e) => updateField("glid", e.target.value)}
          />
        </div>

        <div>
          <label className="block text-sm font-medium">Executive ID</label>
          <input
            className="input"
            value={form.executiveId}
            onChange={(e) => updateField("executiveId", e.target.value)}
          />
        </div>

        <div>
          <label className="block text-sm font-medium">Customer Type</label>
          <input
            className="input"
            value={form.customerType}
            onChange={(e) => updateField("customerType", e.target.value)}
          />
        </div>

        <div>
          <label className="block text-sm font-medium">Customer City</label>
          <input
            className="input"
            value={form.customerCity}
            onChange={(e) => updateField("customerCity", e.target.value)}
          />
        </div>

        <h3 className="text-lg font-semibold mt-4">Call Data</h3>
        {form.callData.map((call, idx) => (
          <textarea
            key={idx}
            className="input h-24"
            placeholder="Enter call transcript..."
            value={call.text}
            onChange={(e) => updateCallText(idx, e.target.value)}
          />
        ))}

        <button
          onClick={addCallRow}
          className="px-3 py-1 bg-gray-200 rounded-md hover:bg-gray-300"
        >
          + Add More Call Data
        </button>

        <button
          onClick={submitForm}
          className="w-full bg-blue-600 text-white py-2 rounded-lg hover:bg-blue-700 flex items-center justify-center"
          disabled={loading}
        >
          {loading ? <Loader2 className="animate-spin" /> : "Generate Insights"}
        </button>
      </div>
    </div>
  );
}
