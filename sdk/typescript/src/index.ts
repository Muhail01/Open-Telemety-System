export type RecommendationTrace = {
  decisionId: string;
  itemId: string;
  rank: number;
  surface: string;
};

export type TelemetryEventName = "recommendation_impression" | "recommendation_click";

export type TelemetryClientOptions = {
  endpoint: string;
  sessionId: string;
  fetchImpl?: typeof fetch;
};

export class TelemetryClient {
  private readonly endpoint: string;
  private readonly sessionId: string;
  private readonly fetchImpl: typeof fetch;
  private readonly sentImpressions = new Set<string>();

  constructor(options: TelemetryClientOptions) {
    this.endpoint = options.endpoint.replace(/\/$/, "");
    this.sessionId = options.sessionId;
    this.fetchImpl = options.fetchImpl ?? fetch;
  }

  async impression(trace: RecommendationTrace): Promise<boolean> {
    const key = `${trace.decisionId}:${trace.surface}:${trace.itemId}:${trace.rank}`;
    if (this.sentImpressions.has(key)) return false;
    await this.send("recommendation_impression", trace);
    this.sentImpressions.add(key);
    return true;
  }

  async click(trace: RecommendationTrace): Promise<void> {
    await this.send("recommendation_click", trace);
  }

  private async send(eventType: TelemetryEventName, trace: RecommendationTrace): Promise<void> {
    const response = await this.fetchImpl(`${this.endpoint}/v1/events`, {
      method: "POST",
      headers: { "content-type": "application/json" },
      body: JSON.stringify({
        event_id: createEventId(),
        event_type: eventType,
        session_id: this.sessionId,
        decision_id: trace.decisionId,
        item_id: trace.itemId,
        surface: trace.surface,
        properties: { rank: trace.rank },
      }),
    });
    if (!response.ok) {
      throw new Error(`telemetry request failed with ${response.status}`);
    }
  }
}

function createEventId(): string {
  if (globalThis.crypto?.randomUUID) return globalThis.crypto.randomUUID();
  return `evt_${Date.now().toString(36)}_${Math.random().toString(36).slice(2)}`;
}
