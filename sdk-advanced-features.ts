/**
 * Advanced Features SDK
 * Penpot Integration | Template Library | Liferay Structures | Multi-Tenant
 */

export interface PenpotProject {
  id: string;
  name: string;
  team_id: string;
  created_at?: string;
}

export interface DesignFile {
  id: string;
  project_id: string;
  name: string;
  created_at?: string;
}

export interface Template {
  id: string;
  name: string;
  category: 'marketing' | 'editorial' | 'social' | 'landing' | 'email';
  tenant_id: string;
  data: Record<string, any>;
  tags: string[];
  created_at?: string;
}

export interface FieldDef {
  name: string;
  type: 'text' | 'rich_text' | 'image' | 'decimal' | 'select' | 'date';
  required?: boolean;
  repeatable?: boolean;
  options?: string[];
}

export interface LiferayStructure {
  id: string;
  name: string;
  fields: FieldDef[];
  definition: Record<string, any>;
  created_at?: string;
}

export interface TenantSettings {
  branding?: {
    logo_url?: string;
    primary_color?: string;
    secondary_color?: string;
  };
  features?: string[];
  [key: string]: any;
}

export interface TenantConfig {
  id: string;
  name: string;
  subdomain: string;
  database: string;
  storage: string;
  settings: TenantSettings;
  users: string[];
  custom_domains: string[];
  created_at?: string;
}

export class PenpotClient {
  private baseUrl: string;
  private apiToken: string;

  constructor(baseUrl: string, apiToken: string) {
    this.baseUrl = baseUrl;
    this.apiToken = apiToken;
  }

  async createProject(name: string, teamId: string): Promise<PenpotProject> {
    const response = await fetch(`${this.baseUrl}/api/penpot/projects`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${this.apiToken}`,
      },
      body: JSON.stringify({ name, team_id: teamId }),
    });
    return response.json();
  }

  async createDesignFile(projectId: string, name: string): Promise<DesignFile> {
    const response = await fetch(`${this.baseUrl}/api/penpot/files`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${this.apiToken}`,
      },
      body: JSON.stringify({ project_id: projectId, name }),
    });
    return response.json();
  }

  async exportDesign(fileId: string, format: 'svg' | 'png' = 'svg'): Promise<Blob> {
    const response = await fetch(
      `${this.baseUrl}/api/penpot/export?file_id=${fileId}&format=${format}`,
      {
        headers: { 'Authorization': `Bearer ${this.apiToken}` },
      }
    );
    return response.blob();
  }
}

export class TemplateLibrary {
  private baseUrl: string;
  private tenantId: string;

  constructor(baseUrl: string, tenantId: string) {
    this.baseUrl = baseUrl;
    this.tenantId = tenantId;
  }

  async addTemplate(
    category: Template['category'],
    name: string,
    data: Record<string, any>,
    tags: string[] = []
  ): Promise<Template> {
    const response = await fetch(`${this.baseUrl}/api/templates/add`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        name,
        category,
        tenant_id: this.tenantId,
        data,
        tags,
      }),
    });
    return response.json();
  }

  async getTemplatesByCategory(category: Template['category']): Promise<Template[]> {
    const response = await fetch(
      `${this.baseUrl}/api/templates/category?category=${category}`,
      { method: 'GET' }
    );
    return response.json();
  }

  async searchTemplates(query: string): Promise<Template[]> {
    const response = await fetch(
      `${this.baseUrl}/api/templates/search?q=${encodeURIComponent(query)}`,
      { method: 'GET' }
    );
    return response.json();
  }

  async cloneTemplate(templateId: string, newName: string): Promise<Template> {
    const response = await fetch(`${this.baseUrl}/api/templates/clone`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ template_id: templateId, new_name: newName }),
    });
    return response.json();
  }
}

export class LiferayStructures {
  private baseUrl: string;

  constructor(baseUrl: string) {
    this.baseUrl = baseUrl;
  }

  async createStructure(
    name: string,
    fields: FieldDef[]
  ): Promise<LiferayStructure> {
    const response = await fetch(`${this.baseUrl}/api/structures/create`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ name, fields }),
    });
    return response.json();
  }

  async getStructure(id: string): Promise<LiferayStructure> {
    const response = await fetch(`${this.baseUrl}/api/structures/get?id=${id}`, {
      method: 'GET',
    });
    return response.json();
  }

  async listStructures(): Promise<LiferayStructure[]> {
    const response = await fetch(`${this.baseUrl}/api/structures/list`, {
      method: 'GET',
    });
    return response.json();
  }

  async createWebContent(
    structureId: string,
    title: string,
    content: Record<string, any>
  ): Promise<any> {
    const response = await fetch(`${this.baseUrl}/api/structures/content`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ structure_id: structureId, title, content }),
    });
    return response.json();
  }
}

export class MultiTenantPlatform {
  private baseUrl: string;
  private adminToken: string;

  constructor(baseUrl: string, adminToken: string) {
    this.baseUrl = baseUrl;
    this.adminToken = adminToken;
  }

  async registerTenant(config: Omit<TenantConfig, 'id' | 'created_at'>): Promise<TenantConfig> {
    const response = await fetch(`${this.baseUrl}/api/tenants/register`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${this.adminToken}`,
      },
      body: JSON.stringify(config),
    });
    return response.json();
  }

  async getTenantConfig(tenantId: string): Promise<TenantConfig> {
    const response = await fetch(
      `${this.baseUrl}/api/tenants/config?tenant_id=${tenantId}`,
      {
        method: 'GET',
        headers: { 'Authorization': `Bearer ${this.adminToken}` },
      }
    );
    return response.json();
  }

  async listTenants(): Promise<TenantConfig[]> {
    const response = await fetch(`${this.baseUrl}/api/tenants/list`, {
      method: 'GET',
      headers: { 'Authorization': `Bearer ${this.adminToken}` },
    });
    return response.json();
  }

  async updateTenantSettings(
    tenantId: string,
    settings: Partial<TenantSettings>
  ): Promise<TenantConfig> {
    const response = await fetch(`${this.baseUrl}/api/tenants/settings`, {
      method: 'PATCH',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${this.adminToken}`,
      },
      body: JSON.stringify({ tenant_id: tenantId, settings }),
    });
    return response.json();
  }

  async addTenantUser(tenantId: string, userId: string, role: string): Promise<any> {
    const response = await fetch(`${this.baseUrl}/api/tenants/users`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${this.adminToken}`,
      },
      body: JSON.stringify({ tenant_id: tenantId, user_id: userId, role }),
    });
    return response.json();
  }

  async addCustomDomain(tenantId: string, domain: string): Promise<any> {
    const response = await fetch(`${this.baseUrl}/api/tenants/domains`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${this.adminToken}`,
      },
      body: JSON.stringify({ tenant_id: tenantId, domain }),
    });
    return response.json();
  }
}

// Example usage
export async function initializeAdvancedFeatures(baseUrl: string, credentials: any) {
  const penpot = new PenpotClient(baseUrl, credentials.penpot_token);
  const templates = new TemplateLibrary(baseUrl, credentials.tenant_id);
  const structures = new LiferayStructures(baseUrl);
  const multiTenant = new MultiTenantPlatform(baseUrl, credentials.admin_token);

  return {
    penpot,
    templates,
    structures,
    multiTenant,
  };
}
