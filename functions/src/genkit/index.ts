import {genkit, z} from "genkit";
import {googleAI, gemini15Pro} from "@genkit-ai/googleai";
import {onRequest} from "firebase-functions/v2/https";
import {logger} from "firebase-functions";

const ai = genkit({
  plugins: [googleAI()],
});

const AnomalySchema = z.object({
  id: z.string(),
  ncd_score: z.number().optional(),
  p_value: z.number().optional(),
  explanation: z.string().optional(),
  stream_type: z.string().optional(),
  detected_at: z.string().optional(),
});

const InputSchema = z.object({
  anomalies: z.array(AnomalySchema),
  query: z.string().optional(),
});

const OutputSchema = z.object({
  analysis: z.string(),
  risk_level: z.enum(["CRITICAL", "HIGH", "MEDIUM", "LOW"]),
  action_items: z.array(z.string()),
});

// Define the flow logic
const explainAnomalyFlow = ai.defineFlow({
  name: "explainAnomaly",
  inputSchema: InputSchema,
  outputSchema: OutputSchema,
}, async (input) => {
  const {anomalies, query} = input;

  const prompt = `
    Analyze these anomalies detected by Driftlock's compression-based system:
    ${anomalies.map((a: any) => `
      - ID: ${a.id}
      - Score: ${a.ncd_score}
      - P-value: ${a.p_value}
      - Explanation: ${a.explanation}
    `).join("\n")}

    User Query: ${query || "Provide insights"}

    Provide a structured analysis including:
    1. Executive summary
    2. Risk assessment
    3. Action items
  `;

  const result = await ai.generate({
    model: gemini15Pro,
    prompt: prompt,
    output: {format: "json", schema: OutputSchema},
  });

  if (!result.output) {
    throw new Error("Failed to generate analysis");
  }

  return result.output;
});

// Expose as Firebase Function
export const explainAnomaly = onRequest({cors: true}, async (request, response) => {
  try {
    // Basic auth check
    // In a real scenario, use request.headers.authorization and verify token
    // or use context from onCall if using callable functions
    // For onRequest, we need to manually verify
    // const auth = request.headers.authorization;
    // if (!auth) throw new Error("Unauthorized");

    const result = await explainAnomalyFlow(request.body);
    response.json(result);
  } catch (error: any) {
    logger.error("Genkit flow error", error);
    response.status(500).json({error: error.message || "Internal Server Error"});
  }
});
