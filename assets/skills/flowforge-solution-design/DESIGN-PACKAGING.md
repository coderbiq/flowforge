# Adaptive design packaging

Start with one design hub containing navigation, cross-cutting decisions, and compact design areas.

Promote an area into `design/<area>.md` only when it has independent consumers, independent review or lifecycle, several unrelated revisions, cross-ticket scope, or enough detail to obscure other decisions. The hub keeps a semantic summary and link rather than copying child authority.

Use one stable area identity and revision until independent consumption demonstrates a need for finer granularity. Splitting a file preserves semantic identity and consumer links.

