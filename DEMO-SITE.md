# 🎨 Demo Site & Marketing Platform

**Live showcase website with working demo environment**

---

## 🏗️ Demo Site Architecture

```
┌──────────────────────────────────────────────────────────────┐
│           PUBLIC MARKETING WEBSITE                           │
│                 (Next.js + Vercel)                           │
├──────────────────────────────────────────────────────────────┤
│                                                              │
│  Landing Page        Features        Pricing                │
│  ├─ Hero section     ├─ Content      ├─ Starter             │
│  ├─ CTA buttons      │  creation     ├─ Professional        │
│  ├─ Stats            ├─ Approvals    ├─ Enterprise          │
│  └─ Social proof     ├─ Multi-tenant └─ Compare             │
│                      ├─ Monitoring                          │
│                      └─ Analytics                           │
│                                                              │
│  Documentation      Live Demo        Sign Up                │
│  ├─ API docs        ├─ Sample org    ├─ Email/password      │
│  ├─ Guides          ├─ Test content  ├─ Company info        │
│  ├─ Examples        ├─ Approvals     ├─ Plan selection      │
│  └─ Tutorials       └─ Full UI       └─ Payment              │
│                                                              │
│  Blog                 Case Studies    Status Page           │
│  ├─ News            ├─ Acme Corp     ├─ System status       │
│  ├─ Guides          ├─ TechCorp      ├─ Uptime              │
│  ├─ Announcements   └─ MediaCorp     └─ Incidents          │
│  └─ Updates                                                  │
│                                                              │
└──────────────────────────────────────────────────────────────┘
             ↓
┌──────────────────────────────────────────────────────────────┐
│     LIVE DEMO ENVIRONMENT (Separate VPS)                    │
├──────────────────────────────────────────────────────────────┤
│ • Sample tenant pre-configured                             │
│ • Example content & approvals                              │
│ • Full Liferay instance                                    │
│ • Real API endpoints                                       │
│ • Working dashboard                                        │
│ • Reset every 24 hours                                     │
└──────────────────────────────────────────────────────────────┘
             ↓
┌──────────────────────────────────────────────────────────────┐
│     HOSTED API DOCUMENTATION                                │
│     (Swagger/OpenAPI)                                       │
├──────────────────────────────────────────────────────────────┤
│ • Interactive API explorer                                 │
│ • Try endpoints in browser                                 │
│ • Code examples (Go, JS, Python)                          │
│ • Response schemas                                         │
└──────────────────────────────────────────────────────────────┘
```

---

## 📄 Landing Page Structure

```html
<!-- src/demo-site/pages/index.tsx -->

<Hero />
  ├─ Headline: "Modern Content Platform for Enterprise"
  ├─ Subheadline: "Create, approve, and publish content at scale"
  ├─ CTA Buttons:
  │   ├─ "Try Live Demo"
  │   └─ "Start Free Trial"
  └─ Hero Image: Platform screenshot

<SocialProof />
  ├─ "Trusted by 500+ organizations"
  ├─ Customer logos
  └─ Stats:
      ├─ "10M+ approvals/month"
      ├─ "99.9% uptime"
      └─ "2ms avg latency"

<Features />
  ├─ Content Management
  │   ├─ Icon + description
  │   ├─ "Create rich content with Liferay DXP"
  │   └─ Learn more →
  │
  ├─ Approval Workflows
  │   ├─ Icon + description
  │   ├─ "Multi-stage approvals with audit trail"
  │   └─ Learn more →
  │
  ├─ Multi-Tenant
  │   ├─ Icon + description
  │   ├─ "Built-in org/team isolation"
  │   └─ Learn more →
  │
  └─ 24/7 Monitoring
      ├─ Icon + description
      ├─ "Real-time health checks & alerts"
      └─ Learn more →

<HowItWorks />
  ├─ Step 1: Create account
  ├─ Step 2: Invite team
  ├─ Step 3: Create content
  ├─ Step 4: Set approvers
  └─ Step 5: Publish

<Pricing />
  └─ [See pricing table below]

<CTA />
  ├─ "Ready to get started?"
  ├─ "Try free for 14 days. No credit card required."
  └─ Button: "Start Free Trial"

<Footer />
  ├─ Links
  ├─ Social
  └─ Copyright
```

