// Package factoriat implements an information-transistor pattern.
//
// A Factoriat receives input factors as signals, accumulates them into live
// memory, evaluates an activation gate, builds facts from an independent
// memory snapshot, stabilizes live memory, and emits the facts outside the
// internal lock.
//
// Lifecycle:
//
//	Capture   -> write the input factor into live memory
//	Evaluate  -> run the activation gate
//	Build     -> create facts from a State snapshot
//	Stabilize -> move live memory into its next stable shape
//	Emit      -> publish each fact after the lock is released
//
// Invariants:
//
//   - Build is always required by Config.Validate.
//   - Capture, Evaluate, and Emit are required by Config.Validate only when Stateful is true.
//   - Evaluate is the activation gate: it decides whether the current state is active.
//   - Build receives an independent State snapshot and must not control live memory.
//   - Stabilize receives live memory and is the only post-fact state transition hook.
//   - Emit runs after the internal mutex is unlocked, once per built fact, in order.
//   - PushResult reports runtime outcomes and optional stage errors through Err.
//   - StateRepository persists StateRecord, including Triggered for edge mode.
//
// Each lifecycle hook returns an error. A non-nil error stops the lifecycle
// with PushStatusFailed.
package factoriat
