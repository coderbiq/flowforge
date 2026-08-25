# Mocking

Use a fake or mock only at an external boundary that is slow, nondeterministic, unavailable, or destructive. Assert the returned behavior first. Verify interaction details only when the interaction itself is the contract; never mock private collaborators to preserve an implementation shape.