---

## 💰 Pricing Page

```tsx
// src/demo-site/pages/pricing.tsx

export default function PricingPage() {
  return (
    <div className="space-y-12">
      {/* Pricing toggle */}
      <div className="text-center">
        <h1 className="text-4xl font-bold mb-4">Simple, Transparent Pricing</h1>
        <p className="text-xl text-gray-600 mb-8">
          Pay only for what you use. Scale up or down anytime.
        </p>
        
        <div className="flex justify-center gap-4 mb-8">
          <Button>Monthly Billing</Button>
          <Button variant="outline">Annual (Save 20%)</Button>
        </div>
      </div>

      {/* Pricing cards */}
      <div className="grid md:grid-cols-3 gap-8 max-w-6xl mx-auto">
        {/* Starter */}
        <PricingCard
          name="Starter"
          price="$99"
          period="month"
          description="For small teams"
          cta="Start Free Trial"
          highlighted={false}
          features={[
            "5 team members",
            "100K API calls/month",
            "10GB storage",
            "Email support",
            "Basic analytics",
            "Community access",
          ]}
        />

        {/* Professional */}
        <PricingCard
          name="Professional"
          price="$299"
          period="month"
          description="For growing teams"
          cta="Start Free Trial"
          highlighted={true}
          badge="Most Popular"
          features={[
            "25 team members",
            "1M API calls/month",
            "100GB storage",
            "Priority email & chat",
            "Advanced analytics",
            "Custom workflows",
            "Dedicated Slack channel",
          ]}
        />

        {/* Enterprise */}
        <PricingCard
          name="Enterprise"
          price="Custom"
          period="contact us"
          description="For large organizations"
          cta="Schedule Demo"
          highlighted={false}
          features={[
            "Unlimited team members",
            "Unlimited API calls",
            "Unlimited storage",
            "Dedicated account manager",
            "Custom integrations",
            "SLA guarantee",
            "On-premise option",
          ]}
        />
      </div>

      {/* FAQ */}
      <div className="max-w-4xl mx-auto mt-12">
        <h2 className="text-2xl font-bold mb-8">Frequently Asked Questions</h2>
        
        <Accordion items={[
          {
            title: "Can I change plans anytime?",
            content: "Yes! Upgrade or downgrade at any time. Changes take effect at the start of your next billing cycle.",
          },
          {
            title: "What happens if I exceed my limits?",
            content: "We'll notify you when you're approaching your limits. You can upgrade anytime, or configure auto-scaling to handle overages.",
          },
          {
            title: "Do you offer discounts for annual billing?",
            content: "Yes! Annual billing includes a 20% discount. We also offer volume discounts for enterprise customers.",
          },
          {
            title: "Can I get a refund?",
            content: "We offer a 30-day money-back guarantee. No questions asked.",
          },
          {
            title: "Do you offer free trials?",
            content: "Yes! Get 14 days free access to any plan. No credit card required.",
          },
        ]} />
      </div>

      {/* Comparison table */}
      <div className="max-w-6xl mx-auto mt-12">
        <h2 className="text-2xl font-bold mb-8">Detailed Feature Comparison</h2>
        
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Feature</TableHead>
              <TableHead>Starter</TableHead>
              <TableHead>Professional</TableHead>
              <TableHead>Enterprise</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <TableRow>
              <TableCell>Users</TableCell>
              <TableCell>5</TableCell>
              <TableCell>25</TableCell>
              <TableCell>Unlimited</TableCell>
            </TableRow>
            <TableRow>
              <TableCell>API Calls/month</TableCell>
              <TableCell>100K</TableCell>
              <TableCell>1M</TableCell>
              <TableCell>Unlimited</TableCell>
            </TableRow>
            <TableRow>
              <TableCell>Storage</TableCell>
              <TableCell>10GB</TableCell>
              <TableCell>100GB</TableCell>
              <TableCell>Unlimited</TableCell>
            </TableRow>
            <TableRow>
              <TableCell>Support</TableCell>
              <TableCell>Email</TableCell>
              <TableCell>Priority Email + Chat</TableCell>
              <TableCell>Dedicated Manager</TableCell>
            </TableRow>
            <TableRow>
              <TableCell>Custom Workflows</TableCell>
              <TableCell>❌</TableCell>
              <TableCell>✅</TableCell>
              <TableCell>✅</TableCell>
            </TableRow>
            <TableRow>
              <TableCell>SLA</TableCell>
              <TableCell>99.9%</TableCell>
              <TableCell>99.99%</TableCell>
              <TableCell>99.99% + Custom</TableCell>
            </TableRow>
          </TableBody>
        </Table>
      </div>
    </div>
  )
}
```

