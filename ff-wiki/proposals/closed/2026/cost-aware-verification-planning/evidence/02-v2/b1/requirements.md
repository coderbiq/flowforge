---
flowforge:
  schema: 1
  role: requirement
  id: b1-requirements
  revision: 1
---
<a id="b1-requirements"></a>
# B1 fixed requirement: ordered order quote

Replace the single-SKU quote with one order quote. Input is a nonempty lines array of objects with string sku and positive integer quantity. Preserve line order and duplicate SKUs. Catalog A=125 cents and B=300 cents. Output has items (sku, quantity, lineCents) and totalCents; no legacy cents wrapper. No mutation of the caller's input. Invalid shape/quantity raises INVALID_ORDER; unknown catalog SKU raises UNKNOWN_SKU. CLI emits the complete result as JSON and exit 0, or no stdout, exact error message on stderr and exit 2. Fixed example B×2,A×1,B×1 yields ordered 600,125,300 and total1025.

Scope: the interface, implementation and CLI are one usable behavior. Keep the existing --help branch and help module untouched in the initial delivery. No network, dependencies, generated artifacts or production use. The later predeclared help spelling correction is an ordinary separate Fix within the same ticket, introduced only by the review owner after the initial delivery; its test is not an initial acceptance item.
