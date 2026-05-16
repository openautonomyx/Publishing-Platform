// Secure multimodal AI endpoint template for serverless deployment.
// Works with Vercel-style serverless functions.
// Keep OPENAI_API_KEY server-side only. Never expose it in GitHub Pages/browser code.

export default async function handler(req, res) {
  if (req.method !== 'POST') {
    return res.status(405).json({ error: 'Method not allowed' });
  }

  const { prompt = '', tone = 'Clear and premium', imageBase64 = null, canonicalUrl = '' } = req.body || {};

  if (!process.env.OPENAI_API_KEY) {
    return res.status(500).json({ error: 'Missing OPENAI_API_KEY environment variable.' });
  }

  const userContent = [
    {
      type: 'input_text',
      text: `Create a multimodal publishing asset plan.\nPrompt: ${prompt}\nTone: ${tone}\nCanonical URL: ${canonicalUrl || 'Not provided'}\n\nReturn only valid JSON with these keys: eyebrow, title, subtitle, category, tags, seoTitle, metaDescription, canonicalUrl, ogTitle, ogDescription, ogImageAlt, suggestedSlug, layout, palette, imageAnalysis, imageEditingInstructions, designNotes.`
    }
  ];

  if (imageBase64) {
    userContent.push({ type: 'input_image', image_url: imageBase64 });
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
            content: 'You are an expert multimodal publishing strategist, SEO editor, and visual designer. Return compact, valid JSON only.'
          },
          {
            role: 'user',
            content: userContent
          }
        ]
      })
    });

    if (!response.ok) {
      return res.status(response.status).json({ error: await response.text() });
    }

    const data = await response.json();
    const text = data.output_text || '{}';
    return res.status(200).json(JSON.parse(text));
  } catch (error) {
    return res.status(500).json({ error: error.message });
  }
}