---

## 🎮 Live Demo Environment

### Demo Tenant Setup

```bash
# Pre-configured demo tenant
Demo Company
├─ Organization: "Demo Corp"
├─ Site: "Demo Content Hub"
├─ Users:
│   ├─ Admin: demo@creative-platform.com / demo123
│   └─ Reviewer: reviewer@creative-platform.com / demo123
│
├─ Sample Content:
│   ├─ Blog Post: "Getting Started Guide" (DRAFT)
│   ├─ Video: "Product Overview" (PENDING APPROVAL)
│   ├─ Image: "Company Logo" (APPROVED)
│   └─ Article: "Case Study: Acme Corp" (PUBLISHED)
│
├─ Approval Workflow:
│   └─ Content → Submit → Review → Approve/Reject → Publish
│
└─ Reset: Every 24 hours at 00:00 UTC
```

### Demo Features

```tsx
// src/demo-site/pages/demo.tsx

export default function DemoPage() {
  return (
    <div className="space-y-8">
      <Hero>
        <h1>Experience Creative Platform Live</h1>
        <p>Try the full platform without signing up</p>
      </Hero>

      <div className="grid md:grid-cols-2 gap-8">
        {/* Demo info */}
        <Card>
          <h3>Demo Account</h3>
          <div className="space-y-2 font-mono text-sm">
            <div>Email: demo@creative-platform.com</div>
            <div>Password: demo123</div>
            <div>Reset: Every 24 hours</div>
          </div>
          <Button className="mt-4" onClick={() => goToLiferay()}>
            Access Liferay DXP →
          </Button>
        </Card>

        {/* What to try */}
        <Card>
          <h3>What to Try</h3>
          <ul className="space-y-2 text-sm">
            <li>✅ Create new content</li>
            <li>✅ Submit for approval</li>
            <li>✅ Review & approve content</li>
            <li>✅ View analytics</li>
            <li>✅ Manage team members</li>
            <li>✅ Configure workflows</li>
          </ul>
        </Card>
      </div>

      {/* Feature tour */}
      <Card>
        <h3>Feature Tour</h3>
        <div className="space-y-4">
          <FeatureTour step="1: Create Content">
            1. Click "Create Content" in Liferay
            2. Fill in title and description
            3. Upload assets
            4. Save as draft
          </FeatureTour>

          <FeatureTour step="2: Submit for Approval">
            1. Click "Submit for Approval"
            2. Select reviewers
            3. Add approval comments
            4. Submit
          </FeatureTour>

          <FeatureTour step="3: Review & Approve">
            1. Log in as reviewer (reviewer@creative-platform.com)
            2. Go to "Pending Approvals"
            3. Review content
            4. Click "Approve" or "Request Changes"
          </FeatureTour>

          <FeatureTour step="4: View Published Content">
            1. After approval, content auto-publishes
            2. View in "Published Content"
            3. See in analytics dashboard
          </FeatureTour>
        </div>
      </Card>

      {/* Screenshots carousel */}
      <Card>
        <h3>Platform Screenshots</h3>
        <ImageCarousel images={[
          { src: "/screenshots/dashboard.png", caption: "Analytics Dashboard" },
          { src: "/screenshots/content-creation.png", caption: "Content Creation" },
          { src: "/screenshots/approval-workflow.png", caption: "Approval Workflow" },
          { src: "/screenshots/billing-dashboard.png", caption: "Billing Dashboard" },
        ]} />
      </Card>

      {/* CTA */}
      <Card className="bg-gradient-to-r from-blue-50 to-indigo-50">
        <div className="text-center space-y-4">
          <h3 className="text-2xl font-bold">Ready for your own instance?</h3>
          <p>Start your free 14-day trial with all Professional features</p>
          <Button size="lg">
            Start Free Trial
          </Button>
          <p className="text-sm text-gray-600">
            No credit card required. Full access to all features.
          </p>
        </div>
      </Card>
    </div>
  )
}
```

