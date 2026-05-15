import { getDb } from './client';

export async function upsertGeoObject(input: {
  plusCode: string;
  name?: string;
  level?: string;
  coordinates: { latitude: number; longitude: number; elevation?: number };
  schema: Record<string, unknown>;
  parentId?: string;
}) {
  const db = await getDb();
  const id = input.plusCode.replace(/[^a-zA-Z0-9]/g, '_');

  return db.query(
    `UPSERT type::thing('geo_object', $id) CONTENT {
      plus_code: $plusCode,
      name: $name,
      level: $level,
      coordinates: $coordinates,
      schema: $schema,
      parent: IF $parentId THEN type::thing('geo_object', $parentId) ELSE NONE END,
      created_at: time::now()
    };`,
    {
      id,
      plusCode: input.plusCode,
      name: input.name,
      level: input.level,
      coordinates: input.coordinates,
      schema: input.schema,
      parentId: input.parentId ?? null,
    },
  );
}

export async function getGeoObjectByPlusCode(plusCode: string) {
  const db = await getDb();
  const [rows] = await db.query<[unknown[]]>(
    'SELECT * FROM geo_object WHERE plus_code = $plusCode LIMIT 1;',
    { plusCode },
  );

  return rows?.[0] ?? null;
}

export async function saveSchemaEntity(input: {
  tenantId?: string;
  subject?: string;
  schemaType: string;
  jsonld: Record<string, unknown>;
}) {
  const db = await getDb();

  return db.create('schema_entity', {
    tenant: input.tenantId ? `tenant:${input.tenantId}` : undefined,
    subject: input.subject,
    schema_type: input.schemaType,
    jsonld: input.jsonld,
    created_at: new Date().toISOString(),
    updated_at: new Date().toISOString(),
  });
}
