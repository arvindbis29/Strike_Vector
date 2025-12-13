export interface CallData {
  call_recording_url: string;
  call_type: string;
  call_date: string;
}

export interface FinalInsightRequest {
  max_call_limit: number;
}

export interface GenerateInsightRequest {
  glid: number;
  executive_id: string;
  customer_type: string;
  customer_city_name: string;
  max_call_limit: number;
  call_data: CallData[];
}

export interface Insight {
  EnsightType: string;
  Concerns: string;
  Resolution: string;
  NextSteps: string;
  Alert: string;
  Sentiment: string;
  KeyPoints: string;
}

export interface InsightResponse {
  code: number;
  status: string;
  error: string;
  response: {
    ensights: Insight[];
  };
}
