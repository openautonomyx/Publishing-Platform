-- Core SurrealDB schema for Publishing Platform.
-- Uses SurrealDB built-in functions for timestamps, ids, slugs, type checks, and defaults.

DEFINE TABLE tenant SCHEMAFULL;
DEFINE FIELD name ON tenant TYPE string ASSERT string::len($value) > 0;
DEFINE FIELD slug ON tenant TYPE string VALUE string::slug($value OR name);
DEFINE FIELD created_at ON tenant TYPE datetime VALUE time::now();
DEFINE FIELD updated_at ON tenant TYPE datetime VALUE time::now();
DEFINE INDEX tenant_slug_unique ON tenant COLUMNS slug UNIQUE;

DEFINE TABLE user SCHEMAFULL;
DEFINE FIELD email ON user TYPE string ASSERT string::is::email($value);
DEFINE FIELD name ON user TYPE option<string>;
DEFINE FIELD created_at ON user TYPE datetime VALUE time::now();
DEFINE FIELD updated_at ON user TYPE datetime VALUE time::now();
DEFINE INDEX user_email_unique ON user COLUMNS email UNIQUE;

DEFINE TABLE role SCHEMAFULL;
DEFINE FIELD name ON role TYPE string ASSERT $value INSIDE ["owner", "admin", "editor", "author", "viewer"];
DEFINE FIELD permissions ON role TYPE array<string> DEFAULT [];
DEFINE FIELD created_at ON role TYPE datetime VALUE time::now();
DEFINE INDEX role_name_unique ON role COLUMNS name UNIQUE;

DEFINE TABLE membership SCHEMAFULL;
DEFINE FIELD tenant ON membership TYPE record<tenant>;
DEFINE FIELD user ON membership TYPE record<user>;
DEFINE FIELD role ON membership TYPE record<role>;
DEFINE FIELD created_at ON membership TYPE datetime VALUE time::now();
DEFINE FIELD updated_at ON membership TYPE datetime VALUE time::now();
DEFINE INDEX membership_unique ON membership COLUMNS tenant, user UNIQUE;

DEFINE TABLE article SCHEMAFULL;
DEFINE FIELD tenant ON article TYPE record<tenant>;
DEFINE FIELD author ON article TYPE record<user>;
DEFINE FIELD title ON article TYPE string ASSERT string::len($value) > 0;
DEFINE FIELD slug ON article TYPE string VALUE string::slug($value OR title);
DEFINE FIELD status ON article TYPE string DEFAULT "draft" ASSERT $value INSIDE ["draft", "review", "published", "archived"];
DEFINE FIELD excerpt ON article TYPE option<string>;
DEFINE FIELD body ON article TYPE string DEFAULT "";
DEFINE FIELD language ON article TYPE string DEFAULT "en";
DEFINE FIELD canonical_url ON article TYPE option<string>;
DEFINE FIELD published_at ON article TYPE option<datetime>;
DEFINE FIELD created_at ON article TYPE datetime VALUE time::now();
DEFINE FIELD updated_at ON article TYPE datetime VALUE time::now();
DEFINE INDEX article_tenant_slug_unique ON article COLUMNS tenant, slug UNIQUE;
DEFINE INDEX article_status_idx ON article COLUMNS status;

DEFINE TABLE webhook_endpoint SCHEMAFULL;
DEFINE FIELD tenant ON webhook_endpoint TYPE record<tenant>;
DEFINE FIELD url ON webhook_endpoint TYPE string ASSERT string::is::url($value);
DEFINE FIELD secret ON webhook_endpoint TYPE string DEFAULT rand::guid();
DEFINE FIELD enabled ON webhook_endpoint TYPE bool DEFAULT true;
DEFINE FIELD created_at ON webhook_endpoint TYPE datetime VALUE time::now();
DEFINE FIELD updated_at ON webhook_endpoint TYPE datetime VALUE time::now();

DEFINE TABLE automation_event SCHEMAFULL;
DEFINE FIELD tenant ON automation_event TYPE record<tenant>;
DEFINE FIELD event_type ON automation_event TYPE string;
DEFINE FIELD payload ON automation_event TYPE object DEFAULT {};
DEFINE FIELD occurred_at ON automation_event TYPE datetime VALUE time::now();
DEFINE INDEX automation_event_type_idx ON automation_event COLUMNS event_type;

DEFINE TABLE webhook_delivery SCHEMAFULL;
DEFINE FIELD endpoint ON webhook_delivery TYPE record<webhook_endpoint>;
DEFINE FIELD event ON webhook_delivery TYPE record<automation_event>;
DEFINE FIELD status ON webhook_delivery TYPE string DEFAULT "pending" ASSERT $value INSIDE ["pending", "delivered", "failed"];
DEFINE FIELD attempt_count ON webhook_delivery TYPE int DEFAULT 0;
DEFINE FIELD next_attempt_at ON webhook_delivery TYPE option<datetime>;
DEFINE FIELD last_error ON webhook_delivery TYPE option<string>;
DEFINE FIELD created_at ON webhook_delivery TYPE datetime VALUE time::now();
DEFINE FIELD updated_at ON webhook_delivery TYPE datetime VALUE time::now();