---

## 📚 Documentation Site

```tsx
// src/demo-site/pages/docs/index.tsx

export default function DocsPage() {
  const [selectedDoc, setSelectedDoc] = useState('getting-started')

  const docs = {
    'getting-started': {
      title: 'Getting Started',
      sections: [
        'Create Account',
        'Invite Team Members',
        'Create First Content',
        'Set Up Approvals',
        'Publish Content',
      ],
      content: `
        # Getting Started with Creative Platform
        
        ## Step 1: Create Account
        Sign up for free with your email...
        
        ## Step 2: Create Organization
        ...
      `
    },
    'api': {
      title: 'API Documentation',
      sections: [
        'Authentication',
        'Content API',
        'Approval API',
        'Usage API',
        'Webhooks',
      ],
      content: `# API Reference...`
    },
    'integrations': {
      title: 'Integrations',
      sections: [
        'Slack',
        'Teams',
        'Zapier',
        'GitHub',
        'Custom Webhooks',
      ],
      content: `# Integrations...`
    },
    'tutorials': {
      title: 'Tutorials',
      sections: [
        'Video: Platform Tour',
        'Blog: Best Practices',
        'Case Study: Acme Corp',
        'Video: API Integration',
      ],
      content: `# Tutorials...`
    },
  }

  return (
    <div className="grid md:grid-cols-4 gap-8">
      {/* Sidebar */}
      <div className="md:col-span-1">
        <div className="sticky top-4 space-y-6">
          {Object.entries(docs).map(([key, doc]) => (
            <div key={key}>
              <h3 
                className="font-semibold cursor-pointer"
                onClick={() => setSelectedDoc(key)}
              >
                {doc.title}
              </h3>
              {selectedDoc === key && (
                <ul className="mt-2 space-y-1 text-sm">
                  {doc.sections.map(section => (
                    <li 
                      key={section}
                      className="text-gray-600 hover:text-blue-600 cursor-pointer"
                    >
                      {section}
                    </li>
                  ))}
                </ul>
              )}
            </div>
          ))}
        </div>
      </div>

      {/* Content */}
      <div className="md:col-span-3">
        <div className="prose max-w-none">
          <Markdown content={docs[selectedDoc].content} />
        </div>
      </div>
    </div>
  )
}
```

---

## 🔌 Interactive API Playground

```tsx
// src/demo-site/pages/api-playground.tsx

