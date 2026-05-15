-- Publishing Platform SurrealDB schema
-- Core primitives: tenancy, users, sessions, RBAC, publishing, schema entities, and linked geo objects.

DEFINE TABLE tenant SCHEMAFULL;
DEFINE FIELD name ON tenant TYPE string;
DEFINE FIELD slug ON tenant TYPE string;
DEFINE FIELD domain ON tenant TYPE option<string>;
DEFINE FIELD settings ON tenant TYPE object DEFAULT {};
DEFINE FIELD created_at ON tenant TYPE datetime DEFAULT time::now();
DEFINE INDEX tenant_slug_unique ON tenant FIELDS slug UNIQUE;

DEFINE TABLE user SCHEMAFULL;
DEFINE FIELD email ON user TYPE string;
DEFINE FIELD name ON user TYPE option<string>;
DEFINE FIELD external_auth_id ON user TYPE option<string>;
DEFINE FIELD created_at ON user TYPE datetime DEFAULT time::now();
DEFINE INDEX user_email_unique ON user FIELDS email UNIQUE;

DEFINE TABLE membership SCHEMAFULL;
DEFINE FIELD tenant ON membership TYPE record<tenant>;
DEFINE FIELD user ON membership TYPE record<user>;
DEFINE FIELD role ON membership TYPE string ASSERT $value INSIDE ['owner', 'admin', 'editor', 'author', 'viewer'];
DEFINE FIELD created_at ON membership TYPE datetime DEFAULT time::now();
DEFINE INDEX membership_unique ON membership FIELDS tenant, user UNIQUE;

DEFINE TABLE session SCHEMAFULL;
DEFINE FIELD user ON session TYPE record<user>;
DEFINE FIELD tenant ON session TYPE option<record<tenant>>;
DEFINE FIELD token_hash ON session TYPE string;
DEFINE FIELD expires_at ON session TYPE datetime;
DEFINE FIELD created_at ON session TYPE datetime DEFAULT time::now();
DEFINE INDEX session_token_unique ON session FIELDS token_hash UNIQUE;

DEFINE TABLE article SCHEMAFULL;
DEFINE FIELD tenant ON article TYPE record<tenant>;
DEFINE FIELD author ON article TYPE record<user>;
DEFINE FIELD title ON article TYPE string;
DEFINE FIELD slug ON article TYPE string;
DEFINE FIELD status ON article TYPE string ASSERT $value INSIDE ['draft', 'in_review', 'approved', 'scheduled', 'published', 'archived'];
DEFINE FIELD content ON article TYPE object DEFAULT {};
DEFINE FIELD language ON article TYPE string DEFAULT 'en';
DEFINE FIELD location ON article TYPE option<record<geo_object>>;
DEFINE FIELD seo ON article TYPE object DEFAULT {};
DEFINE FIELD schema ON article TYPE object DEFAULT {};
DEFINE FIELD scheduled_for ON article TYPE option<datetime>;
DEFINE FIELD published_at ON article TYPE option<datetime>;
DEFINE FIELD updated_at ON article TYPE datetime DEFAULT time::now();
DEFINE INDEX article_tenant_slug_unique ON article FIELDS tenant, slug UNIQUE;

DEFINE TABLE article_revision SCHEMAFULL;
DEFINE FIELD tenant ON article_revision TYPE record<tenant>;
DEFINE FIELD article ON article_revision TYPE record<article>;
DEFINE FIELD actor ON article_revision TYPE record<user>;
DEFINE FIELD content ON article_revision TYPE object;
DEFINE FIELD created_at ON article_revision TYPE datetime DEFAULT time::now();

DEFINE TABLE publishing_event SCHEMAFULL;
DEFINE FIELD tenant ON publishing_event TYPE record<tenant>;
DEFINE FIELD article ON publishing_event TYPE record<article>;
DEFINE FIELD actor ON publishing_event TYPE option<record<user>>;
DEFINE FIELD from_status ON publishing_event TYPE option<string>;
DEFINE FIELD to_status ON publishing_event TYPE string;
DEFINE FIELD created_at ON publishing_event TYPE datetime DEFAULT time::now();

DEFINE TABLE geo_object SCHEMAFULL;
DEFINE FIELD plus_code ON geo_object TYPE string;
DEFINE FIELD name ON geo_object TYPE option<string>;
DEFINE FIELD level ON geo_object TYPE option<string>;
DEFINE FIELD coordinates ON geo_object TYPE object;
DEFINE FIELD schema ON geo_object TYPE object;
DEFINE FIELD parent ON geo_object TYPE option<record<geo_object>>;
DEFINE FIELD created_at ON geo_object TYPE datetime DEFAULT time::now();
DEFINE INDEX geo_plus_code_unique ON geo_object FIELDS plus_code UNIQUE;

DEFINE TABLE schema_entity SCHEMAFULL;
DEFINE FIELD tenant ON schema_entity TYPE option<record<tenant>>;
DEFINE FIELD subject ON schema_entity TYPE option<record>;
DEFINE FIELD schema_type ON schema_entity TYPE string;
DEFINE FIELD jsonld ON schema_entity TYPE object;
DEFINE FIELD created_at ON schema_entity TYPE datetime DEFAULT time::now();
DEFINE FIELD updated_at ON schema_entity TYPE datetime DEFAULT time::now();
