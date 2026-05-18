-- Seed data for Publishing Platform RBAC.
-- Idempotently creates canonical tenant roles and permissions.

UPSERT role:owner CONTENT {
  name: "owner",
  permissions: [
    "tenant:manage",
    "membership:manage",
    "article:create",
    "article:edit",
    "article:publish",
    "article:delete",
    "webhook:manage",
    "media:manage"
  ],
  created_at: time::now()
};

UPSERT role:admin CONTENT {
  name: "admin",
  permissions: [
    "tenant:manage",
    "membership:manage",
    "article:create",
    "article:edit",
    "article:publish",
    "article:delete",
    "webhook:manage",
    "media:manage"
  ],
  created_at: time::now()
};

UPSERT role:editor CONTENT {
  name: "editor",
  permissions: [
    "article:create",
    "article:edit",
    "article:publish",
    "webhook:view",
    "media:manage"
  ],
  created_at: time::now()
};

UPSERT role:author CONTENT {
  name: "author",
  permissions: [
    "article:create",
    "article:edit",
    "media:upload"
  ],
  created_at: time::now()
};

UPSERT role:viewer CONTENT {
  name: "viewer",
  permissions: [
    "article:view",
    "media:view"
  ],
  created_at: time::now()
};
