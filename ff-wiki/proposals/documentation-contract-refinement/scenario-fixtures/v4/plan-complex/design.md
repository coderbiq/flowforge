---
flowforge:
  schema: 1
  role: design
  areas:
    provider-composition-seam:
      revision: 4
      anchor: provider-composition-seam
---
# Provider-composition design fixture

<a id="provider-composition-seam"></a>
## Provider composition seam

Application composition owns provider selection; the web adapter owns concrete construction. Migration expands compatibility, moves consumers independently, then removes the legacy path.
