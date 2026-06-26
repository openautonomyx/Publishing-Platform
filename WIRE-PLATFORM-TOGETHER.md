# Wire Platform Together - Complete Integration Guide

**Connect Liferay DXP UI ↔ OpenAutonomyX Backend Services**

---

## Architecture: End-to-End Wiring

```
┌─────────────────────────────────────────┐
│      Liferay DXP Portal (8080)          │
│  ┌──────────────────────────────────┐   │
│  │  Content Creator Portlet         │   │
│  │  Format Converter Portlet        │   │
│  │  Integrations Manager Portlet    │   │
│  │  Analytics Dashboard Portlet     │   │
│  └──────────────────────────────────┘   │
│              ↓ REST API Calls ↓          │
└─────────────────────────────────────────┘
         ↓ ↓ ↓ ↓ ↓ ↓ ↓ ↓ ↓ ↓ ↓
┌─────────────────────────────────────────┐
│   OpenAutonomyX Backend Services        │
├─────────────────────────────────────────┤
│  API Gateway (3000)                     │
│  ├─ Content (3002)                      │
│  ├─ Blog (3009)                         │
│  ├─ Formats (3011)                      │
│  ├─ Integrations (3010)                 │
│  ├─ Analytics (3005)                    │
│  └─ +15 more services                   │
│              ↓                           │
│  PostgreSQL | Redis | Ollama            │
└─────────────────────────────────────────┘
```

---

## Step 1: Create API Client Library

### File: `/tmp/api-client.ts`

