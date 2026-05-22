# Model Parts

These files are the building blocks for a model document.

They are not a template engine. They exist so a project can either:

- copy one part and adjust it
- copy the entire model template and adapt the result

## Part order

1. `header.md`
2. `purpose.md`
3. `data-structure.md`
4. `responsibilities.md`
5. `lifecycle.md`
6. `validation.md`
7. `references.md`
8. `open-questions.md`

## Customization rule

If a project needs a special field layout, such as an extra `Master table` column, edit the copied `data-structure.md` part or copy the whole model template and adapt it as a unit.

Do not try to infer automatic merging between the default files and a project-local copy. The customization is explicit and visible.
