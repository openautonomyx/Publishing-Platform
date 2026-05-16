// Multi-platform, multilingual distribution endpoint.

export default async function handler(req, res) {
  if (req.method !== 'POST') {
    return res.status(405).json({ error: 'Method not allowed' });
  }

  const {
    title = 'Untitled Content',
    body = '',
    platforms = ['github-pages'],
    languages = ['en'],
    formats = ['article']
  } = req.body || {};

  const jobs = [];

  for (const platform of platforms) {
    for (const language of languages) {
      for (const format of formats) {
        jobs.push({
          platform,
          language,
          format,
          status: 'ready',
          title,
          preview: body.slice(0, 160)
        });
      }
    }
  }

  res.status(200).json({
    title,
    totalJobs: jobs.length,
    jobs
  });
}