```typescript
// Shared API client for all Liferay portlets

export class OpenAutonomyXClient {
  private baseUrl: string;
  private apiKey: string;

  constructor(baseUrl: string = 'http://localhost:3000', apiKey: string) {
    this.baseUrl = baseUrl;
    this.apiKey = apiKey;
  }

  private async request(endpoint: string, options: any = {}) {
    const url = `${this.baseUrl}${endpoint}`;
    const response = await fetch(url, {
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${this.apiKey}`,
        ...options.headers,
      },
      ...options,
    });

    if (!response.ok) {
      throw new Error(`API Error: ${response.statusText}`);
    }

    return response.json();
  }

  // Content API
  async createContent(data: any) {
    return this.request('/api/v1/content/create', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  async getContent(id: string) {
    return this.request(`/api/v1/content/${id}`);
  }

  async listContent(params: any = {}) {
    const query = new URLSearchParams(params).toString();
    return this.request(`/api/v1/content?${query}`);
  }

  async updateContent(id: string, data: any) {
    return this.request(`/api/v1/content/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    });
  }

  async publishContent(id: string) {
    return this.request(`/api/v1/content/${id}/publish`, {
      method: 'POST',
    });
  }

  // Blog API
  async getBlogPosts(params: any = {}) {
    const query = new URLSearchParams(params).toString();
    return this.request(`/api/v1/blog/posts?${query}`);
  }

  async createBlogPost(data: any) {
    return this.request('/api/v1/blog/posts', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  async publishToBlog(id: string) {
    return this.request(`/api/v1/blog/posts/${id}/publish`, {
      method: 'POST',
    });
  }

  // Formats API
  async convertToFormat(format: string, data: any) {
    return this.request(`/api/v1/formats/${format}`, {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  async getFormats() {
    return this.request('/api/v1/formats/supported');
  }

  // Integrations API
  async getIntegrations() {
    return this.request('/api/v1/integrations');
  }

  async publishToIntegration(id: string, data: any) {
    return this.request(`/api/v1/integrations/${id}/publish`, {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  async createIntegration(data: any) {
    return this.request('/api/v1/integrations', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  // Analytics API
  async getAnalytics(params: any = {}) {
    const query = new URLSearchParams(params).toString();
    return this.request(`/api/v1/analytics?${query}`);
  }

  async trackEvent(event: any) {
    return this.request('/api/v1/analytics/events', {
      method: 'POST',
      body: JSON.stringify(event),
    });
  }
}

export default OpenAutonomyXClient;
```

---

## Step 2: Create Liferay Portlet Integration

### Java Class: Content Creator Portlet

```java
package com.openautonomyx.portlets;

import javax.portlet.*;
import org.osgi.service.component.annotations.Component;
import com.liferay.portal.kernel.portlet.bridges.mvc.MVCPortlet;
import java.io.IOException;

@Component(
  immediate = true,
  property = {
    "com.liferay.portlet.display-category=category.publishing",
    "javax.portlet.display-name=Content Creator",
    "javax.portlet.init-param.template-path=/",
    "javax.portlet.init-param.view-template=/create-content.jsp",
    "javax.portlet.resource-bundle=content.Language"
  },
  service = Portlet.class
)
public class ContentCreatorPortlet extends MVCPortlet {

  private static final String API_GATEWAY_URL = 
    System.getenv("OPENAUTONOMYX_API_GATEWAY_URL") != null ? 
    System.getenv("OPENAUTONOMYX_API_GATEWAY_URL") : 
    "http://localhost:3000";

  private static final String API_KEY = 
    System.getenv("OPENAUTONOMYX_API_KEY") != null ? 
    System.getenv("OPENAUTONOMYX_API_KEY") : 
    "dev-key";

  @Override
  public void doView(RenderRequest renderRequest, RenderResponse renderResponse)
    throws IOException, PortletException {
    
    try {
      // Get list of integrations available
      String integrationsUrl = API_GATEWAY_URL + "/api/v1/integrations";
      HttpClient httpClient = HttpClientBuilder.create().build();
      HttpGet request = new HttpGet(integrationsUrl);
      request.setHeader("Authorization", "Bearer " + API_KEY);
      
      HttpResponse response = httpClient.execute(request);
      String integrations = EntityUtils.toString(response.getEntity());
      
      renderRequest.setAttribute("integrations", integrations);
      renderRequest.setAttribute("apiGatewayUrl", API_GATEWAY_URL);
      
    } catch (Exception e) {
      renderRequest.setAttribute("error", e.getMessage());
    }
    
    super.doView(renderRequest, renderResponse);
  }

  public void createContent(ActionRequest actionRequest, ActionResponse actionResponse)
    throws Exception {
    
    String title = actionRequest.getParameter("title");
    String content = actionRequest.getParameter("content");
    String[] integrations = actionRequest.getParameterValues("integrations");
    
    // Prepare API request
    String createUrl = API_GATEWAY_URL + "/api/v1/content/create";
    
    JSONObject json = new JSONObject();
    json.put("title", title);
    json.put("content", content);
    json.put("integrations", integrations);
    
    HttpClient httpClient = HttpClientBuilder.create().build();
    HttpPost request = new HttpPost(createUrl);
    request.setHeader("Content-Type", "application/json");
    request.setHeader("Authorization", "Bearer " + API_KEY);
    request.setEntity(new StringEntity(json.toString()));
    
    HttpResponse response = httpClient.execute(request);
    String result = EntityUtils.toString(response.getEntity());
    
    // Store content ID in session
    PortletSession session = actionRequest.getPortletSession();
    JSONObject resultJson = new JSONObject(result);
    session.setAttribute("lastContentId", resultJson.optString("id"));
    
    // Redirect to formats page
    actionResponse.sendRedirect(actionRequest.getRenderURL() + "?tab=formats");
  }

  public void convertToFormat(ActionRequest actionRequest, ActionResponse actionResponse)
    throws Exception {
    
    PortletSession session = actionRequest.getPortletSession();
    String contentId = (String) session.getAttribute("lastContentId");
    String format = actionRequest.getParameter("format"); // epub, pdf, slides, audio, video
    
    String formatUrl = API_GATEWAY_URL + "/api/v1/formats/" + format;
    
    JSONObject json = new JSONObject();
    json.put("contentId", contentId);
    
    HttpClient httpClient = HttpClientBuilder.create().build();
    HttpPost request = new HttpPost(formatUrl);
    request.setHeader("Content-Type", "application/json");
    request.setHeader("Authorization", "Bearer " + API_KEY);
    request.setEntity(new StringEntity(json.toString()));
    
    HttpResponse response = httpClient.execute(request);
    String result = EntityUtils.toString(response.getEntity());
    
    JSONObject resultJson = new JSONObject(result);
    String downloadUrl = resultJson.optString("downloadUrl");
    
    actionRequest.setAttribute("downloadUrl", downloadUrl);
    actionResponse.sendRedirect(downloadUrl);
  }

  public void publishToIntegrations(ActionRequest actionRequest, ActionResponse actionResponse)
    throws Exception {
    
    PortletSession session = actionRequest.getPortletSession();
    String contentId = (String) session.getAttribute("lastContentId");
    String[] integrationIds = actionRequest.getParameterValues("integrationIds");
    
    for (String integrationId : integrationIds) {
      String publishUrl = API_GATEWAY_URL + "/api/v1/integrations/" + integrationId + "/publish";
      
      JSONObject json = new JSONObject();
      json.put("contentId", contentId);
      
      HttpClient httpClient = HttpClientBuilder.create().build();
      HttpPost request = new HttpPost(publishUrl);
      request.setHeader("Content-Type", "application/json");
      request.setHeader("Authorization", "Bearer " + API_KEY);
      request.setEntity(new StringEntity(json.toString()));
      
      HttpResponse response = httpClient.execute(request);
      // Log result
    }
    
    actionResponse.sendRedirect(actionRequest.getRenderURL() + "?tab=success");
  }
}
```

### JSP Template: create-content.jsp

```jsp
<%@ include file="/init.jsp" %>

<div class="content-creator-portlet">
  <h2>Create & Publish Content</h2>
  
  <% if (request.getAttribute("error") != null) { %>
    <div class="alert alert-danger">
      <%= request.getAttribute("error") %>
    </div>
  <% } %>

  <form action="<portlet:actionURL><portlet:param name="action" value="createContent" /></portlet:actionURL>" method="post">
    
    <div class="form-group">
      <label>Title</label>
      <input type="text" name="title" class="form-control" required />
    </div>

    <div class="form-group">
      <label>Content (Markdown)</label>
      <textarea name="content" class="form-control" rows="10" required></textarea>
    </div>

    <div class="form-group">
      <label>Publish To</label>
      <div class="integrations-list">
        <% 
          String integrations = (String) request.getAttribute("integrations");
          if (integrations != null) {
            JSONArray array = new JSONArray(integrations);
            for (int i = 0; i < array.length(); i++) {
              JSONObject integration = array.getJSONObject(i);
        %>
          <label class="checkbox">
            <input type="checkbox" name="integrations" value="<%= integration.getString("id") %>" />
            <%= integration.getString("name") %>
          </label>
        <% 
            }
          }
        %>
      </div>
    </div>

    <button type="submit" class="btn btn-primary">Create & Continue</button>
  </form>
</div>

<script>
// After creation, show format converter
document.addEventListener('DOMContentLoaded', function() {
  const urlParams = new URLSearchParams(window.location.search);
  if (urlParams.get('tab') === 'formats') {
    // Show format converter section
    document.getElementById('formats-section').style.display = 'block';
  }
});
</script>
```

---

## Step 3: Create Configuration File

### File: `liferay-config.properties`

```properties
# ============================================
# OpenAutonomyX Configuration
# ============================================

# API Gateway
openautonomyx.api.gateway.url=http://api-gateway:3000
openautonomyx.api.key=your-secure-api-key-here

# Services
openautonomyx.blog.url=http://blog:3009
openautonomyx.formats.url=http://formats:3011
openautonomyx.integrations.url=http://integrations:3010
openautonomyx.analytics.url=http://analytics:3005
openautonomyx.content.url=http://content-management:3002

# Database
openautonomyx.database.url=postgresql://postgres:5432/liferay
openautonomyx.database.user=postgres
openautonomyx.database.password=secure_password

# Cache
openautonomyx.cache.url=redis://redis:6379

# Timeouts
openautonomyx.api.timeout=30000
openautonomyx.api.retry.attempts=3

# CORS
cors.allowed.origins=localhost:8080,publishing.openautonomyx.com

# Environment
openautonomyx.environment=production
openautonomyx.debug=false
```

---

## Step 4: Create Environment Setup Script

### File: `setup-integration.sh`

```bash
#!/bin/bash

echo "🔌 Wiring Platform Together..."

# 1. Create Liferay portlet module
echo "📦 Creating Content Creator Portlet..."
blade create -t mvc-portlet \
  -p com.openautonomyx \
  content-creator-portlet

# 2. Copy API client
echo "📋 Setting up API Client..."
cp api-client.ts content-creator-portlet/src/main/java/api/

# 3. Copy portlet code
echo "🔗 Installing Portlet..."
cp ContentCreatorPortlet.java content-creator-portlet/src/main/java/
cp create-content.jsp content-creator-portlet/src/main/resources/

# 4. Build portlet
echo "🔨 Building Portlet..."
cd content-creator-portlet
./gradlew build
cp build/libs/content-creator-portlet.jar $LIFERAY_HOME/deploy/

# 5. Wait for deployment
echo "⏳ Waiting for deployment..."
sleep 10

# 6. Verify connection
echo "✅ Testing API Connection..."
curl -X GET http://localhost:3000/health \
  -H "Authorization: Bearer dev-key"

if [ $? -eq 0 ]; then
  echo "✅ Platform Wired Together Successfully!"
else
  echo "❌ Connection failed. Check API Gateway is running."
fi
```

---

## Step 5: Docker Compose - Fully Wired

```yaml
version: '3.9'

services:
  # Liferay Portal
  liferay:
    image: liferay/dxp:latest
    ports:
      - "8080:8080"
    environment:
      LIFERAY_JDBC_ONE_URL: jdbc:postgresql://postgres:5432/liferay
      LIFERAY_JDBC_ONE_USERNAME: postgres
      LIFERAY_JDBC_ONE_PASSWORD: secure_password
      OPENAUTONOMYX_API_GATEWAY_URL: http://api-gateway:3000
      OPENAUTONOMYX_API_KEY: dev-key
    depends_on:
      - postgres
      - api-gateway
    volumes:
      - ./liferay-config.properties:/opt/liferay/portal-ext.properties
      - ./content-creator-portlet/build/libs:/opt/liferay/deploy
    networks:
      - platform-network

  # PostgreSQL
  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: liferay
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: secure_password
    volumes:
      - postgres_data:/var/lib/postgresql/data
    networks:
      - platform-network

  # Redis
  redis:
    image: redis:7-alpine
    networks:
      - platform-network

  # API Gateway
  api-gateway:
    build: ./services/api-gateway
    ports:
      - "3000:3000"
    environment:
      DATABASE_URL: postgresql://postgres:secure_password@postgres:5432/liferay
      REDIS_URL: redis://redis:6379
      NODE_ENV: production
    depends_on:
      - postgres
      - redis
    networks:
      - platform-network

  # Blog Service
  blog:
    build: ./services/blog
    ports:
      - "3009:3009"
    environment:
      DATABASE_URL: postgresql://postgres:secure_password@postgres:5432/liferay
      REDIS_URL: redis://redis:6379
      EVENT_BUS_URL: http://event-bus:3001
    depends_on:
      - postgres
      - redis
    networks:
      - platform-network

  # Formats Service
  formats:
    build: ./services/formats
    ports:
      - "3011:3011"
    networks:
      - platform-network

  # Integrations Service
  integrations:
    build: ./services/integrations
    ports:
      - "3010:3010"
    networks:
      - platform-network

  # Event Bus
  event-bus:
    build: ./services/event-bus
    ports:
      - "3001:3001"
    environment:
      REDIS_URL: redis://redis:6379
    depends_on:
      - redis
    networks:
      - platform-network

volumes:
  postgres_data:

networks:
  platform-network:
    driver: bridge
```

---

## Step 6: Complete Data Flow

```
User in Liferay Portal
    ↓
[Content Creator Portlet UI]
    ↓
Form Submit (Create Content)
    ↓
API Client sends to http://api-gateway:3000/api/v1/content/create
    ↓
API Gateway routes to Content Service (3002)
    ↓
Content Service:
  • Creates in PostgreSQL
  • Caches in Redis
  • Emits event to Event Bus
    ↓
[User selects Format: EPUB]
    ↓
API Client sends to http://formats:3011/api/v1/formats/epub
    ↓
Formats Service generates EPUB
    ↓
[User selects Integrations: WordPress, Medium, Twitter]
    ↓
API Client sends to http://integrations:3010/api/v1/integrations/{id}/publish
    ↓
Each Integration Service publishes to platform:
  • WordPress API call
  • Medium API call
  • Twitter API call
    ↓
✅ Content published everywhere
    ↓
User sees success in Liferay Portal
```

---

## Testing the Wire

```bash
# 1. Start everything
docker-compose up -d

# 2. Wait for startup
sleep 30

# 3. Test API Gateway
curl http://localhost:3000/health

# 4. Test Blog Service
curl http://localhost:3009/health

# 5. Test Formats Service
curl http://localhost:3011/health

# 6. Test Liferay
curl http://localhost:8080

# 7. Create test content
curl -X POST http://localhost:3000/api/v1/content/create \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Test Post",
    "content": "Testing the wired platform"
  }'

# 8. Convert to EPUB
curl -X POST http://localhost:3011/api/v1/formats/epub \
  -H "Content-Type: application/json" \
  -d '{"contentId": "test-id", "title": "Test", "content": "..."}'

# 9. Publish to integration
curl -X POST http://localhost:3010/api/v1/integrations/wordpress/publish \
  -H "Content-Type: application/json" \
  -d '{"title": "Test", "content": "..."}'
```

---

## Final Checklist

- [ ] API Client library created
- [ ] Content Creator Portlet built
- [ ] Liferay config file set up
- [ ] Docker Compose fully configured
- [ ] All services deployed
- [ ] Portlet deployed to Liferay
- [ ] Test content creation works
- [ ] Test format conversion works
- [ ] Test integration publishing works
- [ ] Analytics tracking works
- [ ] Full end-to-end flow complete

---

**Platform Fully Wired!** 🚀

Liferay Portal ↔ API Gateway ↔ All Microservices ↔ PostgreSQL/Redis/Ollama