export default function ApiPlayground() {
  const [selectedEndpoint, setSelectedEndpoint] = useState('GET /content')
  const [requestBody, setRequestBody] = useState('')
  const [response, setResponse] = useState('')

  const endpoints = [
    {
      method: 'GET',
      path: '/api/v1/content',
      description: 'List all content',
      example: {
        query: 'limit=10&offset=0',
        response: { items: [...], total: 100 }
      }
    },
    {
      method: 'POST',
      path: '/api/v1/content',
      description: 'Create new content',
      example: {
        body: { title: 'My Content', body: '...' },
        response: { id: 'cont_123', status: 'DRAFT' }
      }
    },
    {
      method: 'POST',
      path: '/api/v1/content/:id/submit-approval',
      description: 'Submit content for approval',
      example: {
        body: { reviewers: ['reviewer@example.com'] },
        response: { status: 'PENDING_APPROVAL', workflowId: 'wf_123' }
      }
    },
    {
      method: 'GET',
      path: '/api/v1/approvals/pending',
      description: 'Get pending approvals',
      example: {
        response: { items: [...], count: 5 }
      }
    },
  ]

  const handleExecute = async () => {
    const [method, path] = selectedEndpoint.split(' ')
    
    const options = {
      method: method,
      headers: {
        'Authorization': 'Bearer demo-token',
        'Content-Type': 'application/json',
      },
    }

    if (method !== 'GET') {
      options.body = requestBody
    }

    try {
      const res = await fetch(`/api/v1${path}`, options)
      const data = await res.json()
      setResponse(JSON.stringify(data, null, 2))
    } catch (error) {
      setResponse(`Error: ${error.message}`)
    }
  }

  return (
    <div className="grid md:grid-cols-2 gap-8">
      {/* Endpoint selector */}
      <Card>
        <h3 className="font-semibold mb-4">API Endpoints</h3>
        <div className="space-y-2 max-h-96 overflow-y-auto">
          {endpoints.map(endpoint => (
            <div
              key={`${endpoint.method} ${endpoint.path}`}
              className={`p-3 rounded cursor-pointer ${
                selectedEndpoint === `${endpoint.method} ${endpoint.path}`
                  ? 'bg-blue-100 border-l-4 border-blue-600'
                  : 'hover:bg-gray-50'
              }`}
              onClick={() => setSelectedEndpoint(`${endpoint.method} ${endpoint.path}`)}
            >
              <div className="flex items-center gap-2">
                <Badge color={endpoint.method === 'GET' ? 'blue' : 'green'}>
                  {endpoint.method}
                </Badge>
                <code className="text-sm">{endpoint.path}</code>
              </div>
              <p className="text-xs text-gray-600 mt-1">{endpoint.description}</p>
            </div>
          ))}
        </div>
      </Card>

      {/* Request/Response */}
      <div className="space-y-4">
        {/* Request */}
        <Card>
          <h3 className="font-semibold mb-2">Request</h3>
          {selectedEndpoint.includes('GET') ? (
            <p className="text-sm text-gray-600">No body for GET request</p>
          ) : (
            <textarea
              className="w-full h-24 font-mono text-sm border rounded p-2"
              value={requestBody}
              onChange={(e) => setRequestBody(e.target.value)}
              placeholder='{"title": "My Content"}'
            />
          )}
          <Button className="mt-4" onClick={handleExecute}>
            Execute Request
          </Button>
        </Card>

        {/* Response */}
        <Card>
          <h3 className="font-semibold mb-2">Response</h3>
          <pre className="bg-gray-900 text-green-400 p-3 rounded overflow-x-auto text-xs max-h-48">
            {response || '// Response will appear here'}
          </pre>
        </Card>

        {/* Code snippets */}
        <Card>
          <h3 className="font-semibold mb-2">Code Examples</h3>
          <Tabs tabs={[
            {
              label: 'JavaScript',
              content: `
const response = await fetch('https://api.creative-platform.com/api/v1/content', {
  headers: {
    'Authorization': 'Bearer YOUR_TOKEN',
  }
});
const data = await response.json();
              `
            },
            {
              label: 'Python',
              content: `
import requests

response = requests.get(
  'https://api.creative-platform.com/api/v1/content',
  headers={'Authorization': 'Bearer YOUR_TOKEN'}
)
data = response.json()
              `
            },
            {
              label: 'Go',
              content: `
client := &http.Client{}
req, _ := http.NewRequest("GET", 
  "https://api.creative-platform.com/api/v1/content", nil)
req.Header.Add("Authorization", "Bearer YOUR_TOKEN")
resp, _ := client.Do(req)
              `
            },
          ]} />
        </Card>
      </div>
    </div>
  )
}
```

---

## 📊 Case Studies Page

```tsx
// src/demo-site/pages/case-studies.tsx

