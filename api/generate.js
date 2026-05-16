// Serverless AI generation endpoint template.
// Deploy this with Vercel, Netlify Functions, or another serverless host.
// Do not expose OPENAI_API_KEY in browser-side code.

export default async function handler(req, res) {
  if (req.method !== 'POST') {
    return res.status(405).json({ error: 'Method not allowed' });
  }

  const { prompt = '', tone = 'Clear and premium' } = req.body || {};

  if (!process.env.OPENAI_API_KEY) {
    return res.status(500).json({
      error: 'Missing OPENAI_API_KEY environment variable.'
    });
  }

  try {
    const response = await fetch('https://api.openai.com/v1/responses', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${process.env.OPENAI_API_KEY}`
      },
      body: JSON.stringify({
        model: process.env.OPENAI_MODEL || 'gpt-4.1-mini',
        input: [
          {
            role: 'system',
            content: 'You generate concise publishing asset copy and layout instructions. Return only valid JSON.'
          },
          {
            role: 'user',
            content: `Create a social/blog cover design for this brief: ${prompt}\nTone: ${tone}\nReturn JSON with eyebrow, title, subtitle, palette, template, and notes.`
          }
        ]
      })
    });

    if (!response.ok) {
      const errorText = await response.text();
      return res.status(response.status).json({ error: errorText });
    }

    const data = await response.json();
    const text = data.output_text || '{}';
    const parsed = JSON.parse(text);
    return res.status(200).json(parsed);
  } catch (error) {
    return res.status(500).json({ error: error.message });
  }
}
