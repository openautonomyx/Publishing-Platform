# Schema Grounding

The platform schema should be grounded in public vocabularies and portable identifiers, not private platform-only concepts.

## Core Principle

Internal records may live in SurrealDB, but public meaning should map to open semantic models.

```text
SurrealDB record   = storage and graph representation
Schema.org entity  = public semantic representation
JSON-LD            = public serialization format
Plus Code          = canonical portable geo identity
```

## Schema.org

Schema.org is the primary vocabulary for public structured data.

The platform should map publishing records to Schema.org types wherever possible.

Reference:

- https://schema.org

## JSON-LD

JSON-LD is the public serialization format for Schema.org entities.

Every public page should be able to emit valid JSON-LD.

Reference:

- https://json-ld.org

## Core Schema Types

### Thing

Base type for public semantic entities.

Used for:

- Shared identifiers
- Names
- URLs
- Descriptions
- Images
- Same-as links

Schema.org reference:

- https://schema.org/Thing

### Person

Used for authors, contributors, editors, experts, analysts, and public profiles.

Schema.org reference:

- https://schema.org/Person

### Organization

Used for tenants, publishers, research organizations, institutions, media organizations, and educational organizations.

Schema.org references:

- https://schema.org/Organization
- https://schema.org/NewsMediaOrganization
- https://schema.org/ResearchOrganization
- https://schema.org/EducationalOrganization

### Place

Used for location as a primitive.

The platform should use Schema.org `Place` for public structured data and Plus Codes as canonical location identifiers.

Schema.org references:

- https://schema.org/Place
- https://schema.org/GeoCoordinates
- https://schema.org/PostalAddress

Open Location Code reference:

- https://github.com/google/open-location-code

### CreativeWork

Base type for publishable works.

All publishing content should map to `CreativeWork` or a subtype.

Schema.org reference:

- https://schema.org/CreativeWork

## CreativeWork Subtypes

The platform should support many publishable content types, including:

- Article
- NewsArticle
- OpinionNewsArticle
- AnalysisNewsArticle
- ScholarlyArticle
- TechArticle
- BlogPosting
- Report
- Dataset
- DataCatalog
- Book
- Chapter
- WebPage
- FAQPage
- ProfilePage
- HowTo
- Course
- LearningResource
- DigitalDocument
- PresentationDigitalDocument
- ImageObject
- VideoObject
- AudioObject
- MediaObject
- PodcastEpisode
- PodcastSeries
- Review
- ClaimReview
- Comment
- DiscussionForumPosting
- SoftwareSourceCode
- VisualArtwork
- Photograph

## Internal to Public Mapping

### Tenant

```text
tenant -> Organization / NewsMediaOrganization / ResearchOrganization
```

### User / Author

```text
user -> Person
```

### Article

```text
article -> Article / NewsArticle / OpinionNewsArticle / CreativeWork subtype
```

### Location

```text
location -> Place
location.plusCode -> Place.@id / Place.identifier
location.coordinates -> Place.geo
```

### Media Asset

```text
media_asset -> ImageObject / VideoObject / AudioObject / MediaObject
```

### Dataset

```text
dataset -> Dataset
```

### Feed / Collection

```text
collection -> CollectionPage / CreativeWorkSeries
```

## Location Identity

Location identity should not be locked to SurrealDB record IDs.

Use Plus Codes as canonical location identifiers.

```text
Plus Code        = canonical location identity
GeoCoordinates   = actual latitude/longitude
Schema.org Place = public semantic object
SurrealDB record = internal graph node
```

## Package

The TypeScript schema helpers live in:

```text
packages/schema
```

Key exports:

- `Thing`
- `Person`
- `Organization`
- `Place`
- `GeoCoordinates`
- `PostalAddress`
- `CreativeWork`
- `Article`
- `Dataset`
- `MediaObject`
- `BreadcrumbList`
- `LinkedGeoObject`
- `withContext()`
- `toJsonLd()`
- `createPlaceFromPlusCode()`
- `createLinkedGeoObject()`

## Implementation Rule

Every public publishing object should be able to answer:

1. What is its internal SurrealDB record?
2. What is its Schema.org type?
3. What JSON-LD should it emit?
4. What external identifier makes it portable?

For location, that external identifier is the Plus Code.