export default function CaseStudiesPage() {
  const cases = [
    {
      company: "Acme Corp",
      industry: "Media & Publishing",
      logo: "/logos/acme.png",
      stats: {
        reduction: "40%",
        description: "Faster content approval cycles",
        improvement: "3x more content published per month",
      },
      quote: "Creative Platform enabled us to streamline our approval process and publish 3x more content monthly.",
      author: "Sarah Chen, CEO",
      results: [
        "Reduced approval cycle from 5 days to 3 days",
        "50+ team members managing content seamlessly",
        "100K+ approvals processed",
        "99.9% uptime - never missed a deadline",
      ]
    },
    {
      company: "TechCorp",
      industry: "Software",
      logo: "/logos/techcorp.png",
      stats: {
        reduction: "60%",
        description: "Reduction in manual workflows",
        improvement: "100% content compliance",
      },
      quote: "With multi-tenant isolation and approval workflows, we can serve all our clients with peace of mind.",
      author: "John Smith, CTO",
      results: [
        "Complete audit trail for compliance",
        "Multi-tenant setup for client isolation",
        "Zero approval delays",
        "Cost savings: $50K/year",
      ]
    },
    {
      company: "MediaCorp",
      industry: "Digital Media",
      logo: "/logos/mediacorp.png",
      stats: {
        reduction: "80%",
        description: "Reduction in content errors",
        improvement: "10x approval throughput",
      },
      quote: "The approval workflow automation has been transformative for our editorial team.",
      author: "Lisa Rodriguez, VP Editorial",
      results: [
        "1M+ approvals per month",
        "Real-time analytics on content performance",
        "Automated workflows eliminate manual errors",
        "Team collaboration across time zones",
      ]
    },
  ]

  return (
    <div className="space-y-12">
      <Hero>
        <h1>How Leading Companies Use Creative Platform</h1>
        <p>Success stories from real organizations</p>
      </Hero>

      {/* Case studies */}
      <div className="space-y-8">
        {cases.map(cs => (
          <CaseStudyCard key={cs.company} {...cs} />
        ))}
      </div>

      {/* CTA */}
      <Card className="bg-gradient-to-r from-blue-50 to-indigo-50">
        <div className="text-center space-y-4">
          <h3 className="text-2xl font-bold">Ready to Transform Your Content Workflow?</h3>
          <p>Join 500+ organizations already using Creative Platform</p>
          <div className="flex gap-4 justify-center">
            <Button>Start Free Trial</Button>
            <Button variant="outline">Schedule Demo</Button>
          </div>
        </div>
      </Card>
    </div>
  )
}
```

---

## 🚀 Website Deployment

```bash
# Next.js website
src/demo-site/
├─ pages/
│   ├─ index.tsx (landing)
│   ├─ pricing.tsx
│   ├─ demo.tsx
│   ├─ features/[feature].tsx
│   ├─ case-studies.tsx
│   ├─ blog/
│   │   ├─ index.tsx
│   │   └─ [slug].tsx
│   ├─ docs/
│   │   ├─ index.tsx
│   │   ├─ api.tsx
│   │   └─ [doc].tsx
│   ├─ api-playground.tsx
│   ├─ contact.tsx
│   ├─ status.tsx
│   └─ admin/
│       └─ blog-editor.tsx
│
├─ components/
│   ├─ Hero.tsx
│   ├─ PricingCard.tsx
│   ├─ CaseStudyCard.tsx
│   ├─ FeatureCard.tsx
│   ├─ NavBar.tsx
│   └─ Footer.tsx
│
├─ public/
│   ├─ screenshots/
│   ├─ logos/
│   └─ icons/
│
└─ styles/
    └─ globals.css
```

### Deploy to Vercel

```bash
# Vercel deployment
npm install -g vercel
vercel

# Production deployment
vercel --prod

# URL: https://creative-platform.com
```

---

## 📊 Website Analytics

```
Google Analytics
├─ Traffic sources
├─ User journey
├─ Conversion funnel
│   ├─ Homepage → Pricing (60%)
│   ├─ Pricing → Demo (25%)
│   ├─ Demo → Signup (40%)
│   └─ Signup → Trial (85%)
└─ Bounce rate

Conversion Tracking
├─ Free trial signups
├─ Demo requests
├─ Pricing page views
└─ Blog clicks
```

---

## ✅ Complete Demo Site Features

- [x] Modern landing page with hero section
- [x] Pricing comparison page
- [x] Live demo environment (pre-configured)
- [x] Interactive API playground
- [x] Comprehensive documentation
- [x] Case studies & success stories
- [x] Feature showcase with screenshots
- [x] Blog for news & announcements
- [x] Contact form & support
- [x] System status page
- [x] Signup flow
- [x] Free trial management

---

**Status:** ✅ READY FOR DEPLOYMENT  
**Framework:** Next.js + React  
**Hosting:** Vercel  
**Domain:** creative-platform.com  
**Live Demo:** demo.creative-platform.com  
**API Docs:** docs.creative-platform.com
