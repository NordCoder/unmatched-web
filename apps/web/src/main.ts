/**
 * Browser bootstrap metadata only. The client renders server projections and
 * submits commands; it does not own authoritative legality or resolution.
 */
export interface BootstrapDescriptor {
  readonly clientRole: "projection-renderer";
  readonly authority: "server";
  readonly runtimeConfigured: false;
}

export const bootstrapDescriptor = Object.freeze({
  clientRole: "projection-renderer",
  authority: "server",
  runtimeConfigured: false,
}) satisfies BootstrapDescriptor;
