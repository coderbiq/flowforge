# EXAMPLE — Tangram Backend Responsibility and Composition Design

> Scenario fixture used to validate the `flowforge-solution-design` interface. It is not an approved Tangram design and MUST NOT be used as an implementation specification.

## Requirement context

Realize the backend refinement without allowing application code to depend on implementation adapters, without placing runtime policy inside contract modules, and without breaking configured provider construction or domain lifecycle behavior.

Requirement authority in the real scenario would be linked here rather than copied.

## Current facts that constrain the solution

- `application/tangram-demo-app` previously constructed web implementation types directly.
- `domain/tangram-agent` owns AgentScope orchestration but application composition currently reaches types intended to remain internal.
- Configured web adapters require provider-specific construction inputs such as provider configuration and priority; annotation scanning alone cannot supply those values.
- RPC dispatch currently needs both route discovery and invocation, while reflection handles and runtime beans do not belong in contracts.
- Task trigger lifecycle must start only after the application reaches its running state and stop during shutdown.
- The previous static architecture check could pass while compilation and runtime composition remained broken, so static dependency rules are not sufficient verification.

## Target module responsibilities

| Module | Owns | Does not own |
|---|---|---|
| Contract modules | Stable, framework-neutral declarations and passive values crossing module seams | Mutable builders, reflection handles, runtime beans, provider construction policy |
| Agent domain | Agent capability orchestration and its application-facing composition interface | Application configuration binding or concrete web adapters |
| Task domain | Trigger registration and lifecycle coordination | Scheduler-specific polling implementation |
| Implementation adapters | Concrete scheduler, RPC, persistence, and web behavior behind approved interfaces | Public domain interfaces or application assembly policy |
| Spring/container adapter | Discovery of statically constructible adapters and runtime route invocation mechanics | Domain decisions or provider-specific configuration ownership |
| Application composition | Binding deployment configuration and selecting/configuring adapters through approved composition interfaces | Importing implementation-internal types or exposing domain internals |

## Interfaces and seams

### Agent application composition

The Agent domain exposes one application-facing composition interface that accepts framework-neutral capability descriptions or approved contract interfaces. It hides internal registries and AgentScope implementation types.

- **Owner:** Agent domain.
- **Consumer:** Application composition.
- **Providers:** Agent domain internals and discovered Agent extension adapters.
- **Composition owner:** Application.
- **Invariant:** Application MUST NOT access `internal` registries or AgentScope implementation types.

The exact type shape remains a design question in this fixture; tickets that expose or consume this interface cannot be planned normally until it is selected.

### Configured web adapter construction

Application composition owns provider selection and configuration binding. A web adapter factory interface accepts passive provider settings and returns enabled capability adapters. Concrete construction remains inside the web implementation module.

- **Owner:** Web extension contract owns the factory declaration; the implementation adapter provides it.
- **Consumer and composition owner:** Application composition.
- **Invariant:** Disabled providers MUST NOT be constructed and MUST NOT prevent startup.
- **Invariant:** Application MUST NOT import or instantiate implementation classes.

This simulated decision replaces the invalid assumption that annotation scanning can construct adapters requiring runtime configuration.

### RPC route execution

Contracts expose framework-neutral route identity, input, and result values. Spring/container implementation owns bean discovery, reflection, and invocation. Application dispatch crosses one route-execution interface rather than receiving a bean or `Method` handle.

- **Owner:** RPC extension contract.
- **Consumer:** Application dispatcher.
- **Provider:** Spring RPC adapter.
- **Invariant:** Contract values MUST NOT contain runtime beans or reflection handles.

### Trigger lifecycle

The Task domain owns trigger lifecycle coordination. Scheduler supplies a Trigger adapter. Registration may occur earlier, but start and stop follow application lifecycle transitions owned by TaskEngine.

- **Owner:** Task domain.
- **Provider:** Scheduler adapter.
- **Invariant:** A registered trigger MUST NOT start before the application is running.
- **Invariant:** Every started trigger MUST be stopped during application shutdown.

## Migration shape

1. Expand contract-only interfaces and application-facing composition seams without removing current callers.
2. Prove each new seam with focused behavior tests.
3. Migrate application, Agent, RPC, scheduler, and configured web callers in batches that keep the build coherent.
4. Remove old implementation imports, leaked runtime values, and obsolete construction paths only after their consumers have migrated.
5. Run cross-module integration verification before declaring the refinement complete.

Planning may split independent Agent, RPC, trigger, and web tracer paths after their interfaces are settled. The final integration ticket is blocked by every migrated path.

## Verification strategy

- Compile the application and every changed contract/domain/adapter module together.
- Verify Agent capabilities through the application-facing composition interface, not internal registries.
- Start the application with enabled and disabled provider configurations and observe construction behavior.
- Exercise RPC dispatch through the public route-execution seam without inspecting reflection internals.
- Verify trigger start/stop at application lifecycle transitions.
- Run static dependency rules as an additional check, not as the completion proof.
- Run whitespace and repository-wide relevant tests before completion.

Exact commands belong in the real design once the repository's executable Gradle paths and environment are confirmed; their absence is a visible warning in this fixture.

## Credible rejected alternatives

- **Annotation scanning constructs every adapter:** rejected because configured web adapters require runtime values that scanning cannot derive.
- **Expose Agent registries publicly to application:** rejected because it makes domain internals part of the application interface and leaks AgentScope composition mechanics.
- **Put reflection handles in RPC contracts:** rejected because it couples framework-neutral contracts to runtime invocation implementation.

## Open design questions

### What is the Agent application composition interface shape?

The owner and invariants are settled, but the smallest framework-neutral interface has not been selected. This affects Agent composition and application migration work, not RPC, trigger, or configured web design.

### Which exact Gradle commands form the verification seam?

The required behaviors are known, but executable focused commands must be confirmed in a Java-enabled Tangram environment. This warns all implementation tickets but does not change module responsibility.

## Illustrative machine metadata

```yaml
design_questions:
  agent-application-composition-interface:
    affects:
      - agent-composition
      - application-agent-migration
    status: open
  tangram-focused-gradle-verification:
    affects:
      - agent-composition
      - rpc-routing
      - trigger-lifecycle
      - configured-web-construction
    status: open
    diagnostic: warning
```

The metadata is illustrative and not an approved schema.
